package rewind

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	gameserverpb "github.com/c4t-but-s4d/ctfcup-2024-igra/proto/go/gameserver"
)

const rewindDir = "rewind"

type Rewind struct {
	Moves []*gameserverpb.ClientEvent_KeysPressed `json:"moves"`
}

func filterKeys(keys []int32) []int32 {
	keys = lo.Filter(keys, func(k int32, _ int) bool {
		return k != int32(ebiten.KeyX) && k != int32(ebiten.KeyEqual) && k != int32(ebiten.KeyMinus)
	})
	return keys
}

func (r *Rewind) Record(move *gameserverpb.ClientEvent_KeysPressed) {
	move.KeysPressed = filterKeys(move.KeysPressed)
	move.NewKeysPressed = filterKeys(move.NewKeysPressed)
	r.Moves = append(r.Moves, move)
}

func (r *Rewind) SaveAndReport() {
	if path, err := r.Save(); err != nil {
		logrus.Errorf("failed to save rewind: %v", err)
	} else {
		logrus.Infof("rewind saved to %s", path)
	}
}

func (r *Rewind) Save() (string, error) {
	path := generateFilename()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("creating directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(r); err != nil {
		return "", fmt.Errorf("encoding rewind: %w", err)
	}
	return path, nil
}

func (r *Rewind) Clear() {
	r.Moves = nil
}

func (r *Rewind) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(r); err != nil {
		return fmt.Errorf("decoding rewind: %w", err)
	}

	return nil
}

func generateFilename() string {
	ts := time.Now().Format("2006-01-02_15-04-05")
	return filepath.Join(rewindDir, fmt.Sprintf("%s.json", ts))
}
