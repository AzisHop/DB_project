package httpresponder

import (
"encoding/json"
"log"
"net/http"
)

func Respond(writer http.ResponseWriter, code int, info interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if info != nil {
		err := json.NewEncoder(writer).Encode(info)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
