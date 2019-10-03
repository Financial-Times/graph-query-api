package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type requestHandler struct {
	client *NeoClient
}

type payload struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	IsRelatedWith         []string `json:"isRelatedWith"`
	IsDirectlyRelatedWith []string `json:"isDirectlyRelatedWith"`
}

type responseBody struct {
	UUIDs []string `json:"uuid"`
	Error string   `json:"error"`
}

func (handler *requestHandler) searchEndpoint(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writeError(writer, err)
		return
	}

	var p payload
	err = json.Unmarshal(body, &p)
	if err != nil {
		writeError(writer, err)
		return
	}

	if handler.client != nil {
		writeError(writer, errors.New("failed request, NeoClient is not initialized"))
		return
	}

	content, err := handler.client.Search(payloadToSearchObject(p))
	if err != nil {
		writeError(writer, err)
		return
	}

	body, err = json.Marshal(responseBody{UUIDs: content})

	if err != nil {
		writeError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}

func payloadToSearchObject(data payload) *SearchObject {
	return &SearchObject{
		fromDate:              data.Period.Start.Unix(),
		toDate:                data.Period.End.Unix(),
		isRelatedWith:         data.IsRelatedWith,
		isDirectlyRelatedWith: data.IsDirectlyRelatedWith,
	}
}

func writeError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	body, _ := json.Marshal(responseBody{Error: err.Error()})
	writer.Write(body)
}
