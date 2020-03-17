package message

import "encoding/json"

type Message struct {
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Body   map[string]interface{} `json:"body"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	marshaled, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return marshaled, nil
}

func (m *Message) UnmarshalBinary(bytes []byte) error {
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	return nil
}
