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

* `alerta_username` - User name from [Alerta](https://github.com/alerta/alerta)
* `alerta_password` - User password from [Alerta](https://github.com/alerta/alerta)
* `alerta_url` - URL Address [Alerta](https://github.com/alerta/alerta)
* `alert_query` - Request to [Alerta](https://github.com/alerta/alerta) by which alerts will be received
* `time_sleep` - Sleep time between iterations, in seconds
* `path_to_icon` - The path to the icon, for notification. It does not work on all operating systems.