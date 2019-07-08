package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sud - server upload/download
// Author: Matthew Arnold <matt@mattarnster.co.uk>

type config struct {
	Protocol             string   // The protocol we use - sftp for now.
	Type                 string   // Upload or download
	Hosts                []string // Array of hosts to use
	User                 string   // Username to use when connecting
	KeyPath              string   // Path to the user's private key
	AllFiles             bool     // Should we download all fo the files?
	SourceDirectory      string   // In this specific directory - ABSOLUTE PATH ONLY
	DestinationDirectory string   // Where we should put the downloaded files
	Files                []string // Or these specific ones?
	DeleteOnRetrieve     bool     // Should we delete the files from the server after we retreive them?
}

func main() {
	log.Println("sud - alpha")
	config := readConfig()

	if config.Protocol == "sftp" {
		connectSftp(config)
	}
}

func readConfig() config {
	file, err := os.Open("sud.json")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := config{}
	decodeErr := decoder.Decode(&configuration)

	if decodeErr != nil {
		log.Fatalln(decodeErr.Error())
	}

	return configuration
}

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func connectSftp(c config) {
	for _, host := range c.Hosts {

		log.Println("Connecting to: " + host)

		sshconfig := &ssh.ClientConfig{
			User: c.User,
			Auth: []ssh.AuthMethod{
				publicKey(c.KeyPath),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		conn, err := ssh.Dial("tcp", host, sshconfig)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		log.Println("SSH Connection established")

		// Create ourselves a new client
		client, sftpErr := sftp.NewClient(conn)
		defer client.Close()

		if sftpErr != nil {
			log.Fatalln(host)
		}

		log.Println("SFTP session created")

		if c.Type == "download" && c.AllFiles == true {
			log.Println("(remote) Directory listing: " + c.SourceDirectory)
			info, direrr := client.ReadDir(c.SourceDirectory)
			if direrr != nil {
				log.Println("Hello")
				log.Fatalln(direrr.Error())
			}
			log.Printf("(remote) Got info: %v \n", info)
			for _, file := range info {
				log.Printf("(remote) File info for: " + file.Name())
				if file.IsDir() {
					continue // We don't care about directories right now, just the files.
				}
				log.Println("(local) Does " + path.Join(c.DestinationDirectory, file.Name()) + " exist?")
				if _, err := os.Stat(path.Join(c.DestinationDirectory, file.Name())); os.IsNotExist(err) {
					// We don't have the file locally, create and then download it.
					log.Println("(local) Trying to create: " + path.Join(c.DestinationDirectory, file.Name()))
					f, fileerr := os.Create(path.Join(c.DestinationDirectory, file.Name()))
					if fileerr != nil {
						log.Fatalln(fileerr.Error())
					}

					log.Println("(remote) Trying to open: " + path.Join(c.SourceDirectory, file.Name()))
					sourceFile, err := client.OpenFile(path.Join(c.SourceDirectory, file.Name()), 0)
					if err != nil {
						log.Fatalln(err.Error())
					}

					log.Println("(local/remote) Starting copy operation")
					_, writeErr := io.Copy(f, sourceFile)
					if writeErr != nil {
						log.Println("(local/remote) Copy operation failed")
						log.Fatalln(writeErr.Error())
					} else {
						log.Println("(local/remote) File transfer successful: " + path.Join(c.DestinationDirectory, file.Name()))
						if c.DeleteOnRetrieve == true {
							log.Println("(info) Delete on retreive is set, deleting file...")
							delerr := client.Remove(path.Join(c.SourceDirectory, file.Name()))
							if delerr != nil {
								log.Printf("(remote) Failed to delete file: %s\n", path.Join(c.SourceDirectory, file.Name()))
							}
						}
					}

				} else if err == nil {
					// The file is already there
					log.Printf("(local) File already exists locally: %s\n", path.Join(c.DestinationDirectory, file.Name()))
					err := os.Remove(path.Join(c.DestinationDirectory, file.Name()))
					if err != nil {
						log.Fatalln(err.Error())
					}

					f, fileerr := os.Create(path.Join(c.DestinationDirectory, file.Name()))
					if fileerr != nil {
						log.Fatalln(fileerr.Error())
					}

					sourceFile, err := client.OpenFile(path.Join(c.DestinationDirectory, file.Name()), 0)
					if err != nil {
						log.Fatalln(err.Error())
					}

					_, writeErr := io.Copy(f, sourceFile)
					if writeErr != nil {
						log.Fatalln(writeErr.Error())
					}

					log.Println("File transfer successful: " + path.Join(c.DestinationDirectory, file.Name()))
				}
			}
		}
	}
}
