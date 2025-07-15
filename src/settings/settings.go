package settings

import "sync"

var settings *Settings
var once sync.Once

type Settings struct {
	ApiPort string
}

func GetSettings() *Settings {
	once.Do(func() {
		settings = &Settings{}
	})

	return settings
}
