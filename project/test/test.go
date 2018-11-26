package main

import (
	"fmt"

	d "../document"
)

func main() {
	site := 1
	doc := d.NewDocument(site)

	pos, err := doc.GeneratePositionBetween(d.StartPos, d.EndPos, site)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	doc.Insert(pos, "h")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	pos, err = doc.Move(pos, 1)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	err = doc.Insert(pos, "2")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	pos, err = doc.Move(pos, 1)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	err = doc.Insert(pos, "w")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	content, err := doc.GetContent()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	fmt.Printf("%s\n", content)
}
