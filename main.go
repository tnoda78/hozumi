package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"strings"
	"time"
)

var versionName = "1.2.5"

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
  -g. --graphical                       He dances. (This option is not supporting Windows OS)
`

var opts struct {
	Help bool `short:"h" long:"help" description:"He displays help message."`
	Version bool `short:"v" long:"version" description:"He displays his version."`
	Speed string `short:"s" long:"speed" description:"He displays by specified speed." default:"middle"`
	Cool bool `short:"c" long:"cool" description:"He sometimes shouts, \"Cool\"."`
	Graphical bool `short:"g" long:"graphical" description:"He dances. (This option is not supporting Windows OS)"`
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("\033[0m", sig)
			os.Exit(0)
		}
	}()
  hozumiWriter := setup()
  hozumiWriter.write()
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
	args, err := flags.Parse(&opts)
	if err != nil {
		displayHelpMessage()
	}
	if opts.Help {
		displayHelpMessage()
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
	}
	writer.intervalDisplayRow = intervalDisplayRow
	writer.intervalDisplayOneLetter = intervalDisplayOneLetter
	writer.cool = opts.Cool
	writer.intervalDisplayCool = 10 * time.Millisecond
	if opts.Graphical {
		writer.displayGraphicalLoop()
	}
	return writer
}

func displayHelpMessage() {
	fmt.Printf(message, versionName)
	os.Exit(0)
}

func (writer *HozumiWriter) write() {
	for {
		writer.displayContents()
		if writer.cool {
			writer.displayCool()
		}
	}
}

func (writer *HozumiWriter) displayContents() {
	for _, content := range writer.contents {
		writer.display(content)
		writer.displayAll(content)
	}
}

func (writer *HozumiWriter) display(content string) {
	writer.dashboard = append(writer.dashboard, "")
	row := len(writer.dashboard) - 1
	for _, letter := range strings.Split(content, "") {
		writer.dashboard[row] = letter
		writer.updateDashboard(writer.intervalDisplayRow)
	}
}

func (writer *HozumiWriter) displayAll(content string) {
	row := len(writer.dashboard) - 1
	str := ""
	for _, letter := range strings.Split(content, "") {
		str = str + letter
		writer.dashboard[row] = str
		writer.updateDashboard(writer.intervalDisplayOneLetter)
	}
	writer.updateDashboard(writer.intervalDisplayRow)
}

func (writer *HozumiWriter) displayCool() {
	content := "Cooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooool!"
	writer.dashboard = append(writer.dashboard, "")
	row := len(writer.dashboard) - 1
	for i := 0; i < len(content); i++ {
		writer.dashboard[row] = content[0:i + 1]
		writer.updateDashboard(writer.intervalDisplayCool)
	}
}

func (writer *HozumiWriter) displayGraphicalLoop() {
	str1 := []string{
	"\033[33m-  -  \033[34m----  \033[33m----  \033[34m-  -  \033[33m-   -  \033[34m--- ",
	"\033[33m-  -  \033[34m-  -  \033[33m  -   \033[34m-  -  \033[33m-- --  \033[34m -  ",
	"\033[33m----  \033[34m-  -  \033[33m -    \033[34m-  -  \033[33m- - -  \033[34m -  ",
	"\033[33m-  -  \033[34m-  -  \033[33m-     \033[34m-  -  \033[33m-   -  \033[34m -  ",
	"\033[33m-  -  \033[34m----  \033[33m----  \033[34m----  \033[33m-   -  \033[34m--- ",
	}

	for {
		for i := 0; i <= 40; i++ {
			writer.displayGraphical(str1, i)
		}
		for i := 40; i >= 0; i-- {
			writer.displayGraphical(str1, i)
		}
	}
}

func (writer *HozumiWriter) displayGraphical(content []string, space int) {
	var output = make([]string, len(content))
	for i := 0; i < len(content); i++ {
		output[i] = strings.Repeat(" ", space) + content[i]
	}
	writer.dashboard = output
	writer.updateDashboard(writer.intervalDisplayOneLetter)
}

func (writer *HozumiWriter) updateDashboard(interval time.Duration) {
	str := strings.Join(writer.dashboard, "\n") + "\n"
	os.Stdout.Write([]byte(str))
	time.Sleep(interval)
	for i := 0; i < len(writer.dashboard); i++ {
		fmt.Printf("\033[A\033[2K")
	}
}

