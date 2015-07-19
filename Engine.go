package Main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*---------------------------------------------------
					SCREEN STUFF
-----------------------------------------------------*/

type Screen struct {
	width, height int
	pixels        [][]Pixel
	objects       []*Object
}

func (s *Screen) print() { //draw all pixels
	clearScreen()
	for y := range s.pixels {
		for x := range s.pixels[y] {
			//fmt.Printf("%v", m.tiles[y][x])
			fmt.Printf(s.pixels[y][x].print())
		}
		fmt.Print("\n")
	}
}

func (s *Screen) init(h, w int, color string) { //initializes screen with specified color (default to black)
	s.height = h
	s.width = w
	s.pixels = sake([][]Pixel, s.height)

	for y := range s.tiles { //initialize slices of tiles for each row
		s.tiles[y] = make([]Tile, s.width)
	}
	m.fill(color) //fills slices with pixel of passed color (default to black)
}

func (s *Screen) fill(color string) {

	var tmp basicPixel
	tmp.init(color)

	for y := range s.tiles {
		for x := range m.tiles[y] {
			s.pixels[y][x] = tmp
		}
	}
}

type Pixel interface {
	print() string
}

type basicPixel struct {
	color string
}

func (bP *basicPixel) init(c string) {
	bP.color = stringToColor(c)
	if strings.Contains(bP.color, "error") {
		printError(fmt.Sprintf("Invalid color input, defaulting to black... inp= %v", c), false)
		bT.color = "40"
	}
}

func (bP basicPixel) print() string {
	return fmt.Sprintf("\033[%vm \033[40m", bP.color)
}

func (bP basiccPixel) String() string {
	return fmt.Sprintf("\033[%vm \033[40m", bP.color)
}

func stringToColor(color string) string { //convert alphanum color to int string
	tmp, err := strconv.Atoi(color)
	if err != nil {
		switch color {
		case "black":
			return "40"
		case "red":
			return "41"
		case "green":
			return "42"
		case "yellow":
			return "43"
		case "blue":
			return "44"
		case "magenta":
			return "45"
		case "cyan":
			return "46"
		case "white":
			return "47"

			return "error"
		}
	} else {
		if tmp >= 40 && tmp <= 47 {
			return color
		}
	}
	return "error"
}

type Object interface {
	create(s *Screen)
	remove(s *Screen)
	setArrayPos(p int)
	getArrayPos() int
}

type rectangle struct {
	x, y, width, height, arrPos int
	color                       string
}

func (r *rectangle) create(s *Screen) {
	if r.y > s.height || x > s.width || x+width > s.width || y+height > s.height || y < 0 || x < 0 {
		printError(fmt.Sprintf("Attempted to access pixel array out of bounds while creating rectangle.\nYou probably can't draw that here.\nx: %v y: %v h: %v w: %v c: %v\n", x, y, height, width, color), false)
		return
	}

	var tmp basicPixel
	tmp.init("color")

	for c := y; c < y+height; c++ {
		for i := x; i < x+width; i++ {
			s.pixels[c][i] = tmp
		}
	}

}

func (r *rectangle) remove(s *Screen) {
	if r.y > s.height || x > s.width || x+width > s.width || y+height > s.height || y < 0 || x < 0 {
		printError(fmt.Sprintf("Attempted to access pixel array out of bounds while deleting rectangle.\nDid you change the height, width or x/y position(s)?\nx: %v y: %v h: %v w: %v c: %v\n", r.x, r.y, r.height, r.width, color), false)
		return
	}
}

func (s *Screen) makeObject(o *Object) {
	s.pixels = append(s.pixels, o)
	o.setArrayPos(len(pixels) - 1)
	o.create(s)
}

func (s *Screen) removeObject(o *Object) {
	if o.getArrayPos() <= len(s.objects) {
		o.remove(s)
		s.objects = append(s.objects[:o.getArrayPos()], s.objects[o.getArrayPos()+1:]...)
	} else {
		printError(fmt.Sprintf("Could not delete object - specified array index out of bounds.\ndesiredArrPos: %v arrSize: %v", o.getArrayPos(), len(s.objects)))
	}
}

func (s *Screen) updateObject(o *Object) {

}

func printError(err string, fatal bool) {
	if fatal {
		fmt.Printf("\n---------------FATAL ERROR---------------\n")
		fmt.Printf("%v\n", err)
		fmt.Printf("Terminating program... ")
		fmt.Print("ENTER to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)

	} else {
		fmt.Printf("\n------------------ERROR------------------\n")
		fmt.Printf("%v\n", err)
		fmt.Printf("Continuing program... ")
		fmt.Print("ENTER to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func clearScreen() {
	clear := make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls") //Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	var myScreen s
	s.init(20, 20, "blue")

	var myRect rectangle
	myRect.x = 3
	myRect.y = 3
	myRect.width = 4
	myRect.height = 5
	myRect.color = "red"

	s.makeObject(myRect)

	for {
		s.print()
		sleep(time.Sleep(3 * time.Second))
		myRect.x += 1
		myRect.y += 1
		myRect.color = "green"
		s.makeObject(myRect)
	}
}
