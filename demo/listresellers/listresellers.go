package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kingpin"
	cpanel "github.com/jakdept/cpanelApi"
)

var keyfile *string = kingpin.Flag("keyfile", "location to ssh key").Default("/root/.ssh/id_rsa").String()
var username *string = kingpin.Flag("username", "remote ssh user").Default("root").String()
var host *string = kingpin.Flag("host", "remote ssh host").Default("localhost").String()
var port *int = kingpin.Flag("port", "remote ssh port").Default("22").Int()

func main() {
	_ = kingpin.Parse()

	api, err := cpanel.NewWhmApi(*host)
	if err != nil {
		log.Fatalln(err)
	}

	sshConfig, err := cpanel.InsecureSSHKeyfileConfig(*username, *keyfile)
	if err != nil {
		log.Fatalln(err)
	}

	err = api.SSHSessionAuthenticate(*host, *port, sshConfig)
	if err != nil {
		log.Fatalln(err)
	}

	resellers, err := api.ListAllResellerNames()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All resellers on server:")
	for _, each := range resellers {
		fmt.Println(each)
	}
}
