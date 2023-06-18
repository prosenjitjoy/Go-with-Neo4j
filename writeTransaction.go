package main

// Import the driver
import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		credentials.Uri,
		neo4j.BasicAuth(credentials.Username, credentials.Password, ""),
	)
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	cypher := `
		MATCH (m:Movie {title: $title})
		CREATE (p:Person {name: $name})
		CREATE (p)-[:ACTED_IN]->(m)
		RETURN p`
	params := map[string]any{
		"title": "Matrix, The",
		"name":  "Prosenjit Joy",
	}

	personNode, err := neo4j.ExecuteWrite(ctx, session,
		func(tx neo4j.ManagedTransaction) (neo4j.Node, error) {
			result, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return *new(neo4j.Node), err
			}

			return neo4j.SingleTWithContext(ctx, result,
				func(record *neo4j.Record) (neo4j.Node, error) {
					node, _, err := neo4j.GetRecordValue[neo4j.Node](record, "p")
					return node, err
				})
		})
	if err != nil {
		panic(err)
	}
	fmt.Println(personNode)
}
