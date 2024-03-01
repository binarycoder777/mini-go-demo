package main

import (
	_ "github.com/binarycoder777/mini-go-demo/demo/searchInfo/matchers"
	"github.com/binarycoder777/mini-go-demo/demo/searchInfo/search"
	"log"
	"os"
)

// init在main之前调用
func init() {
	// 日志输出到标准输出
	log.SetOutput(os.Stdout)
}

// 程序入口
func main() {
	search.Run("president")
}
