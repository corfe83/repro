package asset

import (
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/mobile/asset"
)

func validatePathWorksForCaseSensitiveOS(path string) {
	var directory string
	if filepath.IsAbs(path) {
		directory = filepath.Dir(path)
	} else {
		directory = "assets"
	}
	fileName := filepath.Base(path)

	entries, err := os.ReadDir(directory)
	if err != nil {
		panic("Failed to check case sensitivity")
	}

	for i := range entries {
		if entries[i].Name() == fileName {
			// We found exact filename, all is good
			return
		}

		if strings.EqualFold(entries[i].Name(), fileName) {
			// Exact filename was not found, this filename has wrong case. Throw error so we know to fix it, some OS's depend on this!
			panic("File " + path + " has wrong casing! Casing provided: \"" + fileName + "\", actual casing: \"" + entries[i].Name() + "\"")
		}
	}

	// The file was not found at all, that's OK we'll let the normal runtime error process handle this
}

func LoadImage(path string) (*ebiten.Image, error) {
	var err error
	f, err := asset.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	loadedPicture, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	var ebitenImage *ebiten.Image
	ebitenImage = ebiten.NewImageFromImage(loadedPicture)
	return ebitenImage, nil
}

func LoadFile(path string) (asset.File, error) {
	if runtime.GOOS == "android" {
		return asset.Open(path)
	} else {
		return os.Open(filepath.FromSlash("assets/" + path))
	}
}
