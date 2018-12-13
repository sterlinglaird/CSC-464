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
//@TODO the start and end lines should be abstracted away more. User shouldnt be aware at all of them. Currently you can move to them just not delete, maybe start with an initial line?
//@TODO Maybe only use tombstone/site id for deleted positions instead of actually deleting them
//@TODO breaks when commas,periods,brackets etc appear in text

//Returns position in the lines slice, if position doesn't exist then it will return false and the index where it WOULD go
func (this *Document) GetLineIndex(pos []Position) (int, bool) {
	//var midpoint int = 0
	var toMidpoint int = 0
	linesToLook := this.lines

	for {
		if len(linesToLook) == 0 {
			//fmt.Println("ended")
			//fmt.Printf("Returned idx to: %s\n", ToString(this.lines[toMidpoint].Pos))
			return toMidpoint, false
		}

		midpoint := len(linesToLook) / 2

		//fmt.Printf("mid:%d toMid:%d\n", midpoint, toMidpoint)
		posCompare := Compare(pos, linesToLook[midpoint].Pos)
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
	doc := Document{[]Line{Line{StartPos, ""}, Line{EndPos, ""}}, site}
	return doc
}

//@TODO Should this be public?
func (this *Document) Insert(pos []Position, content string) (err error) {
	//fmt.Printf("Input to Insert() got %s and '%s'\n", ToString(pos), content)
	if lineIdx, exists := this.GetLineIndex(pos); !exists {
		this.lines = append(this.lines, Line{})
		copy(this.lines[lineIdx+1:], this.lines[lineIdx:])
		this.lines[lineIdx] = Line{pos, content}
		//fmt.Printf("Now %s\n", this.ToString())
		this.site++
		return
	}

	err = fmt.Errorf("Input to Insert() already exists. Got %s", ToString(pos))
	return
}

func (this *Document) InsertRight(pos []Position, content string) (newPos []Position, err error) {
	//fmt.Printf("Input to InsertRight() got %s and '%s'\n", ToString(pos), content)
	rightPos, err := this.Move(pos, 1)
	if err != nil {
		return
	}

	newPos, err = this.GeneratePositionBetween(pos, rightPos, this.site)
	if err != nil {
		return
	}

	err = this.Insert(newPos, content)

	return
}

func (this *Document) InsertLeft(pos []Position, content string) (newPos []Position, err error) {
	//fmt.Printf("Input to InsertLeft() got %s and '%s'\n", ToString(pos), content)
	leftPos, err := this.Move(pos, -1)
	if err != nil {
		return
	}

	newPos, err = this.GeneratePositionBetween(leftPos, pos, this.site)
	if err != nil {
		return
	}
	err = this.Insert(newPos, content)

	return
}

func (this *Document) Delete(pos []Position) (err error) {
	//fmt.Printf("Input to Delete() got %s\n", ToString(pos))
	if lineIdx, exists := this.GetLineIndex(pos); exists {
		if lineIdx == 0 || lineIdx == len(this.lines)-1 {
			err = fmt.Errorf("Input to Delete() cannot be deleted. Got %s", ToString(pos))
			return
		}

		copy(this.lines[lineIdx:], this.lines[lineIdx+1:])
		this.lines = this.lines[:len(this.lines)-1]
		//fmt.Printf("Now %s\n", this.ToString())
		this.site++
		return
	}

	err = fmt.Errorf("Input to Delete() does not exist. Got %s", ToString(pos))
	return
}

func (this *Document) DeleteRight(pos []Position) (err error) {
	toDeletePos, err := this.Move(pos, 1)
	if err != nil {
		return
	}

	return this.Delete(toDeletePos)
}

func (this *Document) DeleteLeft(pos []Position) (err error) {
	toDeletePos, err := this.Move(pos, -1)
	if err != nil {
		return
	}

	return this.Delete(toDeletePos)
}

//Returns a moved position based on numMove. Positive is right neg is left. Based on input position
func (this *Document) Move(pos []Position, moveAmount int) (newPos []Position, err error) {
	//log.Printf("Input to Move() got %s and %d\n", ToString(pos), moveAmount)
	if lineIdx, exists := this.GetLineIndex(pos); exists {
		if lineIdx+moveAmount < len(this.lines) && lineIdx+moveAmount >= 0 {
			newPos = this.lines[lineIdx+moveAmount].Pos
			return
		} else {
			err = fmt.Errorf("Input to Move() (moveAmount) Too extreme. Got %d, currernt position %s", moveAmount, ToString(pos))
		}
		return
	}

	err = fmt.Errorf("Input to Move() doesnt exist. Got %s", ToString(pos))

	return
}

func (this *Document) MoveRight(pos []Position) (newPos []Position, err error) {
	//fmt.Printf("Input to MoveRight() got %s\n", ToString(pos))
	if lineIdx, exists := this.GetLineIndex(pos); exists {
		if lineIdx+1 < len(this.lines) {
			newPos, err = GeneratePositionBetween(pos, this.lines[lineIdx+1].Pos, this.site)
			return
		} else {
			err = fmt.Errorf("Cannot MoveRight()")
		}
		return
	}

	err = fmt.Errorf("Input to MoveRight() doesnt exist. Got %s", ToString(pos))
	return
}

func (this *Document) MoveLeft(pos []Position) (newPos []Position, err error) {
	//fmt.Printf("Input to MoveLeft() got %s\n", ToString(pos))
	if lineIdx, exists := this.GetLineIndex(pos); exists {
		if lineIdx-1 < len(this.lines) {
			newPos, err = GeneratePositionBetween(pos, this.lines[lineIdx-1].Pos, this.site)
			return
		} else {
			err = fmt.Errorf("Cannot MoveLeft()")
		}
		return
	}

	err = fmt.Errorf("Input to MoveLeft() doesnt exist. Got %s", ToString(pos))
	return
}

func (this *Document) GetContent() (content string, err error) {
	var contentBytes bytes.Buffer
	for lineIdx := range this.lines {
		contentBytes.WriteString(this.lines[lineIdx].Content)
	}

	content = contentBytes.String()

	return
}

func (this *Document) GetContentAt(pos []Position) (content string, err error) {
	return
}

func (this *Document) GetLength() (length int, err error) {
	length = len(this.lines)
	return
}

func (this *Document) GetLine(idx int) (line Line) {
	line = this.lines[idx]
	return
}

//@TODO use the site value from document instead of passing it in
func (this *Document) GeneratePositionBetween(l []Position, r []Position, site int) (pos []Position, err error) {
	return GeneratePositionBetween(l, r, site)
}

func (this *Document) ToString() string {
	var docBytes bytes.Buffer
	for lineIdx := range this.lines {
		docBytes.WriteString("\t")
		docBytes.WriteString(this.lines[lineIdx].ToString())
		docBytes.WriteString(".\n")
	}

	return fmt.Sprintf("<%s>", docBytes.String())
}
