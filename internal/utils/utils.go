package utils

import (
	"os"
	"strconv"

	"github.com/OldTyT/alerta_notify/internal/vars"
	"github.com/OldTyT/notify"
)

func ViewSummary() {
	URL := vars.Notifier.Alerta.URL + vars.Notifier.Alerta.Query
	SendNotify("Alerta query: " + URL + "\nSleep time: " + strconv.Itoa(vars.Notifier.TimeSleep) + "sec" + "\nVersion: " + vars.Version)
}

func ErrorExiting(ErrorMsg string) {
	notify.Alert("Alerta notify", "Alerta Notify", ErrorMsg, vars.Notifier.Path.Icon.Alert, vars.Notifier.Path.Sound.Alert)
	os.Exit(1)
}

func SendNotify(text string) {
	go notify.Notify("Alerta notify", "Alerta Notify", text, vars.Notifier.Path.Icon.Notify, vars.Notifier.Path.Sound.Notify)
}

func SendAlert(text string) {
	go notify.Alert("Alerta notify", "Alerta Notify", text, vars.Notifier.Path.Icon.Alert, vars.Notifier.Path.Sound.Alert)
}
