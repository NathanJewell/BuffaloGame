/*
The MIT License (MIT)

Copyright (c) 2015 Nathan Jewell

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

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

const ()

/*---------------------------------------
		ITEM CLASSES
-----------------------------------------*/
type Item interface {
	use(b *Buffalo) bool
}

type statPot struct {
	name, stat string
	strength   int
	useable    bool
}

func (sP statPot) use(b *Buffalo) bool {
	fmt.Printf("Using %v", sP)
	switch sP.stat {
	case "fullness":
		b.feed(sP.strength)
	case "happiness":
		b.entertain(sP.strength)
	case "energy":
		b.energize(sP.strength)
	}
	return sP.useable
}

func (sP statPot) String() string {
	return fmt.Sprintf("statpot of type %v and strength %v.\n", sP.stat, sP.strength)
}

type emptyItem struct {
	useable bool
}

func (ei emptyItem) use(b *Buffalo) bool {
	fmt.Println("This is not a useable item.\n")
	return false
}

func (ei emptyItem) String() string {
	return "...\n"
}

type weaponItem struct {
	useable bool
}

func (wI weaponItem) use(b *Buffalo) {

}

/*---------------------------------------
		INVENTORY CLASS
-----------------------------------------*/
type Inventory struct {
	items []Item
}

func (i *Inventory) init() {
	i.items = make([]Item, 1)
}

func (i Inventory) String() string {
	var toreturn string = ""
	toreturn += fmt.Sprintf("------------INVENTORY------------\n\n")
	for c := 0; c < len(i.items); c++ {
		toreturn += fmt.Sprintf("%v: %v", c+1, i.items[c])
	}
	toreturn += fmt.Sprintf("\n---------------------------------\n")
	return toreturn
}

func (i *Inventory) add(it Item) {
	i.items = append(i.items, it)
}

func (i *Inventory) use(which int, b *Buffalo) {
	if which <= len(i.items) {
		fmt.Printf("%v", i.items[which-1])
		if i.items[which-1].use(b) {
			i.items = append(i.items[:which-1], i.items[which:]...)

		}
	} else {
		fmt.Println("\nThat is not a valid item.")
	}

}

/*---------------------------------------
		BUFFALO CLASS
-----------------------------------------*/

type Buffalo struct {
	name                               string
	fullness, happiness, energy, combo int
	fBar, hBar, eBar                   string
	inv                                Inventory
}

func (b Buffalo) init() {
	b.inv.init()
}

//feed, entertain and energize simply add/subract a certain stat from the buffalo
//they also give a small print statement to notify the user
func (b *Buffalo) feed(a int) {
	b.fullness += a
	if a < 0 {
		fmt.Printf("%v FULLNESS\n", a)
	} else {
		fmt.Printf("+%v FULLNESS\n", a)
	}

}

func (b *Buffalo) entertain(a int) {
	b.happiness += a
	if a < 0 {
		fmt.Printf("%v HAPPINESS\n", a)
	} else {
		fmt.Printf("+%v HAPPINESS\n", a)
	}

}

func (b *Buffalo) energize(a int) {
	b.energy += a
	if a < 0 {
		fmt.Printf("%v ENERGY\n", a)
	} else {
		fmt.Printf("+%v ENERGY\n", a)
	}
}

func (b Buffalo) String() string {
	var s string = ""
	s += fmt.Sprintf("   FULLNESS: %v [%v] \n", b.fullness, b.fBar)
	s += fmt.Sprintf("  HAPPINESS: %v [%v] \n", b.happiness, b.hBar)
	s += fmt.Sprintf("     ENERGY: %v [%v] \n", b.energy, b.eBar)
	return s
}

func (b *Buffalo) calcBars() {
	//iterates through stats and creates strings for each one that reflect amount of each one
	b.fBar, b.hBar, b.eBar = "", "", ""
	for i := 0; i < 10; i++ {
		if i < b.fullness {
			b.fBar += "|"
		} else {
			b.fBar += " "
		}
		if i < b.happiness {
			b.hBar += "|"
		} else {
			b.hBar += " "
		}
		if i < b.energy {
			b.eBar += "|"
		} else {
			b.eBar += " "
		}
	}
	b.combo = b.energy + b.fullness + b.happiness
}

func (b Buffalo) printInv() {
	fmt.Printf("%v", b.inv)
}

/*---------------------------------------
		MAP FUNCTIONS
-----------------------------------------*/
type Map struct {
	width, height int
	tiles         [][]Tile
}

func (m *Map) print() {
	clearScreen()
	for y := range m.tiles {
		for x := range m.tiles[y] {
			//fmt.Printf("%v", m.tiles[y][x])
			fmt.Printf(m.tiles[y][x].print())
		}
		fmt.Print("\n")
	}
}

func (m *Map) init(h, w int) {
	m.height = h
	m.width = w
	m.tiles = make([][]Tile, m.height)
	fmt.Println("HERE")
	for y := range m.tiles {
		m.tiles[y] = make([]Tile, m.width)
	}
	fmt.Println("HERE")
	m.create()
	fmt.Println("HERE")
}

func (m *Map) create() {
	for y := range m.tiles {
		for x := range m.tiles[y] {
			var tmp basicTile
			tmp.init("green")
			m.tiles[y][x] = tmp
		}
	}
}

type Tile interface {
	print() string
}

type basicTile struct {
	color string
}

func (bT *basicTile) init(c string) {
	bT.color = stringToColor(c)
	if strings.Contains(bT.color, "error") {
		printError(fmt.Sprintf("Invalid color input, defaulting to black... inp= %v", c), false)
		bT.color = "40"
	}
}

