package cpanel

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/publicsuffix"
)

func NewInsecureRemoteSSHWhmAPI(username, keyFile, hostname string, port int) (WhmAPI, error) {
	creds, err := SSHKeyfileInsecureRemote(username, keyFile)
	if err != nil {
		return WhmAPI{}, err
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), &creds)
	if err != nil {
		return WhmAPI{}, err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return WhmAPI{}, err
	}

	cmd := "whmapi1"
	cmd += " create_user_session"
	cmd += " --output=json"
	cmd += " user=root"
	cmd += " service=whostmgrd"
	cmd += " preferred_domain=" + hostname

	output, err := session.Output(cmd)
	if err != nil {
		return WhmAPI{}, err
	}

	token, err := parseUserSessionOutput(output)
	if err != nil {
		return WhmAPI{}, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return WhmAPI{}, err
	}

	api := WhmAPI{
		hostname: hostname,
		token:    token,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Jar: jar,
		},
	}
	return api, nil
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
