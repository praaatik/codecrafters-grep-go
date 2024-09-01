package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	parser := Parser{
		pattern:  pattern,
		position: 0,
		tokens:   []Token{},
	}

	parser.Parse()

	if parser.Match(string(line)) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
	// return false, fmt.Errorf("input does not match")

	//ok, err := matchLine(line, pattern)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	//	os.Exit(1)
	//}

	//if !ok {
	//	os.Exit(1)
	//}

	//// default exit code is 0 which means success
	//os.Exit(0)
}

func matchLine(line []byte, pattern string) (bool, error) {
	parser := Parser{
		pattern:  pattern,
		position: 0,
		tokens:   []Token{},
	}

	parser.Parse()

	if parser.Match(string(line)) {
	}
	return false, fmt.Errorf("input does not match")
}

type TokenType int

const (
	Char TokenType = iota
	DigitClass
)

type Token struct {
	Type  TokenType
	Value string
}

type Parser struct {
	pattern  string
	position int
	tokens   []Token
}

func (p *Parser) parseCharacterClass() Token {
	fmt.Println(p)
	switch p.pattern[p.position] {
	case 'd':
		p.position += 1
		return Token{
			Type:  DigitClass,
			Value: "\\d",
		}
	default:
	}
	return Token{}
}

func (p *Parser) Parse() []Token {
	for p.position < len(p.pattern) {
		char := p.pattern[p.position]

		switch char {
		case '\\':
			p.position += 1
			token := p.parseCharacterClass()
			p.tokens = append(p.tokens, token)
		default:
			p.position += 1
			token := Token{
				Type:  Char,
				Value: string(char),
			}
			p.tokens = append(p.tokens, token)
		}
	}

	return p.tokens
}

func (p *Parser) Match(input string) bool {
	inputLength := len(input)

	// Iterate over each possible starting position in the input
	for start := 0; start < inputLength; start++ {
		inputPos := start
		matched := true

		// Attempt to match the pattern from this starting position
		for _, token := range p.tokens {
			switch token.Type {
			case DigitClass:
				if inputPos >= inputLength || !isDigit(input[inputPos]) {
					matched = false
					break
				}
				inputPos++
			case Char:
				if inputPos >= inputLength || input[inputPos] != token.Value[0] {
					matched = false
					break
				}
				inputPos++

			default:
				matched = false
			}

			if !matched {
				break // Stop checking if this substring doesn't match
			}
		}

		// If all tokens matched successfully, return true
		if matched {
			return true
		}
	}

	// No matching substring found
	return false
}

func isDigit(c byte) bool {
	return c > '0' && c <= '9'
}
