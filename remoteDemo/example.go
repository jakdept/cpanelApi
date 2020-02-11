package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/alecthomas/kingpin"
	"golang.org/x/crypto/ssh"
)

func SSHKeyfileInsecureRemote(username, keyFile string) (ssh.ClientConfig, error) {
	// read the keyfile
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	return ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint
	}, nil
}

// whmapi1 fetch_service_ssl_components parse out certificate
// https://documentation.cpanel.net/display/DD/WHM+API+1+Functions+-+fetch_service_ssl_components

func Connect(proto, host string, port int, creds ssh.ClientConfig) (string, error) {
	conn, err := ssh.Dial(proto, fmt.Sprintf("%s:%d", host, port), &creds)
	if err != nil {
		return nil, err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return nil, err
	}
	output, err := session.Output("whmapi1 create_user_session --output=json user=root service=whostmgrd")
	if err != nil {
		log.Fatalln(err)
	}

	token, err := parseUserSessionOutput(output)
	if err != nil {
		log.Fatalln(err)
	}
	return token, nil
}

var keyfile *string = kingpin.Flag("keyfile", "location to ssh key").Default("/root/.ssh/id_rsa").String()
var username *string = kingpin.Flag("username", "remote ssh user").Default("root").String()
var proto *string = kingpin.Flag("tcp", "ssh network protocol").Default("tcp").String()
var host *string = kingpin.Flag("host", "remote ssh host").Default("localhost").String()
var port *int = kingpin.Flag("port", "remote ssh port").Default("22").Int()

func parseUserSessionOutput(output []byte) (string, error) {
	unmarshalObject := struct {
		Data struct {
			Token string `json:"data"`
		}
	}{}

	err := json.Unmarshal(output, &unmarshalObject)
	if err != nil {
		return "", err
	}
	return unmarshalObject.Data.Token, nil
}

func main() {
	_ = kingpin.Parse()

	creds, err := SSHKeyfileInsecureRemote(*username, *keyfile)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := Connect(*proto, *host, *port, creds)
	if err != nil {
		log.Fatalln(err)
	}
	output, err := conn.Output("whmapi1 create_user_session --output=json user=root service=whostmgrd")
	if err != nil {
		log.Fatalln(err)
	}

	token, err := parseUserSessionOutput(output)
	if err != nil {
		log.Fatalln(err)
	}

	_ = token
}
