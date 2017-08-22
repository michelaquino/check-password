package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"gitlab.globoi.com/michel.aquino/check-password/context"
	"gitlab.globoi.com/michel.aquino/check-password/repository"
)

func main() {
	log := context.GetLogger()

	unprocessCredentials, err := repository.GetUnprocessedCredentials()
	if err != nil {
		log.Error("Get unproccess credentials", "Error", err.Error())
		return
	}

	sha1Credentials := []string{}
	for _, credential := range unprocessCredentials {
		fmt.Println("MD5 hash: ", credential.PasswordMD5Hash)
		sha1Credentials = append(sha1Credentials, credential.PasswordMD5Hash)
	}

	if err := writeHashFile(sha1Credentials); err != nil {
		log.Error("Writing file with hash", "Error", err.Error())
	}

	output, err := exec.Command("hashcat", "-m0", "/tmp/test_file.txt", "/opt/generic/hashcat-3.6.0/all_dictionary.dic").Output()
	if err != nil {
		log.Error("Exec hashcat command", "Error", err.Error())
		return
	}

	log.Info("Exec hashcat command", "Success", fmt.Sprintf("Output: %s", output))
}

func writeHashFile(sha1Credentials []string) error {
	log := context.GetLogger()

	file, err := os.Create("/tmp/test_file.txt")
	if err != nil {
		log.Error("Error on create file", "Error", err.Error())
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range sha1Credentials {
		fmt.Fprintln(w, line)
	}

	if err := w.Flush(); err != nil {
		log.Error("Error on flush buffer", "Error", err.Error())
		return err
	}

	return nil
}
