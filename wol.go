package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)
import "github.com/linde12/gowol"

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
		fmt.Println("Sending WOL magic packet")
		packet.Send("192.168.2.255") // send to broadcast
		// specify receiving port
	}
}

func main() {

	file, err := os.Open("config.txt")
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

	fmt.Println("My password is:", Mypassword)
	fmt.Println("My port: ", MyPort)
	fmt.Println("My MAC is: ", MyMAC)
	err = http.ListenAndServe(":"+MyPort, turnonHandler{})

}
