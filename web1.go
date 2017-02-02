package main

import (
	"io"
	"net/http"
	"os"
	"github.com/urfave/cli"
	"fmt"

)

type MyHandle struct{}

func main() {
	var language string

	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: &language,
		},
	}

	fmt.Println(language)
	app.Before = func(c *cli.Context) error {
		fmt.Println("before")
		if c.GlobalBool("lang") {
			fmt.Println("lang")
		}
		return nil
	}
	app.Action = func(c *cli.Context) error {
		name := "someone"
		if c.NArg() > 0 {
			name = c.Args()[0]
		}
		fmt.Println(language)
		if language == "spanish" {
			fmt.Println("Hola", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil
	}

	app.Run(os.Args)
}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())
}