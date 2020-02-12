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

	api, err := cpanel.NewInsecureRemoteSSHWhmAPI(*username, *keyfile, *host, *port)
	if err != nil {
		log.Fatalln(err)
	}

	accounts, err := api.ListAccounts()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("All accounts on server:")
	for _, each := range accounts {
		fmt.Println(each)
	}
}
