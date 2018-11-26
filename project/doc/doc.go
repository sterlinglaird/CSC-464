package doc

import (
	"math"
)

type Doc struct {
	lines []Line
	site  int
}

type Line struct {
	pos     []Pos
	content string
}

type Pos struct {
	posId int
	site  int
}

var (
	StartLine = []Pos{{0, 0}}
	EndLine   = []Pos{{math.MaxUint32, 0}}
)