func (bT basicTile) print() string {
	return fmt.Sprintf("\033[%vm \033[40m", bT.color)
}

func (bT basicTile) String() string {
	return fmt.Sprintf("\033[%vm \033[40m", bT.color)
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

func (m *Map) createSquare(x, y, width, height int, color string) {
	if y > m.height || x > m.width || x+width > m.width || y+height > m.height || y < 0 || x < 0 {
		printError(fmt.Sprintf("Attempted out of bound to 'pixel' array.\nYou probably can't draw that here.\nx: %v y: %v h: %v w: %v c: %v\n", x, y, height, width, color), false)
		return
	}

	var tmp basicTile
	tmp.init("color")

	for c := y; c < y+height; c++ {
		for i := x; i < x+width; i++ {
			m.tiles[c][i] = tmp
		}
	}

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

/*
40	Black
41	Red
42	Green
43	Yellow
44	Blue
45	Magenta
46	Cyan
47	White
*/

/*---------------------------------------
		MAIN FUNCTIONS
-----------------------------------------*/
func main() {
	fmt.Printf("\03316m")
	var myB Buffalo
	myB.init()
	myB.feed(9)
	myB.entertain(9)
	myB.energize(9)
	myB.calcBars()

	var hapPot statPot
	hapPot.stat = "happiness"
	hapPot.strength = 3
	hapPot.useable = true

	var eI emptyItem
	eI.useable = false

	myB.inv.add(hapPot)
	myB.inv.add(hapPot)
	myB.inv.add(eI)

	alive := true

	var gameMap Map
	gameMap.init(10, 10)
	gameMap.createSquare(2, 2, 3, 3, "blue")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("      WELCOME TO BUFFALO SIMULATOR 2K15 - The wildest one in the west")
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("\nInstructions: Don't kill your furry buffalo friend too much.")
	fmt.Println("                    Use 'help' to view other commands")
	fmt.Println("                             HAVE FUN!!!\n\n ")
	for alive {
		if !checkHealth(&myB) {
			break
		}
		myB.energy -= 1
		myB.happiness -= 1
		myB.fullness -= 1

		for {
			gameMap.print()
			fmt.Println("Enter a command: ")
			cmdInput, _ := reader.ReadString('\n')
			if !doCMDIn(cmdInput, &myB) { //parse and execute command

			} else {
				break
			}

		}

		alive = checkHealth(&myB) //confirm life
	}

	fmt.Println("\n\n\n--------------------------------------------------------------------------")
	fmt.Println("                                GAME OVER")
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("          \nInstruction: Go cry in a corner, you killed your")
	fmt.Println("                         innocent buffalo friend. RIP")
	fmt.Println("\n To start again, start the program again because implementing a long for loop is hard...\nCome again soon.")

}

func printBuffalo(b Buffalo) {
	b.calcBars()
	fmt.Println("\nWelcome to buffalo stats v.01. Please enjoy! \n")
	if b.combo < 10 {
		fmt.Println("You buffalo could die soon! Oh no")
	} else if b.combo < 20 {
		fmt.Println("Your buffalo is doing ok.")
	} else if b.combo <= 30 {
		fmt.Println("You buffalo is doing jolly well tyvm.")
	}

	fmt.Printf("%v\n", b)

}

func checkHealth(b *Buffalo) bool {
	if b.energy-1 <= 0 || b.fullness-1 <= 0 || b.happiness-1 <= 0 {
		return false //buffalo has died
	} else {
		return true //buffalo is well and alive
	}
}

func doCMDIn(s string, b *Buffalo) bool {
	if strings.Contains(s, "help") {
		fmt.Println("\nType 'explore', 'play', 'rest' or 'eat' to do it.")
		fmt.Println("Type stat to see your current energy, happiness and fullness.")
		fmt.Println("View inventory with 'inv' and use item with 'use [num]'.\n")
		return false
	} else if strings.Contains(s, "explore") {
		/*
			fmt.Println("\nYou explored quite the good amount of stuff. Something bad could've happened though...")
			b.entertain(genRand(-3, 5))
			b.energize(genRand(-3, 5))
			b.feed(genRand(-3, 5))
			fmt.Println("")
		*/
	} else if strings.Contains(s, "play") {
		fmt.Println("\nYou frolicked in the flowers and sunshine.")
		b.entertain(2)
		fmt.Println("")
	} else if strings.Contains(s, "rest") {
		fmt.Println("\nNice rest its good cause now you're not so tired.")
		b.energize(2)
		fmt.Println("")
	} else if strings.Contains(s, "eat") {
		fmt.Println("\nYou ate a delicious tasty meal of tofu and brussel sprout salad!")
		b.feed(2)
		fmt.Println("")
	} else if strings.Contains(s, "stat") {
		printBuffalo(*b)
		return false
	} else if strings.Contains(s, "inv") {
		b.printInv()
		return false
	} else if strings.Contains(s, "use") {
		var char string
		var is bool

		for __, i := range s {
			var _ = __
			char = string(i)

			if char == " " {
				is = true
			} else if is {
				break
			}

		}
		if is {
			num, err := strconv.Atoi(char)
			if err != nil {
				fmt.Println("That item does not exist.\n")
			} else {
				b.inv.use(num, b)
			}

		}

		return false

	} else {
		fmt.Println("Invalid command - do it again.\n")
		return false
	}
	return true //command successfully executed
}

func genRand(l, t int) int {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(8)
	randomNum -= 3
	return randomNum
}
