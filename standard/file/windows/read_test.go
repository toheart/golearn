package main

import (
	"os"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/4/12 18:10
@description:
**/

func TestOsRemove(t *testing.T) {
	err := os.RemoveAll("D:\\test\\H\\4")
	if err != nil {
		t.Errorf("%s \n", err)
	}
}
