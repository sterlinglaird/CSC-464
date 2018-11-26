package document

import (
	"bytes"
	"fmt"
)

type Document struct {
	lines []Line
	site  int
}

//@TODO add vector clock support for unique sites

//Returns position in the lines slice, if position doesn't exist then it will return false and the index where it WOULD go
func (this *Document) getLineIndex(pos []Position) (int, bool) {
	//var midpoint int = 0
	var toMidpoint int = 0
	linesToLook := this.lines

	for {
		if len(linesToLook) == 0 {
			fmt.Println("ended")
			return toMidpoint, false
		}

		midpoint := len(linesToLook) / 2

		fmt.Printf("mid:%d toMid:%d\n", midpoint, toMidpoint)
		posCompare := Compare(pos, linesToLook[midpoint].pos)
		if posCompare < 0 {
			linesToLook = linesToLook[0:midpoint]
		} else if posCompare > 0 {
			toMidpoint += midpoint + 1
			linesToLook = linesToLook[midpoint+1:]
		} else {
			return midpoint + toMidpoint, true
		}
	}
}

func NewDocument(site int) Document {
	return Document{[]Line{Line{StartPos, ""}, Line{EndPos, ""}}, site}
}

func (this *Document) Insert(pos []Position, content string) (err error) {
	fmt.Printf("%s\n", this.toString())
	fmt.Printf("Input to Insert() got %s and '%s'\n", ToString(pos), content)
	if lineIdx, exists := this.getLineIndex(pos); !exists {
		fmt.Printf("%d\n", lineIdx)
		this.lines = append(this.lines, Line{})
		copy(this.lines[lineIdx+1:], this.lines[lineIdx:])
		this.lines[lineIdx] = Line{pos, content}
		return
	}

	err = fmt.Errorf("Input to Insert() already exists. Got %s", ToString(pos))

	return
}

func (this *Document) InsertLeft(pos []Position, content string) (err error) {
	return
}

func (this *Document) InsertRight(pos []Position, content string) (err error) {
	return
}

func (this *Document) Delete(pos []Position) (err error) {

	return
}

//Returns a moved position based on numMove. Positive is right neg is left. Based on input position
func (this *Document) Move(pos []Position, moveAmount int) (newPos []Position, err error) {
	fmt.Printf("%s\n", this.toString())
	fmt.Printf("Input to Move() got %s and %d\n", ToString(pos), moveAmount)
	if lineIdx, exists := this.getLineIndex(pos); exists {
		fmt.Printf("%d\n", lineIdx)
		if lineIdx+moveAmount < len(this.lines) {
			if moveAmount > 0 {
				newPos, err = GeneratePositionBetween(pos, this.lines[lineIdx+moveAmount].pos, this.site)
			} else {
				newPos, err = GeneratePositionBetween(this.lines[lineIdx+moveAmount].pos, pos, this.site)
			}
		} else {
			err = fmt.Errorf("Input to Move() (moveAmount) Too extreme. Got %d", moveAmount)
		}
		return
	}

	err = fmt.Errorf("Input to Move() doesnt exist. Got %s", ToString(pos))

	return
}

func (this *Document) MoveRight(pos []Position) (newPos []Position, err error) {

	return
}

func (this *Document) GetContent() (content string, err error) {
	var contentBytes bytes.Buffer
	for lineIdx := range this.lines {
		contentBytes.WriteString(this.lines[lineIdx].content)
	}

	content = contentBytes.String()

	return
}

func (this *Document) GetContentAt(pos []Position) (content string, err error) {

	return
}

//@TODO use the site value from document instead of passing it in
func (this *Document) GeneratePositionBetween(l []Position, r []Position, site int) (pos []Position, err error) {
	return GeneratePositionBetween(l, r, site)
}

func (this *Document) toString() string {
	var docBytes bytes.Buffer
	for lineIdx := range this.lines {
		docBytes.WriteString("\t")
		docBytes.WriteString(this.lines[lineIdx].ToString())
		docBytes.WriteString(".\n")
	}

	return fmt.Sprintf("<%s>", docBytes.String())
}
