package vars

type NotifierCFG struct {
	AlertaUsername string `json:"alerta_username"`
	AlertaPassword string `json:"alerta_password"`
	AlertaURL      string `json:"alerta_url"`
	TimeSleep      int    `json:"time_sleep"`
}

type OtherCFG struct {
	AlertaToken string
}

var (
	Notifier NotifierCFG
	Other    OtherCFG
)
