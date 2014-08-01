package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var versionName = "2.0.0"

var message = `
Hozumi Command

USAGE:
  hozumi [option] param1 param2 param3...

VERSION:
  %s

OPTIONS:
  -h, --help                            He displays help message.
  -v, --version                         He displays his version.
  -s, --speed {low|middle|high}         He displays by specified speed.
  -c. --cool                            He sometimes shouts, "Cool".
  -p  --parallel                        He displays in parallel.
`

var opts struct {
	Help     bool   `short:"h" long:"help" description:"He displays help message."`
	Version  bool   `short:"v" long:"version" description:"He displays his version."`
	Speed    string `short:"s" long:"speed" description:"He displays by specified speed." default:"middle"`
	Cool     bool   `short:"c" long:"cool" description:"He sometimes shouts, \"Cool\"."`
	Parallel bool   `short:"p" long:"parallel" description:"He displeys in parallel."`
}

func main() {
	hozumiWriter := setup()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	go hozumiWriter.write()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break loop
		case termbox.EventResize:
			break loop
		}
	}
}

type HozumiWriter struct {
	contents                 []string
	intervalDisplayRow       time.Duration
	intervalDisplayOneLetter time.Duration
	intervalDisplayCool      time.Duration
	cool                     bool
	dashboard                []string
}

func setup() *HozumiWriter {
	p := flags.NewParser(&opts, flags.PrintErrors)
	args, err := p.Parse()

	if err != nil {
		displayHelpMessage()
		os.Exit(1)
	}

	if opts.Help {
		displayHelpMessage()
		os.Exit(0)
	}

	if opts.Version {
		fmt.Printf("Hozumi Command version (%s)\n", versionName)
		os.Exit(0)
	}

	writer := new(HozumiWriter)

	if len(args) > 0 {
		writer.contents = args
	} else {
		writer.contents = []string{"ほずみ"}
	}
	var intervalDisplayRow time.Duration
	var intervalDisplayOneLetter time.Duration

	if opts.Speed == "low" {
		intervalDisplayRow = 600 * time.Millisecond
		intervalDisplayOneLetter = 270 * time.Millisecond
	} else if opts.Speed == "middle" {
		intervalDisplayRow = 300 * time.Millisecond
		intervalDisplayOneLetter = 180 * time.Millisecond
	} else if opts.Speed == "high" {
		intervalDisplayRow = 150 * time.Millisecond
		intervalDisplayOneLetter = 90 * time.Millisecond
	} else if opts.Speed != "" {
		displayHelpMessage()
		os.Exit(1)
	}
	writer.intervalDisplayRow = intervalDisplayRow
	writer.intervalDisplayOneLetter = intervalDisplayOneLetter
	writer.cool = opts.Cool
	writer.intervalDisplayCool = 10 * time.Millisecond
	return writer
}

func displayHelpMessage() {
	fmt.Printf(message, versionName)
}

func (writer *HozumiWriter) write() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, ymax := termbox.Size()
	for {
		y := 0
		for y < ymax {
			if opts.Parallel {
				writer.writer_contents_pararrel(y)
			} else {
				writer.write_contents(y)
			}
			y = y + len(writer.contents)
			if writer.cool {
				y++
			}
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (writer *HozumiWriter) write_contents(row int) {
	for i, content := range writer.contents {
		writer.draw_row(content, row+i)
		time.Sleep(90 * time.Millisecond)
	}
	if writer.cool {
		writer.displayCool(row + len(writer.contents))
	}
}

func (writer *HozumiWriter) writer_contents_pararrel(row int) {
	var channum int
	if writer.cool {
		channum = len(writer.contents) + 1
	} else {
		channum = len(writer.contents)
	}

	c := make(chan int, channum)

	for i, content := range writer.contents {
		go func() {
			writer.draw_row(content, row+i)
			c <- 1
		}()
		time.Sleep(90 * time.Millisecond)
	}
	if writer.cool {
		go func() {
			writer.displayCool(row + len(writer.contents))
			c <- 1
		}()
	}
	for i := 0; i < channum; i++ {
		<-c
	}
}

func (writer *HozumiWriter) draw_row(message string, y int) {
	str := ""
	for _, letter := range strings.Split(message, "") {
		str = str + letter
		set_row(str, y, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(writer.intervalDisplayRow)
	}
	clear_row(y)
	time.Sleep(writer.intervalDisplayRow)
	str = ""
	for _, letter := range strings.Split(message, "") {
		str = str + letter
		set_row(str, y, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(writer.intervalDisplayOneLetter)
	}
	time.Sleep(writer.intervalDisplayRow)
}

func (writer *HozumiWriter) displayCool(y int) {
	content := "Cooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooool!"
	str := ""
	for _, letter := range strings.Split(content, "") {
		str = str + letter
		set_row(str, y, termbox.ColorYellow)
		termbox.Flush()
		time.Sleep(writer.intervalDisplayCool)
	}
}

func set_row(message string, y int, fg termbox.Attribute) {

	x := 0

	for len(message) > 0 {
		c, w := utf8.DecodeRuneInString(message)
		if c == utf8.RuneError {
			c = '?'
			w = 1
		}
		message = message[w:]
		termbox.SetCell(x, y, c, fg, termbox.ColorDefault)
		x += runewidth.RuneWidth(c)
	}
}

func clear_row(y int) {
	xmax, _ := termbox.Size()
	for x := 0; x < xmax; x++ {
		termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}
