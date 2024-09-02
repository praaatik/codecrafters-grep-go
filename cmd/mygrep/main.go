package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(1) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(1)
	}

	parser := Parser{
		pattern:  pattern,
		position: 0,
		tokens:   []Token{},
	}

	parser.Parse()

	if parser.Match(string(line)) {
		fmt.Println("match found")
		os.Exit(0)
	}
	fmt.Println("match not found")
	os.Exit(1)
}

type TokenType int

const (
	Char TokenType = iota

	DigitCharacterClass
	AlphanumericCharacterClass

	PositiveCharacterGroup
	NegativeCharacterGroup
)

type Token struct {
	Type  TokenType
	Value string
	Chars []rune
}

type Parser struct {
	pattern  string
	position int
	tokens   []Token
}

func (p *Parser) parseCharacterGroup() Token {
	startPosition := p.position

	isNegativeCharacterGroup := false
	if p.pattern[p.position] == '^' {
		isNegativeCharacterGroup = true
		p.position += 1
	}
	var chars []rune

	for p.position < len(p.pattern) && p.pattern[p.position] != ']' {
		chars = append(chars, rune(p.pattern[p.position]))
		p.position += 1
	}

	p.position += 1
	if isNegativeCharacterGroup {
		return Token{
			Type:  NegativeCharacterGroup,
			Value: p.pattern[startPosition:p.position],
			Chars: chars,
		}

	}
	return Token{
		Type:  PositiveCharacterGroup,
		Value: p.pattern[startPosition:p.position],
		Chars: chars,
	}
}

func (p *Parser) parseCharacterClass() Token {
	fmt.Println(p)
	switch p.pattern[p.position] {
	case 'd':
		p.position += 1
		return Token{
			Type:  DigitCharacterClass,
			Value: "\\d",
		}

	case 'w':
		p.position += 1
		return Token{
			Type:  AlphanumericCharacterClass,
			Value: "\\w",
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

		case '[':
			p.position += 1
			token := p.parseCharacterGroup()
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
			case AlphanumericCharacterClass:
				if inputPos >= inputLength || !isAlphanumeric(input[inputPos]) {
					matched = false
					break
				}
				inputPos++

			case DigitCharacterClass:
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

			case PositiveCharacterGroup:
				if !isCharacterInGroup(input[inputPos], token.Chars) {
					matched = false
					break
				}
				inputPos++

			case NegativeCharacterGroup:
				if isCharacterInGroup(input[inputPos], token.Chars) {
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

func isAlphanumeric(c byte) bool {
	return unicode.IsLetter(rune(c)) || c == '_' || isDigit(c)
}

func isCharacterInGroup(c byte, group []rune) bool {
	for _, value := range group {
		if c == byte(value) {
			return true

		}
	}
	return false
}
