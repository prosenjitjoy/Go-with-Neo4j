package main

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()

	cypher := `
      MATCH (p:Person)-[:DIRECTED]->(:Movie {title: $title})
      RETURN p.name AS Director
    `
	params := map[string]any{"title": "Toy Story"}

	driver, err := neo4j.NewDriverWithContext(credentials.Uri, neo4j.BasicAuth(credentials.Username, credentials.Password, ""))
	if err != nil {
		panic(err)
	}

	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		panic(err)
	}

	director, err := neo4j.SingleTWithContext(ctx, result,
		func(record *neo4j.Record) (string, error) {
			director, _, err := neo4j.GetRecordValue[string](record, "Director")
			return director, err
		})
	if err != nil {
		panic(err)
	}
	log.Println(director)
}
