package shape

import (
	"image"
	"image/color"
	"math"
)

type CircleStyle int

const (
	FilledIn CircleStyle = iota
	EdgeGradientInnerFocus
	EdgeGradientOuterFocus
	EdgeNoGradient
)

type quarterReflection struct {
	quarterImage     image.Image
	bounds           image.Rectangle
	midX             int
	midY             int
	reflectionOffset int
}

func createQuarterReflection(img image.Image, oddSize bool) quarterReflection {
	var edgeLength, midX, midY, reflectionOffset int
	if !oddSize {
		edgeLength = img.Bounds().Max.X * 2
		midX = img.Bounds().Max.X
		midY = img.Bounds().Max.X
		reflectionOffset = 1
	} else {
		edgeLength = img.Bounds().Max.X*2 - 1
		midX = img.Bounds().Max.X
		midY = img.Bounds().Max.X
		reflectionOffset = 2
	}
	newBounds := image.Rectangle{image.Pt(0, 0), image.Pt(edgeLength, edgeLength)}
	result := quarterReflection{
		quarterImage:     img,
		bounds:           newBounds,
		midX:             midX,
		midY:             midY,
		reflectionOffset: reflectionOffset,
	}
	return result
}

func (me quarterReflection) ColorModel() color.Model {
	return me.quarterImage.ColorModel()
}

func (me quarterReflection) Bounds() image.Rectangle {
	return me.bounds
}

func (me quarterReflection) At(x, y int) color.Color {
	if x < me.midX && y < me.midY {
		return me.quarterImage.At(x, y)
	} else if x >= me.midX && y < me.midY {
		return me.quarterImage.At(me.midX*2-x-me.reflectionOffset, y)
	} else if x < me.midX && y >= me.midY {
		return me.quarterImage.At(x, me.midY*2-y-me.reflectionOffset)
	} else {
		return me.quarterImage.At(me.midX*2-x-me.reflectionOffset, me.midY*2-y-me.reflectionOffset)
	}
}

// Circle will fill this imageWidth
func CreateCircleImage(imageEdgeLength int, edgeThickness float64, style CircleStyle) image.Image {
	// If it's an odd number, the quarter image needs to cover that middle pixel
	edgeLengthOfQuarterImage := imageEdgeLength/2 + imageEdgeLength%2

	quarterImage := image.NewNRGBA(image.Rectangle{image.Pt(0, 0), image.Pt(edgeLengthOfQuarterImage, edgeLengthOfQuarterImage)})

	// Radius of the full circle, including anti-aliased soft edge
	outerRadius := float64(imageEdgeLength) / 2.0
	// Radius of only the portion that is fully shaded in
	innerRadius := outerRadius - edgeThickness

	outerRadiusSquared := outerRadius * outerRadius
	innerRadiusSquared := innerRadius * innerRadius

	center := float64(imageEdgeLength) / 2.0

	color := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < edgeLengthOfQuarterImage; i++ {
		x := float64(i) + 0.5
		for j := i; j < edgeLengthOfQuarterImage; j++ {
			y := float64(j) + 0.5

			diffX := center - x
			diffY := center - y
			distanceFromCenterSquared := diffX*diffX + diffY*diffY
			if distanceFromCenterSquared < innerRadiusSquared {
				if style == FilledIn {
					color.A = 0xFF
				} else {
					continue
				}
			} else if distanceFromCenterSquared < outerRadiusSquared {
				var ratio float64
				if style == EdgeGradientInnerFocus {
					ratio = 1.0 - (math.Sqrt(distanceFromCenterSquared)-innerRadius)/edgeThickness
				} else if style == EdgeGradientOuterFocus {
					ratio = (math.Sqrt(distanceFromCenterSquared) - innerRadius) / edgeThickness
				} else if style == EdgeNoGradient {
					ratio = 1.0
				}
				color.A = uint8(ratio * 255.0)
			} else {
				continue
			}

			quarterImage.Set(i, j, color)
			quarterImage.Set(j, i, color)
		}
	}

	return createQuarterReflection(quarterImage, imageEdgeLength%2 == 1)
}

// Circle will fill this imageWidth
func CreateCircleImageNoAlpha(imageEdgeLength int, edgeThickness float64, style CircleStyle) image.Image {
	// If it's an odd number, the quarter image needs to cover that middle pixel
	edgeLengthOfQuarterImage := imageEdgeLength/2 + imageEdgeLength%2

	quarterImage := image.NewRGBA(image.Rectangle{image.Pt(0, 0), image.Pt(edgeLengthOfQuarterImage, edgeLengthOfQuarterImage)})

	// Radius of the full circle, including anti-aliased soft edge
	outerRadius := float64(imageEdgeLength) / 2.0
	// Radius of only the portion that is fully shaded in
	innerRadius := outerRadius - edgeThickness

	outerRadiusSquared := outerRadius * outerRadius
	innerRadiusSquared := innerRadius * innerRadius

	center := float64(imageEdgeLength) / 2.0

	color := color.NRGBA{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < edgeLengthOfQuarterImage; i++ {
		x := float64(i) + 0.5
		for j := i; j < edgeLengthOfQuarterImage; j++ {
			y := float64(j) + 0.5

			diffX := center - x
			diffY := center - y
			distanceFromCenterSquared := diffX*diffX + diffY*diffY
			if distanceFromCenterSquared < innerRadiusSquared {
				if style == FilledIn {
					color.R, color.G, color.B = 0xff, 0xff, 0xff
				} else {
					continue
				}
			} else if distanceFromCenterSquared < outerRadiusSquared {
				var ratio float64
				if style == EdgeGradientInnerFocus {
					ratio = 1.0 - (math.Sqrt(distanceFromCenterSquared)-innerRadius)/edgeThickness
				} else if style == EdgeGradientOuterFocus {
					ratio = (math.Sqrt(distanceFromCenterSquared) - innerRadius) / edgeThickness
				} else if style == EdgeNoGradient {
					ratio = 1.0
				}
				value := uint8(ratio * 255.0)
				color.R, color.G, color.B = value, value, value
			} else {
				continue
			}

			quarterImage.Set(i, j, color)
			quarterImage.Set(j, i, color)
		}
	}

	return createQuarterReflection(quarterImage, imageEdgeLength%2 == 1)
}
