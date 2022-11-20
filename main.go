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
	"github.com/martinlindhe/notify"
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
		fmt.Println(err)
		fmt.Println("Error open conf file -", *ConfFile)
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&vars.Notifier)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	LoginAlerta()
	go ViewSummary()
	UpdateAlerts()
}

func ViewSummary() {
	URL := vars.Notifier.AlertaURL + vars.Notifier.AlertaQuery
	notify.Notify("Alerta notify", "Alerta notify summary", "Alerta query: "+URL+"\nSleep time: "+strconv.Itoa(vars.Notifier.TimeSleep)+"sec", vars.Notifier.PathIcon)
}

func ErrorExiting(ErrorMsg string) {
	notify.Alert("Alerta notify", "Alerta Notify", ErrorMsg, vars.Notifier.PathIcon)
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
	AlertaAuthData.Password = vars.Notifier.AlertaPassword
	AlertaAuthData.UserName = vars.Notifier.AlertaUsername
	JsonData, err := json.Marshal(AlertaAuthData)
	if err != nil {
		fmt.Println(err)
	}
	URL := vars.Notifier.AlertaURL + "/auth/login"
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(JsonData))
	if err != nil {
		ErrorExiting("Error auth in alerta. " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ErrorExiting("Response status != 200, when authorization in alerta.\nExit.")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &TokenLocal); err != nil {
		ErrorExiting("Can't unmarshal JSON: " + err.Error())
	}
	vars.Other.AlertaToken = TokenLocal.Token
	go notify.Notify("Alerta notify", "Alerta Notify", "Successful authorization in Alerta", vars.Notifier.PathIcon)
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
	URL := vars.Notifier.AlertaURL + vars.Notifier.AlertaQuery
	for true {
		client := &http.Client{}
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			go notify.Notify("Alerta notify", "Alerta Notify", "Error when receiving alerts.", vars.Notifier.PathIcon)
		}
		req.Header.Set("Authorization", "Bearer "+vars.Other.AlertaToken)
		resp, err := client.Do(req)
		if err != nil {
			go notify.Notify("Alerta notify", "Alerta Notify", "Error when receiving alerts.", vars.Notifier.PathIcon)
		}
		if resp.StatusCode != 200 {
			go notify.Alert("Alerta notify", "Alerta Notify", "Response status != 200, when update alerts.", vars.Notifier.PathIcon)
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &AlertsSummary); err != nil {
				go notify.Alert("Alerta notify", "Alerta Notify", "Can't unmarshal JSON, when update alerts. "+err.Error(), vars.Notifier.PathIcon)
			}
			if AlertsSummary.Total != 0 {
				b, err := json.Marshal(AlertsSummary.Alerts)
				if err != nil {
					go notify.Alert("Alerta notify", "Alerta Notify", "Can't marshal JSON. "+err.Error(), vars.Notifier.PathIcon)
				}
				if err = json.Unmarshal(b, &Alert); err != nil {
					go notify.Alert("Alerta notify", "Alerta Notify", "Can't unmarshal JSON. "+err.Error(), vars.Notifier.PathIcon)
				}
				for key := range Alert {
					go notify.Alert("Alerta notify", Alert[key].ENV+"/"+Alert[key].Severity, Alert[key].AlertName+"\n"+Alert[key].Resource, vars.Notifier.PathIcon)
				}
			}
		}
		time.Sleep(time.Duration(vars.Notifier.TimeSleep) * time.Second)
	}

}
