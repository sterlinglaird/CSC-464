package document

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
)

type Position struct {
	posId uint8
	site  int
}

//These are the pre-existing lines at the start and end of the document. All following inserts will be somewhere between these
var (
	StartPos = []Position{{0, 0}}
	EndPos   = []Position{{math.MaxUint8, 0}}
)

//Random number between x and y non inclusive on both ends (x,y)
func randBetween(x int, y int) int {
	//fmt.Printf("%d %d\n", x, y)
	return int(rand.Intn(int(y-x-1)) + int(x) + 1)
}

//Returns 1 if x > y  -1 if x < y, and 0 if x = y
func Compare(x []Position, y []Position) int {
	for idx := range x {
		if len(y) == idx {
			return 1
		}

		xPos := x[idx].posId
		yPos := y[idx].posId

		xSite := x[idx].site
		ySite := y[idx].site

		if xPos > yPos {
			return 1
		} else if xPos < yPos {
			return -1
		}

		if xSite > ySite {
			return 1
		} else if xSite < ySite {
			return -1
		}
	}

	if len(x) < len(y) {
		return -1
	}

	return 0
}

//According the definitions given the paper
//@TODO make this function nicer, I had to make a couple hotixes so it has a number of edge cases
//@TODO also need to handle other edge cases of over/underflow of digits
//@TODO make sure I am assigning the sites correctly
func GeneratePositionBetween(l []Position, r []Position, site int) (pos []Position, err error) {
	//@TODO verify l < r

	fmt.Printf("GeneratePositionBetween %s %s\n", ToString(l), ToString(r))
	//Loop over left side since we want to make a position bigger than this
	for idx := 0; idx < len(l); idx++ {
		lPos := l[idx]
		rPos := r[idx] //@TODO Could this indec out of bounds ??

		//When they are the same we just go down a digit
		if rPos.posId == lPos.posId && rPos.site == lPos.site {
			pos = append(pos, Position{rPos.posId, rPos.site})
			continue
		}
		var difference = rPos.posId - lPos.posId
		if difference > 1 {
			pos = append(pos, Position{uint8(randBetween(int(lPos.posId), int(rPos.posId))), site})
		} else if difference < 1 {
			if rPos.site > site && lPos.site < site {
				pos = append(pos, Position{lPos.posId, lPos.site})
			} else {
				pos = append(pos, Position{lPos.posId, lPos.site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
			}
		} else {
			//When we have to split an identical position id then our actions depend on the site
			//I believe this is because we just want some consistent behavior that we can count on, the paper didnt really explain...
			if site > lPos.site {
				pos = append(pos, Position{lPos.posId, site})
			} else if site < rPos.site {
				pos = append(pos, Position{rPos.posId, site})
			} else {
				pos = append(pos, Position{lPos.posId, lPos.site}, Position{uint8(randBetween(0, math.MaxUint8)), site})
			}
		}
	}

	if len(l) < len(r) {
		var rightNextId uint8
		for idx := 0; idx < len(r)-len(l); idx++ {
			rightNextId = r[len(l)+idx].posId
			fmt.Printf("%d\n", rightNextId)
			if rightNextId == 1 || rightNextId == 0 {
				pos = append(pos, Position{0, site})
				fmt.Printf("appended %s\n", ToString(pos))
				rightNextId = math.MaxUint8
			} else {
				break
			}
		}

		pos = append(pos, Position{uint8(randBetween(0, int(rightNextId))), site})

		// //Edge case, we need to add a new digit
		// rightNextId := r[len(l)].posId
		// if rightNextId == 1 || rightNextId == 0 {
		// 	pos = append(pos, Position{0, site})
		// 	pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
		// } else {
		// 	pos = append(pos, Position{uint8(randBetween(0, int(rightNextId))), site})
		// }
	}

	fmt.Printf("Created %s\n", ToString(pos))
	return
}

func (this *Position) ToString() string {
	return fmt.Sprintf("<%d,%d>", this.posId, this.site)
}

func ToString(pos []Position) string {
	var posBytes bytes.Buffer
	for posIdx := range pos {
		if posIdx != 0 {
			posBytes.WriteString(".")
		}
		posBytes.WriteString(pos[posIdx].ToString())
	}

	return posBytes.String()
}
