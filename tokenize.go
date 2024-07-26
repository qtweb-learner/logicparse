package main

import "unicode"

func popLastChar(runes *[]rune, c rune) []rune {
	var results []rune
	for {
		length := len(*runes)
		if length == 0 {
			break
		}
		if (*runes)[length-1] != c {
			break
		}
		results = append(results, c)
		(*runes) = (*runes)[0 : length-1]
	}
	return results
}

// tokenize 将查询字符串分词成符号列表
func tokenize(query string) []string {
	runeQuery := []rune(query)

	var tokens []string
	var currentToken []rune

	var quoto rune
	quotoMode := false

	for i := 0; i < len(runeQuery); i++ {
		char := runeQuery[i]
		nextChar := rune(0)
		if i+1 < len(runeQuery) {
			nextChar = runeQuery[i+1]
		}
		if char == '\'' || char == '"' {
			if !quotoMode {
				if len(currentToken) != 0 {
					tokens = append(tokens, string(currentToken))
					currentToken = currentToken[0:0]
				}
				quotoMode = true
				quoto = char
			} else {
				if char != quoto {
					currentToken = append(currentToken, char)
				} else {
					length := len(popLastChar(&currentToken, '\\'))
					for j := 0; j < length/2; j++ {
						currentToken = append(currentToken, '\\')
					}
					if length%2 == 0 {
						quotoMode = false
						tokens = append(tokens, string(currentToken))
						currentToken = currentToken[0:0]
					} else {
						currentToken = append(currentToken, char)
					}
				}
			}
			continue
		}
		if quotoMode {
			currentToken = append(currentToken, char)
			continue
		}
		if char == '>' && nextChar == '=' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, ">=")
			i++
			continue
		}
		if char == '>' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, ">")
			continue
		}
		if char == '<' && nextChar == '=' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, "<=")
			i++
			continue
		}
		if char == '<' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, "<")
			continue
		}

		if char == '=' && nextChar == '=' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, "==")
			i++
			continue
		}
		if char == '=' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, "=")
			continue
		}

		if char == '&' && nextChar == '&' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, "&&")
			i++
			continue
		}
		if char == '|' && nextChar == '|' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, "||")
			i++
			continue
		}
		if char == '(' || char == ')' {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			tokens = append(tokens, string(char))
			continue
		}
		if unicode.IsSpace(char) {
			if len(currentToken) != 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = currentToken[0:0]
			}
			continue
		}
		currentToken = append(currentToken, char)
	}

	if len(currentToken) != 0 {
		tokens = append(tokens, string(currentToken))
	}

	return tokens
}
