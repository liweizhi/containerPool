package cmd

import (

	"github.com/urfave/cli"
	"github.com/samalba/dockerclient"
	"fmt"
	"os"
	"github.com/liweizhi/containerPool/controller/manager"
	"github.com/liweizhi/containerPool/auth"
	log"github.com/Sirupsen/logrus"
)
func Server(c *cli.Context) {
	fmt.Println("hello world")

	listenAddr := c.String("listen")
	rethinkdbAddr := c.String("rethinkdb-addr")
	rethinkdbAuthKey := c.String("rethinkdb-auth-key")
	rethinkdbName := c.String("rethinkdb-database")
	dockerUrl := c.String("docker")



	client, err := dockerclient.NewDockerClient(dockerUrl, nil)

	if err != nil{
		log.Fatalln(err)
	}

	authenticator := auth.NewAuthenticator("builtin")

	controllerManager, err := manager.NewManager(rethinkdbAddr, rethinkdbName, rethinkdbAuthKey, client, authenticator)
	if err != nil{
		log.Fatalln(err)
	}













}
