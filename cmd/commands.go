package cmd

import "github.com/efy/gorecall/subcmd"

var Commands = []subcmd.Command{
	serve,
	createuser,
	migrate,
	index,
	search,
}
