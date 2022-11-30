package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"wolwebservice/api"
	"wolwebservice/config"

	"github.com/kardianos/service"
)

var MyConfig config.Config

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// First need to create certs for HTTPS via LetsEncrypt service and certbot
	// sudo certbot certonly --standalone -d www.your_domain.com
	// Cert bot will create a fullchain and a privateg key file
	// So input these path in line 4 and 5 of config.txt
	var err error
	var handler http.Handler
	MyConfig, handler, err = api.Setup()
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}
	log.Print("Serving at port: ", MyConfig.Port)

	// check TLS mode
	if !MyConfig.TLS {
		err = http.ListenAndServe(":"+MyConfig.Port, handler)
	} else {
		log.Print("tls chain path: ", MyConfig.Fullchain)
		log.Print("tls key path: ", MyConfig.PrivKey)
		err = http.ListenAndServeTLS(":"+MyConfig.Port, MyConfig.Fullchain, MyConfig.PrivKey, handler)
	}
	if err != nil {
		log.Fatal("Error starting web server: ", err)
	}

}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func setupLog() *os.File {
	//Get file path from where the exe is launched
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Create new folder "logs" if it doesn't exist
	if _, err := os.Stat(dir + "/logs"); os.IsNotExist(err) {
		os.Mkdir(dir+"/logs", 0777)
	}

	t := time.Now().Format("2006-01-02_15_04_05")
	logPath := dir + "/logs/info_" + t + ".log"
	log.Print("Storing log file at location :" + logPath)

	//set up log file
	filelog, errlog := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errlog != nil {
		log.Fatal(errlog)
	}
	log.SetOutput(filelog)
	return filelog
}

func main() {

	filelog := setupLog()
	defer filelog.Close()

	svcConfig := &service.Config{
		Name:        "WOL_SERVICE",
		DisplayName: "WoL Service server",
		Description: "WoL Service server",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
