[![Go Report Card](https://goreportcard.com/badge/github.com/OldTyT/alerta_notify)](https://goreportcard.com/report/OldTyT/alerta_notify)
[![GolangCI](https://golangci.com/badges/github.com/OldTyT/alerta_notify.svg)](https://golangci.com/r/github.com/OldTyT/alerta_notify)


# Alerta notify

---

This is simple desktop notify on alerts from [Alerta](https://github.com/alerta/alerta)

Default path to config - `$HOME/.config/alerta_notify.json`

### How to start

```
./alerta_notify -config="path/to/config"
```

### Config

* `alerta.username` - User name from [Alerta](https://github.com/alerta/alerta).
* `alerta.password` - User password from [Alerta](https://github.com/alerta/alerta).
* `alerta.url` - URL Address [Alerta](https://github.com/alerta/alerta).
* `alerta.query` - Request to [Alerta](https://github.com/alerta/alerta) by which alerts will be received.
* `path.icon_notify` - The path to the icon, for notification. It does not work on all operating systems.
* `path.icon_alert` - The path to the icon, for alerts notification. It does not work on all operating systems.
* `path.sound_notify` - Path to sound file for notify message. To turn off the sound, specify a non-existent file and ignore the errors.
* `path.sound_alert` - Path to sound file for alert.
* `time_sleep` - Sleep time between iterations, in seconds.
