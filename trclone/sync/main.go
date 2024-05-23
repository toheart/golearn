package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"time"
)

/**
@file:
@author: levi.Tang
@time: 2024/8/6 19:50
@description:
**/

func main() {
	defer func(start time.Time) {
		fmt.Printf("timecost: %s \n", time.Now().Sub(start).String())
	}(
		time.Now())
	// 尝试删除老文件
	src := "F:\\9105"
	dest := "E:\\9105"
	cmd := exec.Command("rclone.exe", "sync", src, dest,
		"--inplace", "--delete-before", "--retries=1")

	cmdReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("open stderr: %s", err)
		return
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s \n", err)
		return
	}

	fmt.Printf("exec rclone from: %s, to: %s sync success \n", src, dest)
	return
}

func gbk_to_utf8(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	if d, e := ioutil.ReadAll(reader); e == nil {
		return d
	}
	return s
}
