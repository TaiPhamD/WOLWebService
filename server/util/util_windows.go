package util

// #cgo LDFLAGS: -L../../build/ -lutil_windows  -lpowrprof -lstdc++
import (
	"C"
	"log"
	"net/http"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

func delaySuspend(n time.Duration) {
	time.Sleep(n * time.Second)
	// call C code to suspend
	C.SystemSuspend()
}

func Suspend(w http.ResponseWriter) error {
	// call DLL to set system to suspend
	// write status ok to w
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("System is suspending"))
	go delaySuspend(3)
	return nil
}

func Reboot(bootnext string, w http.ResponseWriter) error {
	// call C code to change UEFI boot order and reboot

	var mode uint16
	mode = 0 // Mode = 0 to change BootNext variable
	//decode bootnext to uint16
	data, err := strconv.ParseInt(bootnext, 16, 16)
	if err != nil {
		log.Print(err)
		w.writeHeader(http.StatusBadRequest)
		return err
	}

	c_mode := C.uint16_t(mode)
	c_data := C.uint16_t(data)
	// call C code to change bootnext
	C.SystemChangeBoot(c_mode, c_data)

	// Call C code to reboot
	w.writeHeader(http.StatusOK)
	w.Write([]byte("System is rebooting"))
	mode = 0 // Mode = 0 to reboot
	cmode = C.uint16_t(mode)
	C.SystemShutdown(cmode)
	return nil
}
