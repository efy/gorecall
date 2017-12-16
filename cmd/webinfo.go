package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/importer"
	"github.com/efy/gorecall/subcmd"
)

var webinfo = subcmd.Command{
	UsageLine: "webinfo",
	Short:     "process actual web requests for bookmark links and populate metadata fields",
	Run: func(cmd *subcmd.Command, args []string) {
		dbdriver := cmd.Flag.String("dbdriver", "sqlite3", "driver of the database you intend to use (sqlite3, postgres)")
		dbdsn := cmd.Flag.String("dsn", "gorecall.db", "data source name")
		cmd.ParseFlags(args)

		dbopts := database.Options{
			Driver: *dbdriver,
			DSN:    *dbdsn,
		}

		db, err := database.Connect(dbopts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bookmarkRepo, err := datastore.NewBookmarkRepo(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bookmarks, err := bookmarkRepo.GetAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Starting webinfo...")
		start := time.Now()
		bookmarks = importer.BatchWebinfoSerial(bookmarks)
		dur := time.Now().Sub(start)

		fmt.Println("Saving bookmarks...")
		errcount := 0
		for _, b := range bookmarks {
			_, err := bookmarkRepo.Update(&b)
			if err != nil {
				errcount++
				fmt.Printf("failed to save: %s", b.Title)
			}
		}

		fmt.Printf("Retrieved info for %d links with %d errors. Took %v\n", len(bookmarks), errcount, dur)
	},
}
