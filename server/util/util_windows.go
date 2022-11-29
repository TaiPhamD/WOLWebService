package util

import (
	"log"
	"net/http"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

func delaySuspend(n time.Duration) {
	time.Sleep(n * time.Second)
	loaddll := syscall.MustLoadDLL("efiDLL")
	//defer loaddll.Release()
	SystemSuspendFunc := loaddll.MustFindProc("SystemSuspend")
	SystemSuspendFunc.Call()
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
	// call DLL to change boot NEXT

	var mode uint16
	// Mode = 0 to change BootNext variable
	mode = 0
	//decode bootnext to uint16
	data, err := strconv.ParseInt(bootnext, 16, 16)
	if err != nil {
		log.Print(err)
		return err
	}

	loaddll := syscall.MustLoadDLL("efiDLL")
	//defer loaddll.Release()
	ChangeBootFunc := loaddll.MustFindProc("SystemChangeBoot")
	ChangeBootFunc.Call(uintptr(unsafe.Pointer(&data)), uintptr(unsafe.Pointer(&mode)))

	// write status ok to w
	w.WriteHeader(http.StatusOK)
	// call DLL command to reboot
	// Mode = 0 is restart
	mode = 0
	ShutdownFunc := loaddll.MustFindProc("SystemShutdown")
	ShutdownFunc.Call(uintptr(unsafe.Pointer(&mode)))
	if err != nil {
		log.Println("Error executing reboot command: ", err)
		return err
	}
	return nil
}