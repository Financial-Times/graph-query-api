package main

import (
	"errors"
	"fmt"

	"github.com/Financial-Times/go-logger"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type NeoClient struct{
	driver neo4j.Driver
}

func NewNeoClient(neoURL string) (*NeoClient, error){
	neoDriver, err := neo4j.NewDriver(neoURL, neo4j.NoAuth(), func(config *neo4j.Config) {
		config.Log = neo4j.ConsoleLogger(neo4j.INFO)
	})
	if err != nil {
		return nil, err
	}
	return &NeoClient{neoDriver}, nil
}

type SearchObject struct{
	fromDate int
	toDate int
	isRelatedWith []string
	isDirectlyRelatedWith []string
	limit int
}

func (nc *NeoClient) Search(sObj *SearchObject) ([]string, error){
	var session neo4j.Session
	var records []neo4j.Record
	var err error

	statement := `
		MATCH (:Concept{uuid:"9ab8e36c-4b79-4e96-9aae-cc586b7d19c4"})-[:EQUIVALENT_TO]->(canon:Concept)
		MATCH (canon)<-[:EQUIVALENT_TO]-(leaves:Concept)-[:HAS_BROADER*1..]->(implicit:Concept)<-[]-(c:Content)
		WITH c, leaves, implicit
		WHERE c.publishedDateEpoch>1569179507 AND c.publishedDateEpoch<1570043512 AND c-[:ABOUT]->(:Topic{prefLabel:"Agriculture"})
		RETURN c.uuid as uuid
		LIMIT(5)`
	if session, err = nc.driver.Session(neo4j.AccessModeRead); err != nil {
		return nil, err
	}
	defer session.Close()

	records, err = neo4j.Collect(session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(statement, map[string]interface{}{"contentUUID": "", "annotationLifecycle": ""})
	}))

	logger.Infof("Query returned following results: %v", records)

	if len(records) < 1 {
		return nil, errors.New("no results found")
	}

	var results []string
	for _, record := range records {
		fmt.Printf("%+v\n", record)
		recordUUID, ok := record.Get("uuid")
		if !ok {
			logger.Error("not found uuid for record")
		}

		results = append(results, fmt.Sprintf("%v",recordUUID))
	}

	return results, nil
}

func (nc *NeoClient) Close()error{
	return nc.driver.Close()
}