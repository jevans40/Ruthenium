package file

import (
	"image"
	"io/ioutil"

	"github.com/jevans40/Ruthenium/linmath"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//TODO:: Documentation
//TODO:: Find a way to make this work on multiple threads

type PSFont struct {
	face     font.Face
	location fixed.Point26_6
}

//TODO this needs to be individual to each thread, otherwise I think this will clash the threads

func LoadFont(filename string, fontsize int) PSFont {
	fontBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}

	utf8Font, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Panic(err)
	}

	face, err := opentype.NewFace(utf8Font, &opentype.FaceOptions{
		Size:    float64(fontsize),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Panic(err)
	}

	return PSFont{face, fixed.Point26_6{}}
}

//Returns where to draw the char
func (fnt *PSFont) DrawChar(c rune, p rune) (topCorner linmath.PSPoint) {
	if p >= 0 {
		fnt.location.X += fnt.face.Kern(p, c)
	}
	bnds, adv, err := fnt.face.GlyphBounds(c)
	if !err {
		//logging.Log.Panic("Error loading glyph of rune %s", c)
	}
	fnt.location.X += adv
	return linmath.NewPSPoint(int32(bnds.Min.X+fnt.location.X)>>6, int32(bnds.Min.Y+fnt.location.Y)>>6)

}

func (fnt *PSFont) SetPos(newPos linmath.PSPoint) {
	fnt.location.X = fixed.Int26_6(newPos.X() << 6)
	fnt.location.Y = fixed.Int26_6(newPos.X()<<6) + fnt.face.Metrics().Height
}

func (fnt *PSFont) DrawString(toDraw string) (drawPoints []linmath.PSPoint) {
	lastrune := rune(-1)
	drawPoints = make([]linmath.PSPoint, len(toDraw))
	for i, c := range toDraw {
		drawPoints[i] = fnt.DrawChar(c, lastrune)
		lastrune = c
	}
	return
}

func (fnt *PSFont) GetGlyphImage(c rune) image.Image {
	bnds, _, _ := fnt.face.GlyphBounds(c)
	//wbnds, adv, _ := fnt.face.GlyphBounds('w')
	stringify := string(c)
	log.Trace(stringify)
	//soIcanread := []int{int(bnds.Min.X) >> 6, int(bnds.Max.X) >> 6, int(bnds.Min.Y) >> 6, int(bnds.Max.Y) >> 6}
	dst := image.NewGray((image.Rect(0, 0, (int(bnds.Max.X)>>6)-(int(bnds.Min.X)>>6), (int(bnds.Max.Y)>>6)-(int(bnds.Min.Y)>>6))))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: fnt.face,
		Dot:  fixed.Point26_6{X: ((bnds.Min.X) * -1) + fnt.face.Kern('w', c), Y: bnds.Min.Y * -1},
	}
	d.DrawString(string(c))
	return dst
}
