package main

import (
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// NodeType 表示节点类型
type NodeType int

const (
	NodeTypeAnd NodeType = iota
	NodeTypeOr
	NodeTypeComparison
)

// Node 表示语法树节点
type Node struct {
	Type       NodeType
	Left       *Node
	Right      *Node
	Comparison ComparisonNode
}

type NodeEx struct {
	Type       NodeType
	Childs     []*NodeEx
	Comparison ComparisonNode
}

// ComparisonNode 表示比较操作符节点
type ComparisonNode struct {
	Left     string
	Operator string
	Right    interface{}
}

func toNodeEx(pnode *Node) *NodeEx {
	ptmp := new(NodeEx)
	ptmp.Type = pnode.Type
	if ptmp.Type == NodeTypeComparison {
		ptmp.Comparison = pnode.Comparison
		return ptmp
	}
	left := toNodeEx(pnode.Left)
	right := toNodeEx(pnode.Right)
	if left.Type == ptmp.Type {
		ptmp.Childs = append(ptmp.Childs, left.Childs...)
	} else {
		ptmp.Childs = append(ptmp.Childs, left)
	}
	if right.Type == ptmp.Type {
		ptmp.Childs = append(ptmp.Childs, right.Childs...)
	} else {
		ptmp.Childs = append(ptmp.Childs, right)
	}
	return ptmp
}

func convert(pnode *NodeEx) elastic.Query {
	if pnode.Type == NodeTypeComparison {
		return elastic.NewMatchQuery(pnode.Comparison.Left, pnode.Comparison.Right)
	}
	boolQuery := elastic.NewBoolQuery()
	if pnode.Type == NodeTypeOr {
		boolQuery.MinimumNumberShouldMatch(1)
	}
	for _, v := range pnode.Childs {
		tmpquery := convert(v)
		if pnode.Type == NodeTypeAnd {
			boolQuery.Filter(tmpquery)
		} else {
			boolQuery.Should(tmpquery)
		}
	}
	return boolQuery
}

// operatorType 返回操作符类型
func operatorType(op string) NodeType {
	switch op {
	case "and":
		return NodeTypeAnd
	case "or":
		return NodeTypeOr
	default:
		panic("Unknown operator")
	}
}

func printQuery(q elastic.Query) {
	mapsource, _ := q.Source()
	jsonBytes, _ := json.MarshalIndent(mapsource, "", "\t")
	fmt.Println(string(jsonBytes))
}
