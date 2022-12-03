package util

import (
	"log"
	"net/http"
	"os/exec"
	"time"
)

func delaySuspend(n time.Duration) {
	time.Sleep(n * time.Second)
	cmd := exec.Command("systemctl", "suspend")
	_, err := cmd.Output()
	if err != nil {
		log.Println("Error executing suspend command: ", err)
	}
}

func Suspend(w http.ResponseWriter) error {
	// call DLL to set system to suspend

	// write status ok to w
	w.WriteHeader(http.StatusOK)
	// write "System is suspending" to w
	w.Write([]byte("System is suspending"))
	go delaySuspend(3)
	return nil
}

func Reboot(bootnext string, w http.ResponseWriter) error {
	// call shell command to execute "efibootmgr --bootnext 0000" command
	if len(bootnext) > 0 {
		cmd := exec.Command("efibootmgr", "--bootnext", bootnext)
		_, err := cmd.Output()
		if err != nil {
			log.Println("Error executing efibootmgr command: ", err)
			return err
		}
	}
	// write status ok to w
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("System is rebooting in 1 minute"))

	// force flush of w
	w.(http.Flusher).Flush()

	// call shell command /usr/sbin/shutdown -t 15 -r
	cmd := exec.Command("shutdown", "-t", "1", "-r")
	// cmd = exec.Command("reboot")
	_, err := cmd.Output()
	if err != nil {
		log.Println("Error executing reboot command: ", err)
		return err
	}
	return nil
}
