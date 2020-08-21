package gwjsoner

import json "github.com/json-iterator/go"

func Unmarshal(bytes []byte, out interface{}) error {
	return json.Unmarshal(bytes, &out)
}

func Marshal(out interface{}) ([]byte, error) {
	return json.Marshal(&out)
}
