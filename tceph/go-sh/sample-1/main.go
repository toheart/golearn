package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

/**
@file:
@author: levi.Tang
@time: 2024/9/17 10:54
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
	out, _, err := PowershellOutput(fmt.Sprintf(`
		$disk=Get-Disk -SerialNumber %s 
		if (!$disk) {
            Write-Host false
		}else{
            Write-Host true
        }
	`, "runtime_pool_v61/sttcmnt.g4.v12110.vm101.t0.i0"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}
