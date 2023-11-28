package config

var ssoEnabled bool = false

func IsSsoEnabled() bool {
	return ssoEnabled
}

func EnableSso() {
	ssoEnabled = true
}

func DisableSso() {
	ssoEnabled = false
}
