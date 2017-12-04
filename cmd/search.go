package cmd

import (
	"fmt"
	"os"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/simple"
	"github.com/efy/gorecall/subcmd"
)

var search = subcmd.Command{
	UsageLine: "search",
	Short:     "perform a search against the index",
	Run: func(cmd *subcmd.Command, args []string) {
		indexname := cmd.Flag.String("indexname", "gorecall.idx", "path to index directory")
		query := cmd.Flag.String("q", "", "search query")
		cmd.ParseFlags(args)

		index, err := bleve.Open(*indexname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *query == "" {
			fmt.Println("must provide query")
			os.Exit(1)
		}

		q := bleve.NewMatchQuery(*query)
		s := bleve.NewSearchRequest(q)
		s.Highlight = bleve.NewHighlight()
		s.Highlight.AddField("Title")
		result, err := index.Search(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(result)
	},
}
