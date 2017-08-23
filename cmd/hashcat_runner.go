package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"gitlab.globoi.com/michel.aquino/check-password/context"
	"gitlab.globoi.com/michel.aquino/check-password/repository"
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

	file, err := os.Open(hashcatOutputFile)
	if err != nil {
		log.Error("Open hashcat output file", "Error", err.Error())
		os.Exit(1)
		return
	}
	defer file.Close()

	md5HashListHacked := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		md5HashListHacked = append(md5HashListHacked, scanner.Text())
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
