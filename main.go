package main

import "fmt"

// main 函数示例用法
func main() {
	query := `(((ip==1.1.1.1))||(protocol=http))`
	parser := NewParser(query)
	pnode, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	printQuery(convert(toNodeEx(pnode)))

}
