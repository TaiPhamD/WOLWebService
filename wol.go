package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/linde12/gowol"
)

//import "strings"

var Mypassword string
var MyPort string
var MyMAC string

type auth_struct struct {
	Password string
}

type turnonHandler struct{}

func (h turnonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonAuth auth_struct

	err := decoder.Decode(&jsonAuth)
	if err != nil {
		log.Print("error decoding JSON\n")
		return
	}

	if jsonAuth.Password != Mypassword {
		log.Print("Password from JSON doesn't match\n")
		return
	}
	if packet, err := gowol.NewMagicPacket(MyMAC); err == nil {
		log.Print("Sending WOL magic packet")
		packet.Send("192.168.2.255") // send to broadcast
		// specify receiving port
	}
}

func main() {
	//Get file path from where the exe is launched
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Print(dir)

	//set up log file
	filelog, errlog := os.OpenFile(dir+"/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errlog != nil {
		log.Fatal(errlog)
	}

	defer filelog.Close()

	log.SetOutput(filelog)

	file, err := os.Open(dir + "/config.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	Mypassword = scanner.Text()
	scanner.Scan()
	MyPort = scanner.Text()
	scanner.Scan()
	MyMAC = scanner.Text()
	file.Close()

	log.Print("My password is:", Mypassword)
	log.Print("My port: ", MyPort)
	log.Print("My MAC is: ", MyMAC)
	err = http.ListenAndServe(":"+MyPort, turnonHandler{})

}
