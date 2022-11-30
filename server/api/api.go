package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"wolwebservice/config"
	"wolwebservice/util"

	"golang.org/x/time/rate"
)

type Params struct {
	Alias  string `json:"alias"`
	APIKey string `json:"api_key"`
	Os     string `json:"os"`
}

var MyConfig config.Config

type ctxWoLParam struct{}

var limiter = rate.NewLimiter(1, 1)

// define limit middleware that checks limiter and also check decoded Params password
func LimitAndAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if request is allowed by limiter
		if !limiter.Allow() {
			// log too many requests and api path
			log.Println(r.URL.Path, "Too many requests")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("429 - Too Many Requests"))
			return
		}
		//check password
		var params Params
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Println(r.URL.Path, "Error decoding JSON: ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
			return
		}
		paramApiHash := sha256.Sum256([]byte(params.APIKey))
		if !bytes.Equal(paramApiHash[:], MyConfig.APIKeyHash[:]) {
			log.Println(r.URL.Path, "API Key doesn't match")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
		}
		// add param to r context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxWoLParam{}, params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Setup() (config.Config, http.Handler, error) {
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
	// handle Suspend requests
	mux.HandleFunc("/api/suspend", SuspendHandler)
	// handle OS query requests
	mux.HandleFunc("/api/os", OSQueryHandler)

	// handle all requests through LimitAndAuth middleware
	handler := LimitAndAuth(mux)

	return MyConfig, handler, err
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

func WoLHandler(w http.ResponseWriter, r *http.Request) {
	// get params from context
	params := r.Context().Value(ctxWoLParam{}).(Params)

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
	// get params from context
	params := r.Context().Value(ctxWoLParam{}).(Params)
	if !MyConfig.Master {
		reboot(w, params)
	} else {
		log.Println("Forwarding Restart request for: ", params.Alias)
		forwardRequest(w, r, params)
	}
}

func SuspendHandler(w http.ResponseWriter, r *http.Request) {
	// get params from context
	params := r.Context().Value(ctxWoLParam{}).(Params)
	if !MyConfig.Master {
		util.Suspend(w)
	} else {
		log.Println("Forwarding Suspend request for: ", params.Alias)
		forwardRequest(w, r, params)
	}
}

func OSQueryHandler(w http.ResponseWriter, r *http.Request) {
	// get params from context
	params := r.Context().Value(ctxWoLParam{}).(Params)
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
