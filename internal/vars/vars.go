package vars

type NotifierCFG struct {
	Alerta    AlertaStuct `json:"alerta"`
	Path      PathStruct  `json:"path"`
	TimeSleep int         `json:"time_sleep"`
}

type PathStruct struct {
	IconNotify  string `json:"icon_notify"`
	IconAlert   string `json:"icon_alert"`
	SoundNotify string `json:"sound_notify"`
	SoundAlert  string `json:"sound_alert"`
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
	Version  string = "v0.0.6"
)
