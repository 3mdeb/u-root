// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a basic init script.
package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	commands = []string{
		"/bbin/mount -t ext4 /dev/sda1 /var",
		"/bbin/kexec -l -c \"dom0_mem=512M loglvl=all guest_loglvl=all com1=115200,8n1 console=com1 no-real-mode\" --module \"/var/bzImage console=hvc0 earlyprintk=xen nomodeset root=/dev/sda2\" /var/xen.gz",
		"/bbin/kexec -e",
		"/bbin/shutdown halt",
	}
)

func main() {
	for _, line := range commands {
		log.Printf("Executing Command: %v", line)
		cmdSplit := strings.Split(line, " ")
		if len(cmdSplit) == 0 {
			continue
		}

		cmd := exec.Command(cmdSplit[0], cmdSplit[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}

	}
	log.Print("Uinit Done!")
}
