package main

import (
	"net/http"
	"fmt"
	"os"
	"bufio"
	"log"
)
import "github.com/linde12/gowol"
//import "strings"

var Mypassword string
var MyPort string
var MyMAC string
type turnonHandler struct{}

func (h turnonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pass, _ := r.BasicAuth()
	fmt.Println(pass)
	if(pass != Mypassword){
		fmt.Println("Password Doesnt match\n")
		return
	}
       if packet, err := gowol.NewMagicPacket(MyMAC); err == nil {
                packet.Send("192.168.255.255")          // send to broadcast
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

	fmt.Println("My password is:",Mypassword)

        err = http.ListenAndServe(":"+MyPort, turnonHandler{})

}
