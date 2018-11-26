package document

import (
	"fmt"
	"math"
	"math/rand"
)

type Position struct {
	posId uint
	site  int
}

//These are the pre-existing lines at the start and end of the document. All following inserts will be somewhere between these
var (
	StartPos = []Position{{0, 0}}
	EndPos   = []Position{{math.MaxUint32, 0}}
)

//Random number between x and y non inclusive on both ends (x,y)
func randBetween(x uint, y uint) uint {
	return uint(rand.Intn(int(y-x-1)) + 1 + int(x))
}

//Returns 1 if x > y  -1 if x < y, and 0 if x = y
func Compare(x []Position, y []Position) int {
	for idx := range x {
		xPos := x[idx].posId
		yPos := y[idx].posId

		xSite := x[idx].site
		ySite := y[idx].site

		if len(y) == idx {
			return 1
		}

		if xPos > yPos || xSite > ySite {
			return 1
		}

		if xPos < yPos || xSite < ySite {
			return -1
		}
	}

	if len(x) < len(y) {
		return -1
	}

	return 0
}

//According the definitions given the paper
func GeneratePositionBetween(l []Position, r []Position, site int) (pos []Position, err error) {
	//@TODO verify l < r

	//Loop over left side since we want to make a position bigger than this
	for idx := 0; idx < len(l); idx++ {
		lPos := l[idx]
		rPos := r[idx]

		if rPos.posId == lPos.posId && rPos.site == lPos.site {
			pos = append(pos, Position{rPos.posId, rPos.site})
			continue
		}

		var posId uint
		//var lowerBound uint

		var difference = rPos.posId - lPos.posId
		if difference > 1 {
			posId = rPos.posId
		} else if difference < 1 {
			panic(fmt.Sprintf("right smaller than left! This should never happen"))
		} else {
			//When we have to split an identical position id then our actions depend on the site
			//I believe this is because we just want some consistent behavior that we can count on, the paper didnt really explain...
			if site > lPos.site {
				pos = append(pos, Position{lPos.posId, site})
			} else if site < rPos.site {
				pos = append(pos, Position{rPos.posId, site})
			} else {
				pos = append(pos, Position{lPos.posId, lPos.site}, Position{randBetween(0, math.MaxUint32), site})
			}
		}

		pos = append(pos, Position{posId, site})
	}
	return
}
