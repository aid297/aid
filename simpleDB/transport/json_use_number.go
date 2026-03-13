package transport

import jsoniter "github.com/json-iterator/go"

var jsonUseNumberAPI = jsoniter.Config{UseNumber: true}.Froze()

type JSONAnyMap map[string]any

func (m *JSONAnyMap) UnmarshalJSON(data []byte) error {
	if m == nil {
		return nil
	}
	var decoded map[string]any
	if err := jsonUseNumberAPI.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*m = decoded
	return nil
}

func (m JSONAnyMap) MarshalJSON() ([]byte, error) {
	return jsonUseNumberAPI.Marshal(map[string]any(m))
}

type JSONAnySlice []any

func (s *JSONAnySlice) UnmarshalJSON(data []byte) error {
	if s == nil {
		return nil
	}
	var decoded []any
	if err := jsonUseNumberAPI.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*s = decoded
	return nil
}

func (s JSONAnySlice) MarshalJSON() ([]byte, error) {
	return jsonUseNumberAPI.Marshal([]any(s))
}

