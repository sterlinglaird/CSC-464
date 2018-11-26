package document

//Not neccesarily an actual line with an EOL, more like a set of characters. The paper uses this term so I will too
type Line struct {
	pos     []Position
	content string
}
