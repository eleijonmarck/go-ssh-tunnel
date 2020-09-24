package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/elliotchance/sshtunnel"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	client, err := startTunnel()
	if err != nil {
		log.Printf("errrrrorrr %s", err)
	}
	sftpclient, err := sftp.NewClient(client)
	if err != nil {
		log.Printf("unable to create sftp client with error: %s", err)
	}
	defer sftpclient.Close()

	// check it's there
	fi, err := sftpclient.Lstat("toSEB")
	if err != nil {
		log.Printf("seb upload test failed to upload a file w. error: %s", err)
	}
	log.Printf("fi %v", fi)
}

func startTunnel() (*ssh.Client, error) {
	portForwarded := "2000"
	sftpHostUser := "root"
	// Connection settings
	sftpHostServer := "localhost:2222"
	sftpRemoteServer := "localhost:3333"
	tunnel := sshtunnel.NewSSHTunnel(
		// User and host of tunnel server, it will default to port 22
		// if not specified.
		sftpHostUser+"@"+sftpHostServer,
		// authentication
		// auth, // 1. private key
		ssh.Password("root"), // 1. private key
		// The destination host and port of the actual server.
		sftpRemoteServer,
		// The local port you want to bind the remote port to.
		// Specifying "0" will lead to a random port.
		portForwarded,
	)

	// You can provide a logger for debugging, or remove this line to
	// make it silent.

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	// Start the server in the background. You will need to wait a
	// small amount of time for it to bind to the localhost port
	// before you can start sending connections.
	go tunnel.Start()

	// known io.Copy error:
	// even if we wait 5 seconds we still get logs of io.Copy error use of closed connection
	defer func() {
		time.Sleep(5 * time.Second)
		tunnel.Close()
	}()

	time.Sleep(100 * time.Millisecond)
	log.Printf("started tunnel")

	log.Printf("tunnel %+v", tunnel)

	sftpHostUser = "docker"
	keyPath := "./ssh_host_payout_staging_rsa_key"
	buff, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Printf("unable to read the privatye key %s, w. error: %s", keyPath, err)
	}
	signer, err := ssh.ParsePrivateKey(buff)
	if err != nil {
		log.Printf("unable to parse key")
	}
	auth := ssh.PublicKeys(signer)

	config := &ssh.ClientConfig{
		User: sftpHostUser,
		Auth: []ssh.AuthMethod{
			auth,
		},
		Timeout:         3 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	hostNetwork := "localhost"
	log.Printf("dialing %s:%s", hostNetwork, portForwarded)
	return ssh.Dial("tcp", hostNetwork+":"+portForwarded, config)
}
