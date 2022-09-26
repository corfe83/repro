package asset

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var cachedImages map[string]*ebiten.Image

func init() {
	cachedImages = make(map[string]*ebiten.Image)
}

func LoadPng(path string) (*ebiten.Image, error) {
	image, ok := cachedImages[path]
	if ok {
		return image, nil
	}

	var ebitenImage, err = LoadImage(path)
	if err != nil {
		return nil, err
	}

	cachedImages[path] = ebitenImage

	return ebitenImage, nil
}
