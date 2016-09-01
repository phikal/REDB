// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
)

var words []string

func main() {
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read stdin: %s\n", err)
		os.Exit(1)
	}

	words = strings.Split(string(content), "\n")
	words = words[0:len(words)-1]

	fmt.Printf("Read %d words. Starting server...\n", len(words))

	// start ReGeX game server on main thread
	gameServer()
}
