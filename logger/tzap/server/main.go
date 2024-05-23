package main

import (
	"bytes"
	log_template "github.com/toheart/golearn/logger/log-template"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"math/rand"
	"net/http"
	"time"
)

/**
@file:
@author: levi.Tang
@time: 2024/10/21 12:18
@description:
**/

var logger *zap.Logger
var maxContent string

func init() {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	// 定义可能的字符
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// 30MB大小的字符串，即30,048,576个字符
	const size = 30 * 1024 * 1024
	// 使用bytes.Buffer来高效构建字符串
	var buffer bytes.Buffer

	for i := 0; i < size; {
		b := chars[rand.Intn(len(chars))]
		buffer.WriteRune(b)
		i++
	}
	// 将bytes.Buffer转换为字符串
	maxContent = buffer.String()
}

func main() {
	//Create the default mux
	mux := http.NewServeMux()
	//Handling the /v1/teachers. The handler is a function here
	mux.HandleFunc("/v1/largeLog", LargeLogHandler)
	mux.HandleFunc("/v1/noLog", NoLogHandler)
	//Handling the /v1/students. The handler is a type implementing the Handler interface here
	sHandler := smallLogHandler{}
	mux.Handle("/v1/smallLog", sHandler)
	logger = log_template.InitLogger(&lumberjack.Logger{
		Filename:  "./app.log",
		Compress:  true,
		LocalTime: true,
	})
	//Create the server.
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()

}

func NoLogHandler(res http.ResponseWriter, req *http.Request) {
	data := []byte("no log handler")
	res.WriteHeader(200)
	res.Write(data)
}

func LargeLogHandler(res http.ResponseWriter, req *http.Request) {
	logger.Info("test", zap.String("content", maxContent))
	data := []byte("large log handler")
	res.WriteHeader(200)
	res.Write(data)
}

type smallLogHandler struct{}

func (h smallLogHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logger.Info("input small log")
	data := []byte("small log handler")
	res.WriteHeader(200)
	res.Write(data)
}
