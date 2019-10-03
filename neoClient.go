package main

import "github.com/neo4j/neo4j-go-driver/neo4j"

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
	return nil, nil
}

func (nc *NeoClient) Close()error{
	return nc.driver.Close()
}