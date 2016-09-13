package rest

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/cihub/seelog"
)

// TransactionsInRangeHandler Handles the home page request
func (s Service) TransactionsInRangeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := decodeBody(r.Body)
	if err != nil {
		log.Error("Error decoding body: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error decoding body"))
		return
	}

	return
}

func decodeBody(b io.Reader) (body map[string]interface{}, err error) {
	decoder := json.NewDecoder(b)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, log.Error("Error unmarshalling body: ", err)
	}
	return body, err
}
