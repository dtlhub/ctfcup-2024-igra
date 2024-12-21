package cheats

import (
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

func SetMonitor() {
	desiredMonitor := os.Getenv("MONITOR")
	if desiredMonitor == "" {
		return
	}

	monitors := ebiten.AppendMonitors(nil)
	names := make([]string, len(monitors))

	for i, m := range monitors {
		names[i] = m.Name()
		found := strings.Contains(m.Name(), desiredMonitor)
		if found {
			ebiten.SetMonitor(m)
			return
		}
	}

	logrus.Warnf("Monitor %s not found, available monitors: %s", desiredMonitor, strings.Join(names, ", "))
}
