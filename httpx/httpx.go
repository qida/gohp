package httpx

import (
	"bytes"
	"encoding/json"
)

func SetUrlEncode(url string) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(url)
	return bf.String()
}
