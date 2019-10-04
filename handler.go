package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type requestHandler struct {
	client *NeoClient
}
type Th struct{
	Text string `json:"text"`
	Id string `json:"id"`
}
type payload struct {
	Period *struct {
		Start int64 `json:"startDate"`
		End   int64 `json:"endDate"`
	} `json:"period,omitempty"`
	IsRelatedWith         []Th `json:"isRelatedWith"`
	IsDirectlyRelatedWith []Th `json:"isDirectlyRelatedWith"`
}

type responseBody struct {
	UUIDs []string `json:"uuid"`
	Error string   `json:"error,omitempty"`
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
	writer.Header().Set("Content-Type", "application/json")
	//writer.WriteHeader(http.StatusOK)
	//writer.Header().Set("Content-Type", "application/json")
	writer.Write(body)
}

func payloadToSearchObject(data payload) *SearchObject {
	fmt.Printf("%+v",data)
	var (
		start int64 = 1538582036 // last year this time
		end   int64 = time.Now().Unix()
	)
	//if data.Period != nil {
	//	start = data.Period.Start
	//	end = data.Period.End
	//}
	//start = data.Period.Start
	//end = data.Period.End
	s := []string{}
	for _,v := range(data.IsRelatedWith){
		s = append(s, v.Id)
	}
	s1 := []string{}
	for _,v := range(data.IsDirectlyRelatedWith){
		s1 = append(s1, v.Id)
	}
	return &SearchObject{
		fromDate:              start,
		toDate:                end,
		isRelatedWith:         s,
		isDirectlyRelatedWith: s1,
		limit:                 25,
	}
}

func writeError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	fmt.Println("--------")
	fmt.Println(err)
	fmt.Println("--------")
	body, _ := json.Marshal(responseBody{Error: err.Error()})
	writer.Write(body)
}
