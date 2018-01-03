package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/efy/gorecall/subcmd"
	"github.com/gorilla/securecookie"
)

var secret = subcmd.Command{
	UsageLine: "secret",
	Short:     "generate a base64 encoded secret key",
	Run: func(cmd *subcmd.Command, args []string) {
		length := cmd.Flag.Int("len", 32, "key length recommended 32 or 64")

		key := securecookie.GenerateRandomKey(*length)
		enc := base64.StdEncoding.EncodeToString(key)
		fmt.Println(enc)
		os.Exit(0)
	},
}
