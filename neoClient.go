package main

import (
	"errors"
	"fmt"

	"github.com/Financial-Times/go-logger"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type NeoClient struct {
	driver neo4j.Driver
}

func NewNeoClient(neoURL string) (*NeoClient, error) {
	neoDriver, err := neo4j.NewDriver(neoURL, neo4j.NoAuth(), func(config *neo4j.Config) {
		config.Log = neo4j.ConsoleLogger(neo4j.INFO)
	})
	if err != nil {
		return nil, err
	}
	return &NeoClient{neoDriver}, nil
}

var relatedPattern = `MATCH (:Concept{uuid:"%s"})-[:EQUIVALENT_TO]->(%s:Concept) MATCH (%s)<-[:EQUIVALENT_TO]-(:Concept)-[:HAS_BROADER*1..]->(%s:Concept) `
var timePattern = "WHERE c.publishedDateEpoch>%d AND c.publishedDateEpoch<%d "
var directlyRelatedPattern = "AND (c)-[:ABOUT|MENTIONS]->(:Concept{uuid:\"%s\"}) "

func constructStatement(sObj *SearchObject) string {

	var statement string

	for i, uuid := range sObj.isRelatedWith {
		canonLabel := fmt.Sprintf("canon%d", i)
		implicitLabel := fmt.Sprintf("implicit%d", i)
		related := fmt.Sprintf(relatedPattern, uuid, canonLabel, canonLabel, implicitLabel)
		statement += related
	}
	collectStatement := "WITH collect(implicit0)"
	for idx := 1; idx < len(sObj.isRelatedWith); idx++ {
		collectStatement += fmt.Sprintf(" + collect(implicit%d)", idx)
	}
	statement += fmt.Sprintf("%s as col UNWIND col as implicit MATCH (implicit)<-[]-(c:Content) ", collectStatement)

	statement += fmt.Sprintf(timePattern, sObj.fromDate, sObj.toDate)

	for _, uuid := range sObj.isDirectlyRelatedWith {
		related := fmt.Sprintf(directlyRelatedPattern, uuid)
		statement += related
	}

	statement += fmt.Sprintf("RETURN c.uuid as uuid LIMIT(%d) ", sObj.limit)

	fmt.Println(statement)
	return statement
}

type SearchObject struct {
	fromDate              int64
	toDate                int64
	isRelatedWith         []string
	isDirectlyRelatedWith []string
	limit                 int
}

func (nc *NeoClient) Search(sObj *SearchObject) ([]string, error) {
	var session neo4j.Session
	var records []neo4j.Record
	var err error

	statement := constructStatement(sObj)
	if session, err = nc.driver.Session(neo4j.AccessModeRead); err != nil {
		return nil, err
	}
	defer session.Close()

	records, err = neo4j.Collect(session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(statement, map[string]interface{}{"contentUUID": "", "annotationLifecycle": ""})
	}))

	if err != nil {
		return nil, err
	}

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

		results = append(results, fmt.Sprintf("%v", recordUUID))
	}

	return results, nil
}

func (nc *NeoClient) Close() error {
	return nc.driver.Close()
}
