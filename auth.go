package cpanel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

var CreateUserSessionCmd = []string{
	"/usr/bin/env",
	"whmapi1",
	"create_user_session",
	"--output=json",
	"user=root",
	"service=whostmgrd",
	"preferred_domain=",
}

func (a *WhmAPI) ActivateTokenUrl(url, token, hostname string) error {
	// note this method is documented at:
	// https://documentation.cpanel.net/display/DD/Guide+to+API+Authentication+-+Single+Sign+On
	a.token = &token
	a.hostname = &hostname

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	// TODO: need more error checking?
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func parseUserSessionOutput(output []byte) (string, string, error) {
	var unmarshalObject struct {
		Data struct {
			Activate string `json:"url"`
			Token    string `json:"cp_security_token"`
		}
	}

	err := json.Unmarshal(output, &unmarshalObject)
	if err != nil {
		return "", "", err
	}
	return unmarshalObject.Data.Activate, unmarshalObject.Data.Token, nil
}

func (a *WhmAPI) SSHSessionAuthenticate(
	hostname string,
	port int,
	config ssh.ClientConfig,
) error {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), &config)
	if err != nil {
		return err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return err
	}

	output, err := session.Output(strings.Join(CreateUserSessionCmd, " ") + hostname)
	if err != nil {
		return err
	}

	activateUrl, token, err := parseUserSessionOutput(output)

	err = a.ActivateTokenUrl(activateUrl, token, hostname)
	if err != nil {
		return err
	}

	return nil
}

func (a *WhmAPI) LocalSessionAuthenticate() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	cmd := exec.Command(CreateUserSessionCmd[0],
		append(CreateUserSessionCmd[1:], hostname)...)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	activateUrl, token, err := parseUserSessionOutput(output)
	if err != nil {
		return err
	}

	err = a.ActivateTokenUrl(activateUrl, token, hostname)
	if err != nil {
		return err
	}

	return nil
}

func InsecureSSHKeyfileConfig(username, keyFile string) (ssh.ClientConfig, error) {
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

// https://documentation.cpanel.net/display/DD/Guide+to+API+Authentication+-+API+Tokens+in+WHM
// func (a *WhmAPI) APIKeyAuthenticate(key, domain string) {
// 	// set up the required object for apikey based authentication
// 	a.auth = map[string]string{
// 		"apikey": strings.TrimSpace(key),
// 	}
// 	if domain != "" {
// 		a.auth["domain"] = domain
// 	}
// }

// https://documentation.cpanel.net/display/DD/Guide+to+API+Authentication+-+Username+and+Password+Authentication
// func (a *WhmAPI) UserAuthenticate(username, password, domain string) {
// 	// set up the required object for user based authentication
// 	a.auth = map[string]string{
// 		"email":    username,
// 		"password": password,
// 	}
// 	if domain != "" {
// 		a.auth["domain"] = domain
// 	}
// }
