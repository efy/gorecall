package subcmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Run       func(cmd *Command, args []string)
	UsageLine string
	Short     string
	Flag      flag.FlagSet

	// Store flags that have been explicitly set
	SetFlags map[string]bool
}

// ParseFlags provides a wrapper around flag.Parse
// in order to build the SetFlags map
func (c *Command) ParseFlags(args []string) {
	c.SetFlags = make(map[string]bool)
	c.Flag.Parse(args)
	c.Flag.Visit(func(f *flag.Flag) {
		c.SetFlags[f.Name] = true
	})
}

func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i > 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %\n\n", c.UsageLine)
	os.Exit(1)
}

func (c *Command) Runnable() bool {
	if c.Run != nil {
		return true
	}
	return false
}

func (c *Command) FlagIsSet(name string) bool {
	v, ok := c.SetFlags[name]
	if v && ok {
		return true
	}
	return false
}
