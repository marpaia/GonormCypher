package gonorm

import (
	"testing"
)

var g *Gonorm

func init() {
	g = New("http://localhost", 7474)
}

func TestAsString(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[:KNOWS]->p2
	RETURN p1.name
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsString()

	if err != nil {
		t.Error(err)
	}

	if results != "Mike" {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsStrings(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[:KNOWS]->p2
	RETURN p1.name, p2.name
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsStrings()

	if err != nil {
		t.Error(err)
	}

	if results[0] != "Mike" || results[1] != "Matt" {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsInt(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[:KNOWS]->p2
	RETURN id(p1)
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsInt()

	if err != nil {
		t.Error(err)
	}

	if results <= 0 {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsInts(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[:KNOWS]->p2
	RETURN id(p1), id(p2)
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsInts()

	if err != nil {
		t.Error(err)
	}

	if results[0] <= 0 || results[1] <= 0 {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsNode(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[:KNOWS]->p2
	RETURN p1
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsNode()

	if err != nil {
		t.Error(err)
	}

	if results.Params["name"] != "Mike" {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsNodes(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[:KNOWS]->p2
	RETURN p1, p2
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsNodes()

	if err != nil {
		t.Error(err)
	}

	if results[0].Params["name"] != "Mike" {
		t.Error("Results not what was expected:", results)
	}

	if results[1].Params["name"] != "Matt" {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsRelationship(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[k:KNOWS]->p2
	RETURN k
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsRelationship()

	if err != nil {
		t.Error(err)
	}

	if results.Type != "KNOWS" {
		t.Error("Results not what was expected:", results)
	}
}

func TestAsRelationships(t *testing.T) {
	results, err := g.Cypher(`
	MERGE (p1:Person{name:{name1}})
	MERGE (p2:Person{name:{name2}})
	CREATE UNIQUE p1-[k1:KNOWS]->p2
	CREATE UNIQUE p2-[k2:KNOWS]->p1
	RETURN k1, k2
	`).On(map[string]interface{}{
		"name1": "Mike",
		"name2": "Matt",
	}).Execute().AsRelationships()

	if err != nil {
		t.Error(err)
	}

	for _, r := range results {
		if r.Type != "KNOWS" {
			t.Error("Results not what was expected:", results)
		}
	}
}
