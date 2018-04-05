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
		dbdsn := cmd.Flag.String("dsn", "postgres://recall:recall@localhost/recall?sslmode=disable", "data source name")
		cmd.ParseFlags(args)

		dbopts := database.Options{
			Driver: "postgres",
			DSN:    *dbdsn,
		}

		db, err := database.Connect(dbopts)
		if err != nil {
			fmt.Println("error connecting to database:", err)
			os.Exit(1)
		}

		err = database.Setup(dbopts, db)
		if err != nil {
			fmt.Println("error migrating databaseL", err)
			os.Exit(0)
		}

		fmt.Println("migrate success")
	},
}
