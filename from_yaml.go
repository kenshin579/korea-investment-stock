package kis

// NewClientFromYAML 은 YAML 파일에서 credentials/설정 로드 후 Client 생성.
// 추가 옵션은 functional options 로 override 가능.
func NewClientFromYAML(path string, opts ...Option) (*Client, error) {
	cfg, err := LoadConfigFromYAML(path)
	if err != nil {
		return nil, err
	}
	return newFromConfig(cfg, opts...)
}
