package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/jcrummy/sbconfig/sb"
)

func main() {
	selectInterface()
	discover(selectedInterface())

	t := prompt.New(executor, completer,
		prompt.OptionTitle("sbconfig: Squeeze Box Configurator"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionShowCompletionAtStart(),
	)
	t.Run()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "configure", Description: "Configure selected device (ex: 'configure 0')"},
		{Text: "discover", Description: "Search network for devices"},
		{Text: "exit", Description: "Exit program"},
		{Text: "interface", Description: "Select a new interface to use"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(s string) {
	s = strings.TrimSpace(s)

	cmd := strings.Split(s, " ")[0]

	switch cmd {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)

	case "discover":
		discover(selectedInterface())

	case "interface":
		selectInterface()

	case "configure":
		fields := strings.Split(s, " ")
		if len(fields) < 2 {
			fmt.Println("No device specified. Use 'configure 0' to configure the first device.")
			return
		}
		deviceIndex, err := strconv.Atoi(fields[1])
		if err != nil {
			fmt.Println("Not a number. Use 'configure 0' to configure the first device.")
			return
		}
		if (deviceIndex > len(sbs)-1) || (deviceIndex < 0) {
			fmt.Println("No such device. Use 'configure 0' to configure the first device.")
			return
		}
		fmt.Printf("Configuring device #%d (%+v).\n", deviceIndex, sbs[deviceIndex].MacAddr)
		sb.Configure(&sbs[deviceIndex], selectedInterface())
	}
	return
}
