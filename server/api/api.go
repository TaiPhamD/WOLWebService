package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"wolwebservice/config"
	"wolwebservice/util"
)

type Params struct {
	Alias  string `json:"alias"`
	APIKey string `json:"api_key"`
	Os     string `json:"os"`
}

var MyConfig config.Config

func Setup() (config.Config, *http.ServeMux, error) {
	var err error
	// parse config
	MyConfig, err = config.ParseConfig()
	mux := http.NewServeMux()
	// handle WOL requests
	if MyConfig.Master {
		mux.HandleFunc("/api/wol", WoLHandler)
	}
	// handle Restart requests
	mux.HandleFunc("/api/restart", RestartHandler)
	// handle OS query requests
	mux.HandleFunc("/api/os", OSQueryHandler)

	return MyConfig, mux, err
	//setup endpoints
}

func getSlaveIP(alias string) string {
	var slaveIp string
	slaveIp = ""
	for _, client := range MyConfig.Clients {
		if strings.ToLower(client.Alias) == strings.ToLower(alias) {
			slaveIp = client.IP
		}
	}
	return slaveIp
}

/** Forward exact request to slave**/
func forwardRequest(w http.ResponseWriter, r *http.Request, params Params) {
	slaveIp := getSlaveIP(params.Alias)
	if slaveIp == "" {
		log.Println("missing slave IP ")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
		return
	}

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s:%s%s", "http", slaveIp, MyConfig.Port, r.RequestURI)
	log.Println("request URL: ", url)

	//convert params to json
	jsonParams, err := json.Marshal(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(jsonParams))

	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// read resp body and return

	w.WriteHeader(resp.StatusCode)
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		w.Write(body)
	}
	defer resp.Body.Close()
}

/** Send WOL magic packet to slave **/
func sendWol(w http.ResponseWriter, params Params) {

	// find mac address from alias
	var mac string
	var ip string
	mac = ""
	ip = ""
	for _, client := range MyConfig.Clients {
		if client.Alias == params.Alias {
			mac = client.MAC
			ip = client.IP
		}
	}
	if mac == "" {
		log.Println("alias from JSON doesn't match")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
		return
	}

	log.Println("Found mac and ip for alias: " + params.Alias)
	err := util.SendWol(mac, ip)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Waking up" + params.Alias))

	} else {
		log.Println("Error creating magic packet: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error"))

	}
}

/** Reboot slave system **/
func reboot(w http.ResponseWriter, params Params) {
	// find BootId from params.Os
	log.Println("Received parameter master and os:", MyConfig.Master, params.Os)
	//Find BootID from params.Os
	var bootId string
	bootId = ""
	for _, client := range MyConfig.Clients {
		if client.Alias == params.Alias {
			for _, os := range client.Os {
				if strings.ToLower(os.Name) == strings.ToLower(params.Os) {
					bootId = os.BootID
				}
			}
		}
	}
	if bootId == "" {
		log.Println("alias from JSON doesn't match: " + params.Alias)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
		return
	}

	log.Println("Found bootId for alias: " + params.Alias + " and os: " + params.Os)
	// do the actual restart and change BootNext based on bootId
	err := util.Reboot(bootId, w)
	if err != nil {
		log.Println("Error rebooting (possibly permission issues): ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error"))
		return
	}

	//write good status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rebooting computer to: " + params.Os))
}

/** Get parameters from request and check api_key **/
func GetAuthParams(r *http.Request) (Params, error) {
	var params Params
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return Params{}, err
	}
	if params.APIKey != MyConfig.APIKey {
		return Params{}, errors.New("APIKey from JSON doesn't match")
	}
	return params, nil
}

func WoLHandler(w http.ResponseWriter, r *http.Request) {
	params, err := GetAuthParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	log.Println("Received WoL request for: ", params.Alias)
	if params.Alias != "" {
		//no os specified, send wol
		sendWol(w, params)
	} else {
		log.Println("mode from JSON doesn't match")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
	}

}

func RestartHandler(w http.ResponseWriter, r *http.Request) {
	params, err := GetAuthParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	if !MyConfig.Master {
		reboot(w, params)
	} else {
		log.Println("Forwarding Restart request for: ", params.Alias)
		forwardRequest(w, r, params)
	}
}

func OSQueryHandler(w http.ResponseWriter, r *http.Request) {
	params, err := GetAuthParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	if !MyConfig.Master {
		// If slave then process request
		w.WriteHeader(http.StatusOK)
		// determine if OS is linux or windows
		currentOS := runtime.GOOS
		w.Write([]byte(currentOS))

	} else {
		log.Println("Forwarding OSQuery request for: ", params.Alias)
		forwardRequest(w, r, params)
	}
}
