package domestic

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
)

// decimalType 은 필드가 shopspring decimal 인지 판별하기 위한 reflect 타입.
var decimalType = reflect.TypeOf(decimal.Decimal{})

// zeroKeyCache: reflect.Type -> map[jsonKey]0치환바이트. 타입별 1회 계산.
var zeroKeyCache sync.Map

// decodeLenientNumbers 는 KIS 가 무거래/거래정지 종목의 숫자 필드를 빈 문자열("")로
// 내려줄 때 파싱 실패를 막는다. target(구조체 포인터)의 숫자 필드(decimal.Decimal 및
// 정수/실수)에 한해 JSON 값이 "" 또는 null 이면 0 으로 치환한 뒤 표준 디코드한다.
// 문자열 필드는 절대 변형하지 않는다.
func decodeLenientNumbers(data []byte, target any) error {
	zeros := numericZeroKeys(reflect.TypeOf(target).Elem())
	if len(zeros) == 0 {
		return json.Unmarshal(data, target)
	}
	// 빈 문자열이 아예 없으면(대다수 정상 응답) map 왕복 없이 바로 디코드.
	if !bytes.Contains(data, []byte(`""`)) {
		return json.Unmarshal(data, target)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return json.Unmarshal(data, target)
	}
	changed := false
	for k, repl := range zeros {
		if raw, ok := m[k]; ok && isEmptyOrNull(raw) {
			m[k] = repl
			changed = true
		}
	}
	if !changed {
		return json.Unmarshal(data, target)
	}
	patched, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(patched, target)
}

func isEmptyOrNull(raw json.RawMessage) bool {
	s := strings.TrimSpace(string(raw))
	return s == `""` || s == "null"
}

// numericZeroKeys 는 구조체 t 의 숫자 필드 json 키 -> 0 치환 바이트 맵을 반환한다(캐시).
// decimal 및 ,string 숫자는 따옴표 "0", 평문 숫자는 bare 0 으로 치환한다.
// 주의: 최상위(flat) 필드만 검사한다. 현재 대상 리프 구조체들은 모두 평면 구조라 충분하다.
func numericZeroKeys(t reflect.Type) map[string]json.RawMessage {
	if v, ok := zeroKeyCache.Load(t); ok {
		return v.(map[string]json.RawMessage)
	}
	out := map[string]json.RawMessage{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name, opts := tag, ""
		if idx := strings.IndexByte(tag, ','); idx >= 0 {
			name, opts = tag[:idx], tag[idx+1:]
		}
		if name == "" || name == "-" {
			continue
		}
		switch {
		case f.Type == decimalType:
			out[name] = json.RawMessage(`"0"`)
		case isNumericKind(f.Type.Kind()):
			if hasOption(opts, "string") {
				out[name] = json.RawMessage(`"0"`)
			} else {
				out[name] = json.RawMessage(`0`)
			}
		}
	}
	zeroKeyCache.Store(t, out)
	return out
}

func hasOption(opts, want string) bool {
	for _, o := range strings.Split(opts, ",") {
		if o == want {
			return true
		}
	}
	return false
}

func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
