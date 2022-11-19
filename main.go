package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/OldTyT/alerta_notify/internal/vars"
	"github.com/martinlindhe/notify"
)

func main() {
	ConfFile := flag.String("conf", "config.json", "Path to conf file.")
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
	UpdateAlerts()
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
		fmt.Println("Error auth in alerta.")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Response status != 200, exiting")
		notify.Alert("Alerta Notify", "Alerta Notify", "Response status != 200, exiting", "path/to/icon.png")
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &TokenLocal); err != nil {
		fmt.Println("Can't unmarshal JSON")
		fmt.Println(err)
		os.Exit(1)
	}
	vars.Other.AlertaToken = TokenLocal.Token
	notify.Notify("Alerta Notify", "Alerta Notify", "Success auth in Alerta", "path/to/icon.png")
}

func UpdateAlerts() {
	type AlertaAlertList struct {
		Total int `json:"total"`
	}
	var AlertsSummary AlertaAlertList
	URL := vars.Notifier.AlertaURL + "/alerts?sort-by=lastReceiveTime&status=open"
	for true {
		client := &http.Client{}
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			notify.Notify("Alerta Notify", "Alerta Notify", "Failure to get alerts.", "path/to/icon.png")
		}
		req.Header.Set("Authorization", "Bearer "+vars.Other.AlertaToken)
		resp, err := client.Do(req)
		if err != nil {
			notify.Notify("Alerta Notify", "Alerta Notify", "Failure to get alerts.", "path/to/icon.png")
		}
		if resp.StatusCode != 200 {
			fmt.Println("Response status != 200")
			notify.Alert("Alerta Notify", "Alerta Notify", "Response status != 200", "path/to/icon.png")
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &AlertsSummary); err != nil {
				fmt.Println("Can't unmarshal JSON")
				fmt.Println(err)
				os.Exit(1)
			}
			if AlertsSummary.Total != 0 {
				notify.Alert("Alerta Notify", "Alerta Notify", "Find active alerts!", "path/to/icon.png")
			}
		}
		time.Sleep(time.Duration(vars.Notifier.TimeSleep) * time.Second)
	}

}
