package main

import (
	"image"
	_ "image/png"
	"strconv"
	"time"

	"github.com/corfe83/repro/shape"

	"github.com/corfe83/repro/asset"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 960
	screenHeight = 540
)

const circleImagesToCreate = 20
type Game struct {
	images []*ebiten.Image
	circleImages []*ebiten.Image
	baseCircle *ebiten.Image

	firstFrameProcessed bool
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) firstFrame() {
	for i := range imagesToLoad {
		img, err := asset.LoadPng(imagesToLoad[i])
		if err != nil {
			panic(err)
		}
		g.images = append(g.images, img)
	}
	g.baseCircle = ebiten.NewImageFromImage(baseCircleImg)
}

func (g *Game) ensureCircleImages() {
	if len(g.circleImages) >= circleImagesToCreate {
		return
	}
	newImage := ebiten.NewImage(g.baseCircle.Size())
	var op ebiten.DrawImageOptions
	op.ColorM.Scale(0.0, 0.2 + float64(len(g.circleImages))*0.01, 0.0, 1.0)
	newImage.DrawImage(g.baseCircle, &op)
	g.circleImages = append(g.circleImages, newImage)
	time.Sleep(time.Millisecond * 50)
}

const renderWidth = 60
const imagesPerRow = screenWidth / renderWidth
const renderHeight = 50
func (g *Game) Draw(screen *ebiten.Image) {
	if !g.firstFrameProcessed {
		g.firstFrame()
		g.firstFrameProcessed = true
	}
	allImages := append(g.images, g.circleImages...)
	g.ensureCircleImages()
	for i := range allImages {
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterLinear
		xOffset := (i % imagesPerRow) * renderWidth
		yOffset := (i / imagesPerRow) * renderHeight
		xSize, ySize := allImages[i].Size()
		scaleX := (renderWidth - 1) / float64(xSize)
		scaleY := (renderHeight - 1) / float64(ySize)
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(xOffset), float64(yOffset))
		screen.DrawImage(allImages[i], op)
	}
	for i := range imagesToLoad {
		xOffset := (i % imagesPerRow) * renderWidth
		yOffset := (i / imagesPerRow) * renderHeight
		modifiedString := stringModifier(imagesToLoad[i])
		text.Draw(screen, modifiedString, basicfont.Face7x13, xOffset, yOffset+15, colornames.White)
	}
	for i := range g.circleImages {
		index := i + len(imagesToLoad)
		xOffset := (index % imagesPerRow) * renderWidth
		yOffset := (index / imagesPerRow) * renderHeight
		var stringToWrite string
		if i < circleImagesToCreate - 1 {
			stringToWrite = "Circ " + strconv.Itoa(i+1)
		} else {
			stringToWrite = "Final circle!"
		}
		maxX := g.circleImages[i].Bounds().Max.X
		maxY := g.circleImages[i].Bounds().Max.Y
		stringToWrite += ":" + strconv.Itoa(maxX)
		stringToWrite += "," + strconv.Itoa(maxY)
		text.Draw(screen, stringModifier(stringToWrite), basicfont.Face7x13, xOffset, yOffset+15, colornames.Gray)
	}
}

const charsPerLine = 7
func stringModifier(input string) string {
	var result string
	var offsetSoFar int
	for len(input) - offsetSoFar > charsPerLine {
		result += input[offsetSoFar:offsetSoFar+charsPerLine]
		result += "\n"
		offsetSoFar += charsPerLine
	}
	if offsetSoFar < len(input) {
		result += input[offsetSoFar:]
	}
	return result
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var baseCircleImg image.Image
func main() {
	baseCircleImg = shape.CreateCircleImage(4000, 2000, shape.EdgeGradientInnerFocus)

	var g Game
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Android Image Invalidation Bug Repro")
	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}