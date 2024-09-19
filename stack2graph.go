package stacktracetograph

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var GLOBAL_STACK_TO_GRAPH *StackToGraph

type StackToGraph struct {
	neo4jDriver neo4j.Driver
	sync.Mutex
	cacheReportedStacks map[string]bool
}

func NewStackToGraph(uri, username, password string) (*StackToGraph, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	return &StackToGraph{
		neo4jDriver: driver,
	}, nil
}

func (s *StackToGraph) SetupGlobal() {
	GLOBAL_STACK_TO_GRAPH = s
}

func (s *StackToGraph) ReportStacktrace() error {
	// Capture the stack trace
	stack := captureStackTrace()

	// Lock the cache to avoid race
	s.Lock()
	if s.cacheReportedStacks == nil {
		s.cacheReportedStacks = make(map[string]bool)
	}
	if _, ok := s.cacheReportedStacks[stack]; ok {
		s.Unlock()
		// Skip reporting the same stack trace
		return nil
	}
	s.Unlock()

	// Parse the stack trace to extract function calls
	parsedStack := parseStackTrace(stack)

	// Report the stack trace to Neo4j
	err := s.reportStackTraceToNeo4j(parsedStack)
	if err != nil {
		log.Printf("Error reporting stack trace to Neo4j: %v\n", err)
		return err
	}
	s.Lock()
	s.cacheReportedStacks[stack] = true
	s.Unlock()

	return nil
}

func (s *StackToGraph) reportStackTraceToNeo4j(stackTraceData []ParsedStackEntry) error {
	// Create a new session
	session := s.neo4jDriver.NewSession(neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close()

	// Execute a write transaction
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		var previousNodeID interface{}

		// Reverse the stack to represent the top-down call flow
		for i := len(stackTraceData) - 1; i >= 0; i-- {
			frame := stackTraceData[i]

			// Merge node for each function call to avoid duplicates
			result, err := tx.Run(`
			    MERGE (f:Function {name: $name, package: $package})
			    SET f.file = $file, f.line = $line
			    RETURN id(f) AS nodeID
			`, map[string]interface{}{
				"name":    frame.Function,
				"file":    frame.File,
				"line":    frame.Line,
				"package": frame.Package,
			})
			if err != nil {
				return nil, err
			}

			record, err := result.Single()
			if err != nil {
				return nil, err
			}
			currentNodeID := record.Values[0]

			// Create "CALLS" relationship from the previous node to the current node
			if previousNodeID != nil {
				_, err = tx.Run(`
					MATCH (caller), (callee)
					WHERE id(caller) = $callerID AND id(callee) = $calleeID
					MERGE (caller)-[:CALLS]->(callee)
				`, map[string]interface{}{
					"callerID": previousNodeID,
					"calleeID": currentNodeID,
				})
				if err != nil {
					return nil, err
				}
			}

			previousNodeID = currentNodeID
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute write transaction: %w", err)
	}

	return nil
}

// Close the Neo4j driver when the application exits
func (s *StackToGraph) Close() {
	if s.neo4jDriver != nil {
		s.neo4jDriver.Close()
	}
}

// ReportStacktrace encapsulates capturing, parsing, and reporting the stack trace to Neo4j.
func ReportStacktrace() error {
	if GLOBAL_STACK_TO_GRAPH == nil {
		return fmt.Errorf("global StackToGraph instance is not set. Call SetupGlobal before using ReportStacktrace.")
	}
	return GLOBAL_STACK_TO_GRAPH.ReportStacktrace()
}

// captureStackTrace captures the current call stack as a string.
func captureStackTrace() string {
	// Adjust buffer size if needed (1<<16 is 64KB, which is usually sufficient)
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, false)
	return string(buf[:stackSize])
}

// parseStackTrace extracts function names, file paths, and line numbers from the stack trace.
// It also cleans up function names by removing arguments.
func parseStackTrace(stackTrace string) []ParsedStackEntry {
	// Regular expression to capture function calls in the stack trace
	// Example stack trace lines:
	// main.isError({0x104486073?, 0xc?})
	//     /path/to/file/main.go:24 +0x9f
	re := regexp.MustCompile(`(?m)^(.*?)\n\s+(.*?)\:(\d+)(?: \+0x[0-9a-f]+)?$`)

	matches := re.FindAllStringSubmatch(stackTrace, -1)
	var parsedData []ParsedStackEntry

	for _, match := range matches {
		functionWithArgs := strings.TrimSpace(match[1])
		// Remove arguments from the function name using regex
		cleanFunction := cleanFunctionName(functionWithArgs)
		pkg := ParsePackageName(cleanFunction)
		cleanFunction = strings.TrimPrefix(cleanFunction, pkg+".")
		cleanFunction = strings.TrimPrefix(cleanFunction, pkg)

		parsedData = append(parsedData, ParsedStackEntry{
			Function: cleanFunction,
			File:     match[2],
			Line:     match[3],
			Package:  pkg,
		})
	}

	return parsedData
}

func cleanFunctionName(s string) string {
	// Remove any trailing whitespace
	s = strings.TrimRight(s, " \t\n\r")

	// Start scanning from the end of the string backwards
	parenLevel := 0
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c == ')' {
			parenLevel++
		} else if c == '(' {
			parenLevel--
			if parenLevel == 0 {
				// We have found the outermost '('
				// Remove everything from position i onwards
				return s[:i]
			}
		}
	}
	// If no outermost '(', return the original string
	return s
}

// ParsePackageName extracts the package path from a fully qualified method/function string.
// It returns the package path up to the first '.' after the last '/'.
// If there is no '/', it takes the string up to the first '.'.
func ParsePackageName(s string) string {
	// Find the index of the last '/'
	lastSlash := strings.LastIndex(s, "/")
	if lastSlash == -1 {
		lastSlash = 0
	} else {
		lastSlash += 1 // Move past the '/'
	}

	// Find the index of the first '.' after the last '/'
	dotIndex := strings.Index(s[lastSlash:], ".")
	if dotIndex == -1 {
		// No '.' found after the last '/', return the entire string
		return s
	}

	// The package path is up to lastSlash + dotIndex
	return s[:lastSlash+dotIndex]
}

type ParsedStackEntry struct {
	Function string
	File     string
	Line     string
	Package  string
}
