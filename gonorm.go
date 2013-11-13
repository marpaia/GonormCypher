package gonorm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// The Gonorm type represents the relevant parameters of a Neo4j connection
type Gonorm struct {
	Host      string
	Port      int
	cypherUrl string
}

// The Cypher type represents the relevant parameters of a cypher query
type Cypher struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
	gonorm *Gonorm
}

// The Results type represents the relevant parameters of a response from the
// Neo4j API
type Results struct {
	Columns []interface{} `json:"columns"`
	Data    []interface{} `json:"data"`
	Error   error
}

// The Node type represent the relevant parameters of a Neo4j node
type Node struct {
	Params map[string]interface{}
}

// The Relationship type represent the relevant parameters of a Neo4j
// relationship
type Relationship struct {
	Params map[string]interface{}
	Type   string
	Start  string
	End    string
}

// Neo4jError defines an error type that is populated when an erroneous query is
// executed
type Neo4jError struct {
	Message    string   `json:"message"`
	Exception  string   `json:"exception"`
	Fullname   string   `json:"fullname"`
	Stacktrace []string `json:"stacktrace"`
}

// The Error method is attached to the Neo4jError type so that it satisfies the
// Error interface
func (neo4jError *Neo4jError) Error() string {
	return neo4jError.Message
}

// New is a factory method that creates a Gonorm type. You should use this
// instead of manually creating a Gonorm type yourself
func New(host string, port int) *Gonorm {
	return &Gonorm{
		Host:      host,
		Port:      port,
		cypherUrl: fmt.Sprintf("%s:%d/db/data/cypher", host, port),
	}
}

// Cypher allows you to create a Cypher type on top of a Gonorm type
func (g *Gonorm) Cypher(query string) *Cypher {
	return &Cypher{Query: query, gonorm: g}
}

// In a function-chaining fashion, On allows you to specify parameters for a
// cypher query
func (c *Cypher) On(params map[string]interface{}) *Cypher {
	c.Params = params
	return c
}

// Execute executes the query stored in the Cypher type
func (c *Cypher) Execute() *Results {
	results := new(Results)

	requestContent, err := json.Marshal(c)
	if err != nil {
		results.Error = err
		return results
	}
	buf := bytes.NewBuffer(requestContent)

	// Create a new request with the appropriate endpoints and POST body data
	req, err := http.NewRequest("POST", c.gonorm.cypherUrl, buf)
	if err != nil {
		results.Error = err
		return results
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		results.Error = err
		return results
	}
	defer resp.Body.Close()

	// Given the response body, which is an io.Reader, return a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		results.Error = err
		return results
	}

	if resp.StatusCode != http.StatusOK {
		neo4jError := new(Neo4jError)
		json.Unmarshal(body, &neo4jError)
		results.Error = neo4jError
		return results
	}

	// Return the results
	json.Unmarshal(body, &results)
	return results
}

// AsInt should be used when you return one integer in your cypher query
func (r *Results) AsInt() (int, error) {
	if r.Error != nil {
		return -1, r.Error
	}

	return int(r.Data[0].([]interface{})[0].(float64)), nil
}

// AsInts should be used when you return several integers in your cypher query
func (r *Results) AsInts() ([]int, error) {
	if r.Error != nil {
		return []int{}, r.Error
	}

	var ints = []int{}

	for _, value := range r.Data[0].([]interface{}) {
		i := int(value.(float64))
		ints = append(ints, i)
	}

	return ints, nil
}

// AsString should be used when you return one string in your cypher query
func (r *Results) AsString() (string, error) {
	if r.Error != nil {
		return "", r.Error
	}

	return r.Data[0].([]interface{})[0].(string), nil
}

// AsStrings should be used when you return several strings in your cypher query
func (r *Results) AsStrings() ([]string, error) {
	if r.Error != nil {
		return []string{}, r.Error
	}

	var strings = []string{}

	for _, value := range r.Data[0].([]interface{}) {
		strings = append(strings, value.(string))
	}

	return strings, nil
}

// AsNode should be used when you return one node in your cypher query
func (r *Results) AsNode() (*Node, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	nodeData := r.Data[0].([]interface{})[0].(map[string]interface{})
	node := new(Node)
	node.Params = nodeData["data"].(map[string]interface{})

	return node, nil
}

// AsNodes should be used when you return several nodes in your cypher query
func (r *Results) AsNodes() ([]*Node, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	var nodes = []*Node{}

	for _, value := range r.Data[0].([]interface{}) {
		node := new(Node)
		node.Params = value.(map[string]interface{})["data"].(map[string]interface{})
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AsRelationship should be used when you return one relationship in your cypher
// query
func (r *Results) AsRelationship() (*Relationship, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	relationshipData := r.Data[0].([]interface{})[0].(map[string]interface{})
	relationship := new(Relationship)
	relationship.Params = relationshipData["data"].(map[string]interface{})
	relationship.Type = relationshipData["type"].(string)
	relationship.Start = relationshipData["start"].(string)
	relationship.End = relationshipData["end"].(string)

	return relationship, nil
}

// AsRelationships should be used when you return several relationships in your
// cypher query
func (r *Results) AsRelationships() ([]*Relationship, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	var relationships = []*Relationship{}

	for _, value := range r.Data[0].([]interface{}) {
		relationship := new(Relationship)
		relationship.Params = value.(map[string]interface{})["data"].(map[string]interface{})
		relationship.Type = value.(map[string]interface{})["type"].(string)
		relationship.Start = value.(map[string]interface{})["start"].(string)
		relationship.End = value.(map[string]interface{})["end"].(string)
		relationships = append(relationships, relationship)
	}

	return relationships, nil
}
