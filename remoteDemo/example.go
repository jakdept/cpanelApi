package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/alecthomas/kingpin"
	"golang.org/x/crypto/ssh"
)

func SshKeyfileInsecureRemote(username, keyFile string) (ssh.ClientConfig, error) {
	// read the keyfile
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)

	return ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

func Connect(proto, host, port string, creds ssh.ClientConfig) (*ssh.Session, error) {
	conn, err := ssh.Dial(proto, host+":"+port, &creds)
	if err != nil {
		return nil, err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return session, nil
}

var keyfile *string = kingpin.Flag("keyfile", "location to ssh key").Default("/root/.ssh/id_rsa").String()
var username *string = kingpin.Flag("username", "remote ssh user").Default("root").String()
var proto *string = kingpin.Flag("tcp", "ssh network protocol").Default("tcp").String()
var host *string = kingpin.Flag("host", "remote ssh host").Default("localhost").String()

func main() {
	_ = kingpin.Parse()

	creds, err := SshKeyfileInsecureRemote(*username, *keyFile)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := Connect(proto, host, port, creds)
	if err != nil {
		log.Fatalln(err)
	}

	output, err := conn.Output("whmapi1 create_user_session output=json user=root service=whostmgrd locale=en")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(output)

}
