package cheats

import "os"

var Enabled = false

func init() {
	value, ok := os.LookupEnv("CHEATS")
	Enabled = ok && value != "0"
}
