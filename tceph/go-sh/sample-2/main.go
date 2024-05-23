package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

/**
@file:
@author: levi.Tang
@time: 2024/9/17 11:17
@description:
**/

func PowershellOutput(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("powershell.exe", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func main() {
	output, errout, err := PowershellOutput(
		fmt.Sprintf(`
		$disk=Get-Disk -SerialNumber %s
		if ($disk.IsOffline) {
			Set-Disk -Number $disk.Number -IsOffline $False
		}
		$part=Get-Partition -DiskNumber $disk.Number | Where-Object {$_.Type -Eq "Basic"}
		if (!$part) {
			Write-Error "no partition available"
 			exit 3
		}
		$output=$part.AccessPaths -join ","
		Write-Output $output
`, os.Args[1]))
	if err != nil {
		log.Fatal("exec GetCephDriverLetter found errout: %s, err: %s", errout, err)
	}
	var MountPoint string
	for _, p := range strings.Split(output, ",") {
		if strings.HasPrefix(p, `\\`) {
			continue
		}
		MountPoint = strings.TrimSpace(p)
		break
	}
	fmt.Printf("mountPoint: %s", MountPoint)
}
