package sb

import (
	"fmt"
	"strings"
)

func (c *configurator) executor(s string) {
	s = strings.TrimSpace(s)

	cmd := strings.Split(s, " ")[0]

	switch cmd {

	case "show":
		c.showValues()

	case "set":
		c.setValue(s)

	case "save":
		c.saveValues()

	case "exit":
		fmt.Println("Press Ctrl-D to exit.")
	}
	return
}
