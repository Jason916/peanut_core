//jasonxu
package json

import (
	"net/http"
	"encoding/json"
)

func Json(respwrite http.ResponseWriter, statuscode int, data interface{}){
	respwrite.Header().Set("Content-Type", "application/json; charset=UTF-8")
	respwrite.WriteHeader(statuscode)
	body, err := json.Marshal(data)

	if err != nil {
		body = []byte(
			`{"error_message": "encoding json failed!"}`)
	}

	respwrite.Write(body)
}