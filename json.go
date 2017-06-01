package apidoc

import (
	"bytes"
	"encoding/json"
)

// JSONPrettyPrint print rich json
func JSONPrettyPrint(in []byte) ([]byte, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, in, "", "  "); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
