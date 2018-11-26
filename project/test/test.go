package main

import (
	"fmt"

	d "../document"
)

func main() {
	site := 1
	doc := d.NewDocument(site)

	pos, _ := doc.GeneratePositionBetween(d.StartPos, d.EndPos, site)
	doc.Insert(pos, "h")

	content, _ := doc.GetContent()

	fmt.Printf("%s\n", content)
}
