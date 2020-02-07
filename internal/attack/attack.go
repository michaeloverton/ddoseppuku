package attack

import "encoding/json"

type Attack struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	Body   struct {
		URL string `json:"url"`
	} `json:"body"`
}

func (a Attack) MarshalBinary() ([]byte, error) {
	marshaled, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return marshaled, nil
}

func (a *Attack) UnmarshalBinary(bytes []byte) error {
	err := json.Unmarshal(bytes, &a)
	if err != nil {
		return err
	}
	return nil
}
