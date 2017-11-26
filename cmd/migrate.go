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
		dbname := cmd.Flag.String("dbname", "gorecall.db", "path to database file")
		cmd.ParseFlags(args)

		db, err := database.InitDatabase(*dbname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		database.MigrateDatabase(db)
	},
}
