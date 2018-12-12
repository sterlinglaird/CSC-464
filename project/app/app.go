package main

import (
	"fmt"
	"log"
	"os"

	"../document"
	"github.com/jroimartin/gocui"
)

var doc = document.NewDocument(1)
var currentPos = document.StartPos

type Editor struct{}

func (e *Editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	var err error
	log.Printf("Current pos %s\n", document.ToString(currentPos))
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		currentPos, err = doc.InsertRight(currentPos, string(ch))
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		currentPos, err = doc.Move(currentPos, 1)
		err = doc.DeleteLeft(currentPos)
	// case key == gocui.KeyDelete:
	// 	v.EditDelete(false)
	// 	err = doc.DeleteRight(currentPos)
	// case key == gocui.KeyInsert:
	// 	v.Overwrite = !v.Overwrite
	// case key == gocui.KeyEnter:
	// 	v.EditNewLine()
	// case key == gocui.KeyArrowDown:
	// 	v.MoveCursor(0, 1, false)
	// case key == gocui.KeyArrowUp:
	// 	v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
		currentPos, err = doc.Move(currentPos, -1)
	case key == gocui.KeyArrowRight:
		currentPos, err = doc.Move(currentPos, 1)
		v.MoveCursor(1, 0, false)
	}

	log.Printf("Now %s\n", doc.ToString())

	if err != nil {
		log.Printf("Error in edit(): %s\n", err.Error())
		log.Printf("%s\n", doc.ToString())
		panic(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		v.Editor = &Editor{}
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var f, _ = os.OpenFile("testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Printf("Error in main loop: %s\n", err.Error())
		fmt.Printf("%s\n", doc.ToString())
		return
	}

}
