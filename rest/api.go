package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/cihub/seelog"
)

// TransactionsInRangeHandler Handles the home page request
func (s Service) TransactionsInRangeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		log.Error("Error decoding body: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error decoding body"))
		return
	}

	start, sFound := body["start"].(string)
	end, eFound := body["end"].(string)
	if !eFound || !sFound {
		log.Error("Error asserting time type: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error asserting time type"))
		return
	}

	sTime, _ := time.Parse("01/02/2006", start)
	eTime, _ := time.Parse("01/02/2006", end)
	results, err := s.U.QueryDateRange(sTime, eTime)
	if err != nil {
		log.Error("Error querying date range: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error querying date range"))
	}

	rStr, err := json.Marshal(results)
	if err != nil {
		log.Error("Error marshalling to json: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling to json"))
	}

	w.Write(rStr)
	return
}

// TransactionsInDateHandler Handles the home page request
func (s Service) TransactionsInDateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		log.Error("Error decoding body: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error decoding body"))
		return
	}

	date, sFound := body["date"].(string)
	if !sFound {
		log.Error("Error asserting time type: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error asserting time type"))
		return
	}

	dateTime, _ := time.Parse("01/02/2006", date)
	results, err := s.U.QueryDay(dateTime)
	if err != nil {
		log.Error("Error querying day: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error querying date range"))
	}

	rStr, err := json.Marshal(results)
	if err != nil {
		log.Error("Error marshalling to json: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling to json"))
	}

	w.Write(rStr)

	return
}

// TransactionsAmountHandler Handles the home page request
func (s Service) TransactionsAmountHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		log.Error("Error decoding body: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error decoding body"))
		return
	}

	amount, sFound := body["amount"].(float64)
	if !sFound {
		log.Error("Error asserting time type: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error asserting time type"))
		return
	}

	results, err := s.U.QueryAmount(amount)
	if err != nil {
		log.Error("Error querying day: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error querying date range"))
	}

	rStr, err := json.Marshal(results)
	if err != nil {
		log.Error("Error marshalling to json: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling to json"))
	}

	w.Write(rStr)

	return
}

// TransactionsLocationHandler Handles the home page request
func (s Service) TransactionsLocationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		log.Error("Error decoding body: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error decoding body"))
		return
	}

	location, sFound := body["location"].(string)
	if !sFound {
		log.Error("Error asserting time type: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error asserting time type"))
		return
	}

	results, err := s.U.QueryLocation(location)
	if err != nil {
		log.Error("Error querying day: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error querying date range"))
	}

	rStr, err := json.Marshal(results)
	if err != nil {
		log.Error("Error marshalling to json: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling to json"))
	}

	w.Write(rStr)

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
