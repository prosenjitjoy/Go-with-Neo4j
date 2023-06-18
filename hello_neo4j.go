package main

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext("neo4j://127.0.0.1:7687", neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		log.Fatal("error connecting to database")
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	fmt.Println("Successfully connected to database 🎉")

	result, err := session.Run(
		ctx, // (1)
		`
	MATCH (p:Person)-[:DIRECTED]->(:Movie {title: $title})
	RETURN p
	`, // (2)
		map[string]any{"title": "The Matrix"},
	)

	people, err := neo4j.CollectTWithContext(ctx, result,
		func(record *neo4j.Record) (neo4j.Node, error) {
			person, _, err := neo4j.GetRecordValue[neo4j.Node](record, "p")
			return person, err
		})
	if err != nil {
		panic(err)
	}
	fmt.Println(people)
}
