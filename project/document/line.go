package document

import (
	"bytes"
	"fmt"
)

//Not neccesarily an actual line with an EOL, more like a set of characters. The paper uses this term so I will too
type Line struct {
	pos     []Position
	content string
}

func (this *Line) ToString() string {
	var posBytes bytes.Buffer
	for posIdx := range this.pos {
		if posIdx != 0 {
			posBytes.WriteString(".")
		}
		posBytes.WriteString(this.pos[posIdx].ToString())
	}

	return fmt.Sprintf("<%s,%s>", posBytes.String(), this.content)
}
