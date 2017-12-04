package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/subcmd"
)

var index = subcmd.Command{
	UsageLine: "index",
	Short:     "update the search index",
	Run: func(cmd *subcmd.Command, args []string) {
		dbname := cmd.Flag.String("dbname", "gorecall.db", "path to database file")
		indexname := cmd.Flag.String("indexname", "gorecall.idx", "path to search index")
		cmd.ParseFlags(args)

		db, err := database.Init(*dbname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		mapping := bleve.NewIndexMapping()
		index, err := bleve.New(*indexname, mapping)
		if err != nil {
			fmt.Println(err)
			index, err = bleve.Open(*indexname)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
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

		idxErrs := make([]error, 0)
		idxCount := 0

		for _, bm := range bookmarks {
			id := strconv.FormatInt(bm.ID, 10)
			err := index.Index(id, bm)
			if err != nil {
				idxErrs = append(idxErrs, err)
				continue
			}
			idxCount++
		}

		fmt.Printf("Index complete. %d indexed, %d errors\n", idxCount, len(idxErrs))
	},
}
