package main

import "fmt"

// Parser 解析器结构体
type Parser struct {
	tokens  []string
	current int
}

// NewParser 创建一个新的解析器
func NewParser(query string) *Parser {
	// 分词，将查询字符串转换成符号列表
	tokens := tokenize(query)
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse 解析查询字符串，返回根节点
func (p *Parser) Parse() (*Node, error) {
	return p.parseExpression(false)
}

// parseExpression 解析表达式
func (p *Parser) parseExpression(sub bool) (*Node, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		token, ok := p.consume()
		if !ok {
			if !sub {
				return node, nil
			} else {
				return nil, fmt.Errorf("expected ) to close sub expression but encounter end")
			}
		}
		if token != "and" && token != "or" && token != "&&" && token != "||" && token != ")" {
			if sub {
				return nil, fmt.Errorf("expected and or && || ) but encounter %s", token)
			} else {
				return nil, fmt.Errorf("expected and or && || but encounter %s", token)
			}
		}
		if token == ")" {
			if !sub {
				return nil, fmt.Errorf("expected and or && || but encounter %s", token)
			} else {
				return node, nil
			}
		}
		logicOper := token
		if logicOper == "||" {
			logicOper = "or"
		}
		if logicOper == "&&" {
			logicOper = "and"
		}
		right, err := p.parseTerm() // 解析右边的项
		if err != nil {
			return nil, err
		}
		node = &Node{
			Type:  operatorType(logicOper),
			Left:  node,
			Right: right,
		}
	}
}

// parseTerm 解析单个项
func (p *Parser) parseTerm() (*Node, error) {
	token, ok := p.consume()
	if !ok {
		return nil, fmt.Errorf("parse term but encounter end")
	}
	if token == "(" {
		expression, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		return expression, nil
	}

	left := token
	operator, ok := p.consume()
	if !ok {
		return nil, fmt.Errorf("expected operator but encounter end")
	}
	right, ok := p.consume()
	if !ok {
		return nil, fmt.Errorf("expected right value but encounter end")
	}
	return &Node{
		Type: NodeTypeComparison,
		Comparison: ComparisonNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		},
	}, nil
}

// consume 消耗符号
func (p *Parser) consume() (string, bool) {
	if p.current >= len(p.tokens) {
		return "", false
	}
	token := p.tokens[p.current]
	p.current++
	return token, true
}
