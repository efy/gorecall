package cmd

import (
	"fmt"
	"os"

	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/subcmd"
)

var migrate = subcmd.Command{
	UsageLine: "migrate",
	Short:     "run any pending database migrations",
	Run: func(cmd *subcmd.Command, args []string) {
		dbdriver := cmd.Flag.String("dbdriver", "sqlite3", "driver of the database you intend to use (sqlite3, postgres)")
		dbdsn := cmd.Flag.String("dsn", "gorecall.db", "data source name")
		cmd.ParseFlags(args)

		db, err := database.Connect(database.Options{
			Driver: *dbdriver,
			DSN:    *dbdsn,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = database.Setup(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("success")
	},
}
