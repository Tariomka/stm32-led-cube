package global

import "time"

var HandleFatal = func() {
	for {
		time.Sleep(100 * time.Second)
	}
}

var HandleError = func() {}

var HandleStartup = func() {}
