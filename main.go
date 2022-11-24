package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/OldTyT/alerta_notify/internal/vars"
	"github.com/OldTyT/notify"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	ConfFile := flag.String("config", homedir+"/.config/alerta_notify.json", "Path to conf file.")
	flag.Parse()
	file, err := os.Open(*ConfFile)
	if err != nil {
		ErrorMsg := err.Error()
		fmt.Println(ErrorMsg)
		notify.Alert("Alerta notify", "Alerta Notify", ErrorMsg, "None", "default")
		ErrorMsg = "Error open conf file -" + *ConfFile
		fmt.Println(ErrorMsg)
		notify.Alert("Alerta notify", "Alerta Notify", ErrorMsg, "None", "default")
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&vars.Notifier)
	if err != nil {
		ErrorMsg := "error:" + err.Error()
		fmt.Println(ErrorMsg)
		notify.Alert("Alerta notify", "Alerta Notify", ErrorMsg, "None", "default")
		os.Exit(1)
	}
	go ViewSummary()
	LoginAlerta()
	for {
		go UpdateAlerts()
		time.Sleep(time.Duration(vars.Notifier.TimeSleep) * time.Second)
	}
}

func ViewSummary() {
	URL := vars.Notifier.Alerta.URL + vars.Notifier.Alerta.Query
	SendNotify("Alerta query: " + URL + "\nSleep time: " + strconv.Itoa(vars.Notifier.TimeSleep) + "sec" + "\nVersion: " + vars.Version)
}

func ErrorExiting(ErrorMsg string) {
	notify.Alert("Alerta notify", "Alerta Notify", ErrorMsg, vars.Notifier.Path.Icon, vars.Notifier.Path.SoundAlert)
	os.Exit(1)
}

func LoginAlerta() {
	type AlertaAuth struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	type AlertaToken struct {
		Token string `json:"token"`
	}
	var (
		AlertaAuthData AlertaAuth
		TokenLocal     AlertaToken
	)
	AlertaAuthData.Password = vars.Notifier.Alerta.Password
	AlertaAuthData.UserName = vars.Notifier.Alerta.Username
	JsonData, err := json.Marshal(AlertaAuthData)
	if err != nil {
		fmt.Println(err)
	}
	URL := vars.Notifier.Alerta.URL + "/auth/login"
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(JsonData))
	if err != nil {
		ErrorExiting("Error auth in alerta. " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ErrorExiting("Response status != 200, when authorization in alerta.\nExit.")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorExiting("Can't read JSON: " + err.Error())
	}
	if err := json.Unmarshal(body, &TokenLocal); err != nil {
		ErrorExiting("Can't unmarshal JSON: " + err.Error())
	}
	vars.Other.AlertaToken = TokenLocal.Token
	SendNotify("Successful authorization in Alerta")
}

func SendNotify(text string) {
	go notify.Notify("Alerta notify", "Alerta Notify", text, vars.Notifier.Path.Icon, vars.Notifier.Path.SoundNotify)
}

func SendAlert(text string) {
	go notify.Alert("Alerta notify", "Alerta Notify", text, vars.Notifier.Path.Icon, vars.Notifier.Path.SoundAlert)
}

func UpdateAlerts() {
	type AlertaAlertList struct {
		Alerts []map[string]interface{} `json:"alerts"`
		Total  int                      `json:"total"`
	}
	type AlertSummary struct {
		AlertName string `json:"event"`
		Resource  string `json:"resource"`
		ENV       string `json:"environment"`
		Severity  string `json:"severity"`
	}
	var (
		AlertsSummary AlertaAlertList
		Alert         []AlertSummary
	)
	URL := vars.Notifier.Alerta.URL + vars.Notifier.Alerta.Query
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		SendNotify("Error when receiving alerts.")
	}
	req.Header.Set("Authorization", "Bearer "+vars.Other.AlertaToken)
	resp, err := client.Do(req)
	if err != nil {
		SendNotify("Error when receiving alerts.")
	}
	if resp.StatusCode != 200 {
		SendAlert("Response status != 200, when update alerts.")
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			SendAlert("Can't read JSON, when update alerts. " + err.Error())
		}
		if err = json.Unmarshal(body, &AlertsSummary); err != nil {
			SendAlert("Can't unmarshal JSON, when update alerts. " + err.Error())
		}
		if AlertsSummary.Total != 0 {
			b, err := json.Marshal(AlertsSummary.Alerts)
			if err != nil {
				SendAlert("Can't marshal JSON. " + err.Error())
			}
			if err = json.Unmarshal(b, &Alert); err != nil {
				SendAlert("Can't unmarshal JSON. " + err.Error())
			}
			for key := range Alert {
				SendAlert("Can't unmarshal JSON. " + err.Error())
				go notify.Alert("Alerta notify", Alert[key].ENV+"/"+Alert[key].Severity, Alert[key].AlertName+"\n"+Alert[key].Resource, vars.Notifier.Path.Icon, vars.Notifier.Path.SoundAlert)
			}
		}
	}
}
