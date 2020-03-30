package main

import (
	"fmt"

	"github.com/altid/libs/fs"
)

func errorWrite(c *fs.Control, err error) {
	ew, _ := c.ErrorWriter()
	defer ew.Close()

	fmt.Fprintf(ew, "discordfs: %v\n", err)
}
