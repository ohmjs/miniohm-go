package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

 "github.com/ohmjs/miniohm-go"
)

func main() {
	fmt.Println("Ohm WebAssembly Matcher - Go Implementation")
	// Parse command line arguments
	wasmFile := flag.String("wasm", "test/data/_add.wasm", "Path to WebAssembly file")
	inputText := flag.String("input", "", "Input text to match against the grammar")
	inputFile := flag.String("file", "", "Path to file containing input text to match")
	verbose := flag.Bool("verbose", false, "Display verbose information about CST nodes")
	flag.Parse()

	// Create a context
	ctx := context.Background()

	// Create a new WasmMatcher
	matcher := miniohm.NewWasmMatcher(ctx)
	defer matcher.Close()

	// Load the WebAssembly module
	wasmPath := *wasmFile
	err := matcher.LoadModule(wasmPath)
	if err != nil {
		fmt.Printf("Error loading module: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded WebAssembly module: %s\n", filepath.Base(wasmPath))

	// Set the input text - either from direct text or from file
	if *inputFile != "" {
		// Read input from file
		err = matcher.SetInputFromFile(*inputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}
	} else if *inputText != "" {
		// Set the input text directly
		matcher.SetInput(*inputText)
	} else {
		fmt.Println("No input provided. Use -input flag to provide text or -file to specify an input file.")
		os.Exit(0)
	}

	// Attempt to match the input
	if *inputFile != "" {
		fmt.Printf("Matching input file: %s\n", *inputFile)
	} else {
		fmt.Printf("Matching input: %q\n", matcher.GetInput())
	}
	success, err := matcher.Match()
	if err != nil {
		fmt.Printf("Error during matching: %v\n", err)
		os.Exit(1)
	}

	if success {
		fmt.Println("Match succeeded")

		// Try to get the CST root
		node, err := matcher.GetCstRoot()
		if err != nil {
			fmt.Printf("Error getting CST root: %v\n", err)
		} else {
			nodeType := node.Type()
			if *verbose {
				fmt.Printf("CST Node - Type: %d\n", nodeType)
				// Print the rule name if available
				ruleName, err := node.RuleName()
				if err == nil {
					fmt.Printf("Rule Name: %s\n", ruleName)
				}
			}

			// Unparse the CST to get the original text
			unparsedText := unparse(node, matcher.GetInput())
			if unparsedText == matcher.GetInput() {
				fmt.Println("Unparsed text matches input")
			} else {
				fmt.Println("ERROR: Unparsed text does not match input")
				fmt.Printf("Unparsed text: %q\n", unparsedText)
				fmt.Printf("Original input: %q\n", matcher.GetInput())
			}
		}
	} else {
		fmt.Println("Match failed")
		os.Exit(1)
	}
}

// unparse walks the CST starting from the given node and reconstructs the original text
// It returns the reconstructed text from the terminal nodes
func unparse(node *miniohm.CstNode, input string) string {
	var result strings.Builder
	pos := uint32(0)
	unparseNode(node, &pos, input, &result)
	return result.String()
}

// unparseNode is a helper function that recursively processes nodes and builds the result
func unparseNode(node *miniohm.CstNode, pos *uint32, input string, result *strings.Builder) {
	// Handle terminal nodes - append the consumed text to the result
	if node.IsTerminal() {
		matchLen, err := node.MatchLength()
		if err != nil {
			fmt.Printf("Error getting match length: %v\n", err)
			return
		}

		if *pos < uint32(len(input)) && matchLen > 0 {
			end := *pos + matchLen
			if end > uint32(len(input)) {
				end = uint32(len(input))
			}
			matchedText := input[*pos:end]
			result.WriteString(matchedText)

			// Update position only after processing terminal nodes
			*pos += matchLen
		}
		return
	}

	// For all other node types (nonterminal, iteration, etc.), process children recursively
	children, err := node.Children()
	if err != nil {
		fmt.Printf("Error getting children: %v\n", err)
		return
	}

	for _, child := range children {
		unparseNode(child, pos, input, result)
	}
}
