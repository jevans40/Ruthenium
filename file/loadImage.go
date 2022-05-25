package file

import (
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"

	log "github.com/sirupsen/logrus"
)

//TODO:: Documentation

type notOpaqueRGBA struct {
	*image.RGBA
}

func (i *notOpaqueRGBA) Opaque() bool {
	return false
}

//TODO:: Handle errors gracefully; no need to crash if a image is missing.
func LoadImageFromFile(filename string) *image.RGBA {
	f, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	img, fmtName, err := image.Decode(f)
	log.Debug(fmtName)
	if err != nil {
		panic(err)
	}
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m
}

func SaveImageToFile(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Panic(err)
	}
}

func DrawToImage(src, dst image.Image, sp image.Point) *image.RGBA {
	b := dst.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), dst, b.Min, draw.Src)
	r := image.Rectangle{sp, sp.Add(src.Bounds().Size())}
	draw.Draw(m, r, src, src.Bounds().Min, draw.Src)
	return m
}
