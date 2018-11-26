package document

import (
	"bytes"
	"fmt"
)

type Document struct {
	lines []Line
	site  int
}

//@TODO add vector clock support

//Returns position in the lines slice, if position doesn't exist then it will return false and the index where it WOULD go
func (this *Document) getLineIndex(pos []Position) (int, bool) {
	//var midpoint int = 0
	var toMidpoint int = 0
	linesToLook := this.lines

	for {
		if len(linesToLook) == 0 {
			return toMidpoint, false
		}

		midpoint := len(linesToLook) / 2
		posCompare := Compare(pos, linesToLook[midpoint].pos)
		if posCompare > 0 {
			linesToLook = linesToLook[0:midpoint]
		} else if posCompare < 0 {
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
	fmt.Printf("Input to Insert() got %s and %s\n", ToString(pos), content)
	fmt.Printf("%d\n", len(this.lines))
	if lineIdx, exists := this.getLineIndex(pos); !exists {
		this.lines = append(this.lines, Line{})
		copy(this.lines[lineIdx+1:], this.lines[lineIdx:])
		this.lines[lineIdx] = Line{pos, content}
		return
	}

	err = fmt.Errorf("Input to Insert() doesnt exist. Got %s", ToString(pos))

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
	fmt.Printf("Input to Move() got %s and %d\n", ToString(pos), moveAmount)
	fmt.Printf("%d\n", len(this.lines))
	if lineIdx, exists := this.getLineIndex(pos); exists {
		if lineIdx+moveAmount < len(this.lines) {
			newPos = this.lines[lineIdx+moveAmount].pos
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

func (this *Document) GeneratePositionBetween(l []Position, r []Position, site int) (pos []Position, err error) {
	return GeneratePositionBetween(l, r, site)
}
