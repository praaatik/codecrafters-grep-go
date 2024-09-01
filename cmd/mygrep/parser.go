package main

// type TokenType int
//
// const (
// 	Char TokenType = iota
// 	DigitClass
// )
//
// type Token struct {
// 	Type  TokenType
// 	Value string
// }
//
// type Parser struct {
// 	pattern  string
// 	position int
// 	tokens   []Token
// }
//
// func (p *Parser) parseCharacterClass() Token {
// 	switch p.pattern[p.position] {
// 	case 'd':
// 		p.position += 1
// 		return Token{
// 			Type:  DigitClass,
// 			Value: "\\d",
// 		}
// 	}
// 	return Token{}
// }
//
// func (p *Parser) Parse() []Token {
// 	for p.position < len(p.pattern) {
// 		char := p.pattern[p.position]
//
// 		switch char {
// 		case '\\':
// 			p.position += 1
// 			token := p.parseCharacterClass()
// 			p.tokens = append(p.tokens, token)
// 		}
// 	}
// 	return p.tokens
// }
//
// func (p *Parser) Match(input string) bool {
// 	inputLength := len(input)
//
// 	// Iterate over each possible starting position in the input
// 	for start := 0; start < inputLength; start++ {
// 		inputPos := start
// 		matched := true
//
// 		// Attempt to match the pattern from this starting position
// 		for _, token := range p.tokens {
// 			switch token.Type {
// 			case DigitClass:
// 				if inputPos >= inputLength || !isDigit(input[inputPos]) {
// 					matched = false
// 					break
// 				}
// 				inputPos++
// 			// case Char:
// 			// 	if inputPos >= inputLength || input[inputPos] != token.Value[0] {
// 			// 		matched = false
// 			// 		break
// 			// 	}
// 			// 	inputPos++
//
// 			default:
// 				matched = false
// 			}
//
// 			if !matched {
// 				break // Stop checking if this substring doesn't match
// 			}
// 		}
//
// 		// If all tokens matched successfully, return true
// 		if matched {
// 			return true
// 		}
// 	}
//
// 	// No matching substring found
// 	return false
// }
//
// func isDigit(c byte) bool {
// 	return c > '0' && c <= '9'
// }
