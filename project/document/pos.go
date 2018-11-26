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

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
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

	diffenceLen := len(r) - len(l)
	smallestLen := min(len(r), len(l))
	addFinalDigit := false
	var lastIdx int

	for idx := 0; idx < smallestLen; idx++ {
		lastIdx = idx

		lPos := l[idx]
		rPos := r[idx]

		var difference = rPos.posId - lPos.posId

		if difference == 0 {
			//Add the digit, same so it doesnt matter
			pos = append(pos, Position{rPos.posId, rPos.site})
			addFinalDigit = true
		} else if difference == 1 {
			//pos = append(pos, Position{lPos.posId, lPos.site})
			if idx < len(l)-1 {
				pos = append(pos, Position{lPos.posId, lPos.site})
				addFinalDigit = true
				break
			} else {
				pos = append(pos, Position{lPos.posId, lPos.site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), lPos.site})
				addFinalDigit = false
			}
		} else if difference > 1 {
			//pos = append(pos, Position{lPos.posId, lPos.site})
			if idx < len(l)-1 {
				pos = append(pos, Position{lPos.posId, lPos.site})
				addFinalDigit = true
				break
			} else {
				pos = append(pos, Position{uint8(randBetween(int(lPos.posId), int(rPos.posId))), lPos.site})
				addFinalDigit = false
			}
		} else {
			panic("Difference GeneratePositionBetween() is less than 0! in This should never happen")
		}
	}

	if addFinalDigit {
		if diffenceLen > 0 {
			nextRightPos := r[lastIdx+1].posId
			for nextRightPos == 0 {
				pos = append(pos, Position{0, site})
				lastIdx++
				nextRightPos = r[lastIdx].posId
			}

			if nextRightPos == 1 {
				pos = append(pos, Position{0, site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
			} else {
				pos = append(pos, Position{uint8(randBetween(0, int(nextRightPos))), site})
			}
		} else if diffenceLen < 0 {
			nextLeftPos := l[lastIdx+1].posId
			for nextLeftPos == math.MaxUint8 {
				pos = append(pos, Position{math.MaxUint8, site})
				lastIdx++
				nextLeftPos = l[lastIdx].posId
			}

			//fmt.Printf("nextLeftPos: %d\n", nextLeftPos)
			if nextLeftPos == math.MaxUint8-1 {
				pos = append(pos, Position{math.MaxUint8, site})
				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
			} else {
				pos = append(pos, Position{uint8(randBetween(int(nextLeftPos), math.MaxUint8)), site})
			}
		} else {
			//No need to do anything when they are the same length
		}
	}

	// lIdx := 0
	// rIdx := 0
	// for {
	// 	lPos := l[lIdx]
	// 	rPos := r[rIdx]

	// 	if rPos.posId == lPos.posId && rPos.site == lPos.site {
	// 		pos = append(pos, Position{rPos.posId, rPos.site})
	// 	} else {
	// 		var difference = rPos.posId - lPos.posId
	// 		if difference > 1 {
	// 			pos = append(pos, Position{uint8(randBetween(int(lPos.posId), int(rPos.posId))), site})
	// 			break
	// 		} else if difference == 0 {
	// 			panic("Hello")
	// 			if rPos.site > site && lPos.site < site {
	// 				pos = append(pos, Position{lPos.posId, lPos.site})
	// 			} else {
	// 				pos = append(pos, Position{lPos.posId, lPos.site})
	// 				pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
	// 			}
	// 		} else if difference == 1 {
	// 			//When we have to split an identical position id then our actions depend on the site
	// 			//I believe this is because we just want some consistent behavior that we can count on, the paper didnt really explain...
	// 			if site > lPos.site {
	// 				pos = append(pos, Position{lPos.posId, site})
	// 			} else if site < rPos.site {
	// 				pos = append(pos, Position{rPos.posId, site})
	// 			} else {
	// 				pos = append(pos, Position{lPos.posId, lPos.site}, Position{uint8(randBetween(0, math.MaxUint8)), site})
	// 			}
	// 			break
	// 		} else {
	// 			panic("Difference GeneratePositionBetween() is less than 0! in This should never happen")
	// 		}
	// 	}

	// 	atLEnd := lIdx == len(l)-1
	// 	atREnd := rIdx == len(r)-1

	// 	if !atLEnd {
	// 		lIdx++
	// 	}

	// 	if !atREnd {
	// 		rIdx++
	// 	}

	// 	if atLEnd && atREnd {
	// 		break
	// 	}

	// }

	///ADSASD

	// for idx := 0; idx < len(l); idx++ {
	// 	if idx > len(r)-1 {
	// 		break
	// 	}

	// 	lPos := l[idx]
	// 	rPos := r[idx] //@TODO Could this index out of bounds ??

	// 	//When they are the same we just go down a digit
	// 	if rPos.posId == lPos.posId && rPos.site == lPos.site {
	// 		pos = append(pos, Position{rPos.posId, rPos.site})
	// 		continue
	// 	}
	// 	var difference = rPos.posId - lPos.posId
	// 	if difference > 1 {
	// 		pos = append(pos, Position{uint8(randBetween(int(lPos.posId), int(rPos.posId))), site})
	// 	} else if difference < 1 {
	// 		if rPos.site > site && lPos.site < site {
	// 			pos = append(pos, Position{lPos.posId, lPos.site})
	// 		} else {
	// 			pos = append(pos, Position{lPos.posId, lPos.site})
	// 			pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
	// 		}
	// 	} else {
	// 		//When we have to split an identical position id then our actions depend on the site
	// 		//I believe this is because we just want some consistent behavior that we can count on, the paper didnt really explain...
	// 		if site > lPos.site {
	// 			pos = append(pos, Position{lPos.posId, site})
	// 		} else if site < rPos.site {
	// 			pos = append(pos, Position{rPos.posId, site})
	// 		} else {
	// 			pos = append(pos, Position{lPos.posId, lPos.site})
	// 			pos = append(pos, Position{uint8(randBetween(0, math.MaxUint8)), site})
	// 		}
	// 	}
	// }

	// if len(r) < len(l) {
	// 	var leftNextId uint8
	// 	for idx := 0; idx < len(l)-len(r); idx++ {
	// 		leftNextId = l[len(r)+idx].posId
	// 		//fmt.Printf("%d\n", rightNextId)
	// 		if leftNextId == math.MaxUint8 {
	// 			pos = append(pos, Position{math.MaxUint8, site})
	// 			//fmt.Printf("appended %s\n", ToString(pos))
	// 			leftNextId = 0
	// 		} else {
	// 			break
	// 		}
	// 	}

	// 	pos = append(pos, Position{uint8(randBetween(int(leftNextId), math.MaxUint8)), site})
	// }

	// if len(l) < len(r) {
	// 	var rightNextId uint8
	// 	for idx := 0; idx < len(r)-len(l); idx++ {
	// 		rightNextId = r[len(l)+idx].posId
	// 		//fmt.Printf("%d\n", rightNextId)
	// 		if rightNextId == 1 || rightNextId == 0 {
	// 			pos = append(pos, Position{0, site})
	// 			//fmt.Printf("appended %s\n", ToString(pos))
	// 			rightNextId = math.MaxUint8
	// 		} else {
	// 			break
	// 		}
	// 	}

	// 	pos = append(pos, Position{uint8(randBetween(0, int(rightNextId))), site})
	// }

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
