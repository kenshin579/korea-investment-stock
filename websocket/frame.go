package websocket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type frameKind int

const (
	frameKindUnknown frameKind = iota
	frameKindRealtime
	frameKindJSON
	frameKindPingPong
)

// frame 은 파싱된 raw WebSocket 메시지.
type frame struct {
	Kind      frameKind
	Encrypted bool     // 0/1 flag (realtime 만 의미)
	TrID      string   // realtime 만
	Count     int      // realtime 만
	Fields    []string // realtime caret-separated payload
	JSON      jsonFrame
}

type jsonFrame struct {
	Header struct {
		TrID string `json:"tr_id"`
	} `json:"header"`
	Body struct {
		RtCd  string `json:"rt_cd"`
		MsgCd string `json:"msg_cd"`
		Msg1  string `json:"msg1"`
	} `json:"body"`
	// 호환성 위해 top-level 필드도 받아둠
	RtCd  string `json:"rt_cd,omitempty"`
	MsgCd string `json:"msg_cd,omitempty"`
	Msg1  string `json:"msg1,omitempty"`
}

// parseFrame 은 raw text 를 frame 으로 파싱. 첫 글자로 종류 분기.
func parseFrame(raw string) (frame, error) {
	if len(raw) == 0 {
		return frame{}, fmt.Errorf("%w: empty", ErrWSInvalidFrame)
	}
	switch raw[0] {
	case '{':
		return parseJSONFrame(raw)
	case '0', '1':
		return parseRealtimeFrame(raw)
	default:
		return frame{}, fmt.Errorf("%w: unknown leader %q", ErrWSInvalidFrame, raw[0])
	}
}

func parseRealtimeFrame(raw string) (frame, error) {
	parts := strings.SplitN(raw, "|", 4)
	if len(parts) != 4 {
		return frame{}, fmt.Errorf("%w: realtime frame requires 4 pipe-parts", ErrWSInvalidFrame)
	}
	encrypted := parts[0] == "1"
	count, err := strconv.Atoi(parts[2])
	if err != nil {
		return frame{}, fmt.Errorf("%w: bad count %q", ErrWSInvalidFrame, parts[2])
	}
	fields := strings.Split(parts[3], "^")
	return frame{
		Kind:      frameKindRealtime,
		Encrypted: encrypted,
		TrID:      parts[1],
		Count:     count,
		Fields:    fields,
	}, nil
}

func parseJSONFrame(raw string) (frame, error) {
	var jf jsonFrame
	if err := json.Unmarshal([]byte(raw), &jf); err != nil {
		return frame{}, fmt.Errorf("%w: json: %v", ErrWSInvalidFrame, err)
	}
	// PINGPONG 분기
	if jf.Header.TrID == "PINGPONG" {
		return frame{Kind: frameKindPingPong, JSON: jf}, nil
	}
	// body 우선, top-level fallback
	if jf.Body.RtCd == "" && jf.RtCd != "" {
		jf.Body.RtCd = jf.RtCd
		jf.Body.MsgCd = jf.MsgCd
		jf.Body.Msg1 = jf.Msg1
	}
	return frame{Kind: frameKindJSON, JSON: jf}, nil
}

// chunkFields 는 fields 를 count 개의 chunk 로 분리. 길이 mismatch 면 panic-free zero return.
func chunkFields(fields []string, count, fieldsPerChunk int) [][]string {
	chunks, _ := chunkFieldsErr(fields, count, fieldsPerChunk)
	return chunks
}

func chunkFieldsErr(fields []string, count, fieldsPerChunk int) ([][]string, error) {
	if count*fieldsPerChunk != len(fields) {
		return nil, fmt.Errorf("%w: chunk mismatch: expect %d*%d=%d, got %d",
			ErrWSInvalidFrame, count, fieldsPerChunk, count*fieldsPerChunk, len(fields))
	}
	out := make([][]string, count)
	for i := 0; i < count; i++ {
		out[i] = fields[i*fieldsPerChunk : (i+1)*fieldsPerChunk]
	}
	return out, nil
}

// IsSubscribeSuccess 는 jsonFrame 이 등록 성공 응답인지 확인.
func (j jsonFrame) IsSubscribeSuccess() bool {
	return j.Body.RtCd == "0" && strings.Contains(j.Body.Msg1, "SUBSCRIBE SUCCESS")
}
