package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j" //Go 1.8
)

func main() {
	fmt.Println("Connecting to Neo4j")
	s, err := runQuery("bolt://localhost:7687", "Graph DBMS", "123456")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func runQuery(uri, username, password string) (string, error) {
	var (
		greeting interface{}
	)
	configForNeo4j4 := func(conf *neo4j.Config) { conf.Encrypted = false }
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), configForNeo4j4)
	if err != nil {
		return "", err
	}
	defer driver.Close()
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: "transaction"}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		return "", err
	}
	defer session.Close()

	greeting, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "3"})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}
