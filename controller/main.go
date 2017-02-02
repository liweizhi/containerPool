package main

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/liweizhi/containerPool/controller/cmd"
)
func main() {
	app := cli.NewApp()
	app.Name = "ContainerPool"
	app.Version = "1.0"
	app.Authors = []cli.Author{
		{
			Name: "liweizhi",
			Email: "li41898@163.com",
		},
	}
	app.Usage = "UI fo docker management"
	app.Before = cli.BeforeFunc(func(c *cli.Context) error{
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	})

	app.Commands = []cli.Command{
		{
			Name: "server",
			Usage: "server for ContainerPool",
			Action: cmd.Server,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "listen address",
					Value: ":8080",
				},
				cli.StringFlag{
					Name:  "rethinkdb-addr",
					Usage: "RethinkDB address",
					Value: "rethinkdb:28015",
				},
				cli.StringFlag{
					Name:  "rethinkdb-auth-key",
					Usage: "RethinkDB auth key",
					Value: "",
				},
				cli.StringFlag{
					Name:  "rethinkdb-database",
					Usage: "RethinkDB database name",
					Value: "dockerpool",
				},

				cli.StringFlag{
					Name:   "docker, d",
					Value:  "tcp://127.0.0.1:2375",
					Usage:  "docker swarm addr",
					EnvVar: "DOCKER_HOST",
				},

				cli.StringSliceFlag{
					Name:  "auth-whitelist-cidr",
					Usage: "whitelist CIDR to bypass auth",
					Value: &cli.StringSlice{},
				},
			},

		},
	}

	app.Flags = []cli.Flag{

		cli.BoolFlag{
			Name: "debug, D",
			Usage: "enabel debug mode",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}


}
