package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

type Config struct {
	ServerHost string
}

var config *Config

func init() {
	config = &Config{}
	flag.StringVar(&config.ServerHost, "p", "0.0.0.0:80", "local server address")
	flag.Usage = func() {
		fmt.Println(fmt.Sprintf("Usage of %s:\n", os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
}
func main() {
	var mux = http.NewServeMux()
	//加载http server file
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./dist/learning-angular"))))
	go func(serverAddr string, m *http.ServeMux) {
		if err := http.ListenAndServe(serverAddr, m); err != nil {
			fmt.Println("Cant start server:", err)
			os.Exit(1)
		}
	}(config.ServerHost, mux)
	// 捕捉ctrl+C 退出信号
	handleSignals()
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
