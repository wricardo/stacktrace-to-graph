package main

import (
	"fmt"
	"log"

	stacktracetograph "github.com/wricardo/stacktrace-to-graph"
	"github.com/wricardo/stacktrace-to-graph/example/subdir"
)

type Person struct {
	Name string
}

func (p *Person) SayHello() {
	stacktracetograph.ReportStacktrace()
	fmt.Printf("Hello, my name is %s\n", p.Name)
}

func functionC() {
	subdir.DoSomething()
	stacktracetograph.ReportStacktrace()
}

// Sample functions to simulate application calls
func functionA() {
	functionB()
}

func main() {
	s2g, err := stacktracetograph.NewStackToGraph("neo4j://localhost", "neo4j", "wallace123")
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j driver: %v", err)
	}
	s2g.SetupGlobal()
	defer s2g.Close()

	// Simulate application flow
	functionA()

	fmt.Println("Application finished.")
}
