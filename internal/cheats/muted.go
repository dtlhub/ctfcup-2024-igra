package cheats

import "os"

var Muted = false

func init() {
	value, ok := os.LookupEnv("MUTED")
	Muted = ok && value != "0"
}
