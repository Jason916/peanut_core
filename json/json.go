//jasonxu
package json

import (
	"net/http"
	"encoding/json"
)

type JsonResp struct {
	data interface{}
}

type ErrorMsg struct {
	Msg string `json:"message"`
}

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

func (j *JsonResp) PrettyPrint() ([]byte, error) {
	return json.MarshalIndent(&j.data, "", "  ")
}

func NewErrorMsg(msg string) *ErrorMsg {
	return &ErrorMsg{Msg: msg}
}