package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/michelaquino/check-password/context"
	"github.com/michelaquino/check-password/repository"
)

const (
	hashcatInputFile  = "/tmp/hashcat_intput.txt"
	hashcatOutputFile = "/tmp/hashcat_output.txt"
)

func main() {
	log := context.GetLogger()
	log.Debug("Hashcat runner", "", "Start")

	unprocessCredentials, err := repository.GetUnprocessedCredentials()
	if err != nil {
		log.Error("Get unproccess credentials", "Error", err.Error())
		os.Exit(1)
		return
	}

	if len(unprocessCredentials) == 0 {
		os.Exit(0)
		return
	}

	md5Credentials := []string{}
	for _, credential := range unprocessCredentials {
		md5Credentials = append(md5Credentials, credential.PasswordMD5Hash)
	}

	if err := writeHashFile(md5Credentials); err != nil {
		log.Error("Writing file with hash", "Error", err.Error())
		os.Exit(1)
		return
	}

	exec.Command("rm", "/tmp/hashcat_output.txt").Run()
	exec.Command("hashcat",
		"-m", "0",
		"-a", "0",
		"-r", "/opt/generic/hashcat-3.6.0/rules/best64.rule",
		"--potfile-disable",
		"-o", hashcatOutputFile,
		"--outfile-format=1",
		hashcatInputFile,
		"/opt/generic/hashcat-3.6.0/all_dictionary.dic").Run()

	md5HashListHacked := []string{}
	if _, err := os.Stat(hashcatOutputFile); !os.IsNotExist(err) {
		file, err := os.Open(hashcatOutputFile)
		if err != nil {
			log.Error("Open hashcat output file", "Error", err.Error())
			os.Exit(1)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			md5HashListHacked = append(md5HashListHacked, scanner.Text())
		}
	}

	repository.UpdateCredentialsProcessed(md5Credentials, md5HashListHacked)
	os.Exit(0)
}

func writeHashFile(md5Credentials []string) error {
	log := context.GetLogger()

	file, err := os.Create(hashcatInputFile)
	if err != nil {
		log.Error("Error on create file", "Error", err.Error())
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range md5Credentials {
		fmt.Fprintln(w, line)
	}

	if err := w.Flush(); err != nil {
		log.Error("Error on flush buffer", "Error", err.Error())
		return err
	}

	return nil
}
