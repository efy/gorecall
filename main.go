package main

import (
	"flag"
	"fmt"
)

var (
	fcgi = flag.Bool("fcgi", false, "Use Fast CGI")
	addr = flag.String("addr", ":8080", "Bind address")
)

func main() {
	flag.Parse()

	if *fcgi {
		fmt.Println("Run FastCGI server")
		return
	}

	fmt.Println("Run http server", *addr)
}
