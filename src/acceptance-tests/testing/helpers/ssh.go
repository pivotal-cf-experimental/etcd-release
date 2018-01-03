package helpers

import (
	"bytes"
	"fmt"
	"log"

	"github.com/cloudfoundry/bosh-cli/director"
	directorCli "github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"golang.org/x/crypto/ssh"
)

func RunSSHCommand(server string, port int, username string, privateKey string, command string) (string, error) {
	parsedPrivateKey, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		log.Println(err)
		return "", err
	}

	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(parsedPrivateKey),
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, port), config)
	if err != nil {
		return "", errors.WrapError(err, "Cannot dial")
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", errors.WrapError(err, "Cannot create session")
	}
	defer session.Close()

	var output bytes.Buffer

	session.Stdout = &output

	session.Run(command)

	return output.String(), nil
}

func GetSSHCreds(deploymentName string, boshCliClient director.Director) (string, string, string, error) {
	deployment, err := boshCliClient.FindDeployment(deploymentName)
	if err != nil {
		return "", "", "", err
	}

	sshOpts, privateKey, err := directorCli.NewSSHOpts(uuid.NewGenerator())
	if err != nil {
		return "", "", "", err
	}

	slug, err := directorCli.NewAllOrInstanceGroupOrInstanceSlugFromString("etcd")
	if err != nil {
		return "", "", "", err
	}

	sshResult, err := deployment.SetUpSSH(slug, sshOpts)
	if err != nil {
		return "", "", "", err
	}

	return sshResult.Hosts[0].Host, sshOpts.Username, privateKey, nil
}
