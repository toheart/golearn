package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

/**
@file:
@author: levi.Tang
@time: 2024/9/16 18:15
@description:
**/

type CephMapped struct {
	Id          int    `json:"id"`
	Device      string `json:"device"`
	Pool        string `json:"pool"`
	Namespace   string `json:"namespace"`
	Image       string `json:"image"`
	Snap        string `json:"snap"`
	Disk_number int    `json:"disk_number"`
	Status      string `json:"status"`
}

func main() {
	cmd := exec.Command("rbd", "showmapped", "--format", "json")
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	if err := cmd.Run(); err != nil {
		log.Fatal(err.Error())
		return
	}
	var mapped []*CephMapped
	err := json.Unmarshal(sout.Bytes(), &mapped)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Printf("output: %+v", mapped)
}
