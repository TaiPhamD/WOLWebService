package util

import (
	"net/http"
)

func Suspend(w http.ResponseWriter) error {
	// TO DO  set system to suspend in macOS
	return nil
}
func Reboot(bootnext string, w http.ResponseWriter) error {
	// TO DO  change BootNext Variable in macOS API some how
	return nil
}
