package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/efy/gorecall/auth"
	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/subcmd"
)

var createuser = subcmd.Command{
	UsageLine: "createuser",
	Short:     "create a new user in the database",
	Run: func(cmd *subcmd.Command, args []string) {
		dbdsn := cmd.Flag.String("dsn", "postgres://recall:recall@localhost/recall?sslmode=disable", "data source name")
		cmd.ParseFlags(args)

		db, err := database.Connect(database.Options{
			Driver: "postgres",
			DSN:    *dbdsn,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		uRepo, err := datastore.NewUserRepo(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Creating user...")

		fmt.Println("Enter username:")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		user, err := uRepo.GetByUsername(username)
		if err != nil && user != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if user != nil {
			fmt.Println("User already exists")
			os.Exit(1)
		}

		fmt.Println("Enter password:")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if len(password) < 8 {
			fmt.Println("Password should be at least 8 characters long")
			os.Exit(1)
		}

		hash, err := auth.HashPassword(password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		u := &datastore.User{
			Username: username,
			Password: hash,
		}

		u, err = uRepo.Create(u)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("created user:", u.Username)
		os.Exit(0)
	},
}
