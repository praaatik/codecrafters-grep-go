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

	//uncomment in case of debugging
	/*
		pattern := `ca?t`
		line := "act"
	*/

	parser := Parser{
		pattern:  pattern,
		position: 0,
		tokens:   []Token{},
	}
	parser.Parse()

	if parser.Match(string(line)) {
		fmt.Println("substring -> ", parser.substring)
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

type TokenType int

const (
	Char TokenType = iota

	DigitCharacterClass
	AlphanumericCharacterClass

	PositiveCharacterGroup
	NegativeCharacterGroup

	StartAnchor
	EndAnchor

	OneOrMoreQuantifier
	ZeroOrOneQuantifier
)

type Token struct {
	Type  TokenType
	Value string
	Chars []rune
}

type Parser struct {
	pattern   string
	position  int
	tokens    []Token
	substring string
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

	case '\\':
		p.position += 1
		return Token{
			Type:  Char,
			Value: `\`,
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

		case '^':
			p.position += 1
			token := Token{
				Type:  StartAnchor,
				Value: "^",
			}
			p.tokens = append(p.tokens, token)

		case '$':
			p.position += 1
			token := Token{
				Type:  EndAnchor,
				Value: "$",
			}
			p.tokens = append(p.tokens, token)

		case '+':
			p.position += 1
			if len(p.tokens) > 0 {
				lastParsedToken := p.tokens[len(p.tokens)-1]
				token := Token{
					Type:  OneOrMoreQuantifier,
					Value: lastParsedToken.Value + "+",
				}
				p.tokens = p.tokens[:len(p.tokens)-1]
				p.tokens = append(p.tokens, token)
			}

		case '?':
			p.position += 1
			if len(p.tokens) > 0 {
				lastParsedToken := p.tokens[len(p.tokens)-1]
				token := Token{
					Type:  ZeroOrOneQuantifier,
					Value: lastParsedToken.Value + "?",
				}
				p.tokens = p.tokens[:len(p.tokens)-1]
				p.tokens = append(p.tokens, token)
				//p.tokens = append(p.tokens, token)
			}

		default:
			//fmt.Println("default condition true")
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

	for start := 0; start < inputLength; start++ {
		inputPos := start
		matched := true

		for i := 0; i < len(p.tokens); i++ {
			token := p.tokens[i]
			switch token.Type {
			case StartAnchor:
				if inputPos != 0 {
					matched = false
					break
				}

			case EndAnchor:
				if inputPos != inputLength {
					matched = false
					break
				}

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

			case OneOrMoreQuantifier:
				count := 0

				char := p.tokens[i].Value

				for inputPos < inputLength && input[inputPos] == char[0] {
					inputPos++
					count++
				}

				if count < 1 {
					matched = false
					break
				}

			case ZeroOrOneQuantifier:
				// Match the preceding character or group zero or one time
				count := 0
				char := p.tokens[i].Value

				for inputPos < inputLength && input[inputPos] == char[0] {
					inputPos++
					count++
				}

				if count > 1 {
					matched = false
					break
				}

				//if i == 0 {
				//	matched = false
				//	break
				//}
				//precedingToken := p.tokens[i-1]
				//if inputPos < inputLength && matchToken(precedingToken, input[inputPos]) {
				//	inputPos++
				//}

			default:
				matched = false
			}

			if !matched {
				p.substring = ""
				break // Stop checking if this substring doesn't match
			}
			if inputPos > 0 {
				p.substring += string(input[inputPos-1])
			}
		}

		// If all tokens matched successfully, return true
		if matched {
			return true
		}
	}

	// No matching substring found
	p.substring = ""
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

func matchToken(token Token, c byte) bool {
	//char := token.Value[0]
	//fmt.Println(char)

	switch token.Type {
	case Char:
		return c == token.Value[0]
	case AlphanumericCharacterClass:
		return isAlphanumeric(c)
	case DigitCharacterClass:
		return isDigit(c)
	case PositiveCharacterGroup:
		return isCharacterInGroup(c, token.Chars)
	case NegativeCharacterGroup:
		return !isCharacterInGroup(c, token.Chars)
	default:
		return false
	}
}

func removePreviousToken(tokenList []Token) []Token {
	if len(tokenList) == 0 {
		return tokenList
	}
	return tokenList[:len(tokenList)-1]
}
