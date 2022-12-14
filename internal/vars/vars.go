package vars

type NotifierCFG struct {
	Alerta    AlertaStuct `json:"alerta"`
	Path      PathStruct  `json:"path"`
	TimeSleep int         `json:"time_sleep"`
}

type PathStruct struct {
	Icon  NotifyAlert `json:"icon"`
	Sound NotifyAlert `json:"sound"`
}

type NotifyAlert struct {
	Notify string `json:"notify"`
	Alert  string `json:"alert"`
}

type AlertaStuct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Query    string `json:"query"`
}

type OtherCFG struct {
	AlertaToken string
}

var (
	Notifier NotifierCFG
	Other    OtherCFG
	Version  string = "v0.0.7"
)
