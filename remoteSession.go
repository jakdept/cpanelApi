package cpanel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func NewRemoteWhmAPI() WhmAPI {
	return WhmAPI{}
}

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

func Connect(proto, host string, port int, creds ssh.ClientConfig) (string, error) {
	conn, err := ssh.Dial(proto, fmt.Sprintf("%s:%d", host, port), &creds)
	if err != nil {
		return "", err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return "", err
	}
	output, err := session.Output("whmapi1 create_user_session --output=json user=root service=whostmgrd")
	if err != nil {
		return "", err
	}

	token, err := parseUserSessionOutput(output)
	if err != nil {
		return "", err
	}
	return token, nil
}

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
