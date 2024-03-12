package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Host         string `toml:"host"`
	Port         uint   `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	IdentityFile string `toml:"identity_file"`
}

func main() {
	var (
		host         string
		port         uint
		user         string
		password     string
		identityFile string
		configPath   string
	)
	hostUsage := "the target host (required if no config file)"
	portUsage := "the port to connect"
	portDefualt := uint(22)
	userUsage := "the login user (required if no config file)"
	passwordUsage := "the login password"
	identityFileUsage := "the identity file"
	configPathUsage := "the path of config file (ignore other args if a config file exists)"
	configPathDefualt := "./config.toml"

	flag.StringVar(&host, "t", "", hostUsage)
	flag.UintVar(&port, "p", portDefualt, portUsage)
	flag.StringVar(&user, "u", "", userUsage)
	flag.StringVar(&password, "s", "", passwordUsage)
	flag.StringVar(&identityFile, "i", "", identityFileUsage)
	flag.StringVar(&configPath, "c", configPathDefualt, configPathUsage)

	flag.Parse()

	var cfg config
	var handler *sshHandler
	if _, err := toml.DecodeFile(configPath, &cfg); errors.Is(err, os.ErrNotExist) {
		if host == "" {
			log.Fatal("host can not be empty")
		}
		if user == "" {
			log.Fatal("user can not be empty")
		}
		if password == "" && identityFile == "" {
			log.Fatal("password can not be empty")
		}
		addr := fmt.Sprintf("%s:%d", host, port)
		handler = &sshHandler{addr: addr, user: user, secret: password}
	} else if err != nil {
		log.Fatal("could not parse config file: ", err)
	} else {
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		if cfg.Password != "" {
			handler = &sshHandler{addr: addr, user: cfg.User, secret: cfg.Password}
		} else {
			handler = &sshHandler{addr: addr, user: cfg.User, keyfile: cfg.IdentityFile}
		}

	}

	http.Handle("/", http.FileServer(http.Dir("./front/")))
	http.HandleFunc("/web-socket/ssh", cors(handler.webSocket))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息，注意前端也要设置withCredentials: true
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                                              //返回数据格式是json

		f(w, r)
	}
}
