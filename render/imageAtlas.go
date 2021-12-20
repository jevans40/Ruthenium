package render

import (
	"image"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/jevans40/Ruthenium/file"
	"github.com/jevans40/Ruthenium/linmath"
	log "github.com/sirupsen/logrus"
)

//TODO:: Add padding so mipmaps work

type ImageAtlas struct {
	imageSize int32
	images    map[string]*imageIndex
	atlases   []image.Image
}

type atlasRec struct {
	linmath.PSRectangle
	atlasnum int
}

type atlasPoint struct {
	linmath.PSPoint
	atlasnum int
}

type imageIndex struct {
	boundingRect atlasRec
	//TODO:: maybe in the future I can save image space by rotating images if they are wider then they are long
	//rotated      bool
	indexedImage image.Image
}

//The atlas should be given all its required images before it is used
//Once all images have been loaded into the atlas init() should be called
//This will setup the atlas, pack the sprites and no more images can be added
//TODO:: In the future I could add a repack function that just recalculates the atlas
func ImageAtlasFactory(ImageSize int32, Images int32) ImageAtlas {
	atlas := new(ImageAtlas)
	atlas.images = make(map[string]*imageIndex)
	atlas.atlases = make([]image.Image, Images)
	for i := range atlas.atlases {
		atlas.atlases[i] = image.NewRGBA(image.Rect(0, 0, int(ImageSize), int(ImageSize)))
	}
	atlas.imageSize = ImageSize

	return *atlas
}

func (i *ImageAtlas) AddImagesFromFolder(folderPath string) {
	//TODO:: Add a check to make sure that the x and y are both less then the imagesize
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		matched, _ := regexp.MatchString(`.png$|.PNG$`, path)
		if matched {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, s := range files {
		newimage := file.LoadImageFromFile(s)
		newRect := linmath.NewPSRectangle(int32(newimage.Rect.Max.X)-int32(newimage.Rect.Min.X), int32(newimage.Rect.Max.Y)-int32(newimage.Rect.Min.Y), 0, 0)
		i.images[s] = &imageIndex{atlasRec{newRect, -1}, newimage}
	}
}

func (i *ImageAtlas) AddImageFromImage(newImage image.Image, name string) {
	//TODO:: Add a check to make sure that the x and y are both less then the imagesize
	b := newImage.Bounds()
	newRect := linmath.NewPSRectangle(int32(b.Max.X-b.Min.X), int32(b.Max.Y-b.Min.Y), 0, 0)
	i.images[name] = &imageIndex{atlasRec{newRect, -1}, newImage}

}

func (i *ImageAtlas) Init() {
	imageSizes := map[string]int{}
	for s, l_image := range i.images {
		imageSizes[s] = int(l_image.boundingRect.GetSize().Y())
	}

	//This is to sort the map and give back a sorted list
	type kv struct {
		Key   string
		Value int
	}

	var sortedSizes []kv
	for k, v := range imageSizes {
		sortedSizes = append(sortedSizes, kv{k, v})
	}

	sort.Slice(sortedSizes, func(i, j int) bool {
		return sortedSizes[i].Value > sortedSizes[j].Value
	})

	var empty []atlasRec
	var rows []atlasPoint
	var rowheights []int
	smallestY := sortedSizes[len(sortedSizes)-1].Value
	for r := 0; r < len(i.atlases); r++ {
		rows = append(rows, atlasPoint{linmath.NewPSPoint(0, int32(sortedSizes[0].Value)), r})
		rowheights = append(rowheights, 0)
	}
	//TODO:: later I can sort empty by y values and only check the entries where its possible to fit a new object.
	for _, v := range sortedSizes {
		allocated := false
		currentImage := i.images[v.Key]
		for atnum := range i.atlases {
			for boxindex, box := range empty {
				if currentImage.boundingRect.GetSize().X() <= box.GetSize().X() && currentImage.boundingRect.GetSize().Y() <= box.GetSize().Y() && atnum == box.atlasnum {
					allocated = true
					es := box.GetSize()
					ep := box.GetPoint()
					bs := currentImage.boundingRect.GetSize()

					if int(es.Y()-bs.Y()) >= smallestY {
						empty = append(empty, atlasRec{linmath.NewPSRectangle(bs.X(), es.Y()-bs.Y(), ep.X(), ep.Y()+bs.Y()), atnum})
					}
					if int(es.X()-bs.X()) >= smallestY {
						empty = append(empty, atlasRec{linmath.NewPSRectangle(es.X()-bs.X(), es.Y(), ep.X()+bs.X(), ep.Y()), atnum})
					}
					empty[boxindex] = empty[len(empty)-1]
					empty = empty[:len(empty)-1]
					i.images[v.Key].boundingRect.SetPoint(ep)
					i.images[v.Key].boundingRect.atlasnum = atnum
					break
				}
			}
			if allocated {
				break
			}
			lastHeight := 0
			for rownum, row := range rows {
				if row.atlasnum == atnum {
					lastHeight = int(row.Y())
					if i.imageSize > row.X()+currentImage.boundingRect.GetSize().X() {
						allocated = true
						i.images[v.Key].boundingRect.SetPoint(linmath.NewPSPoint(row.X(), int32(rowheights[rownum])))
						i.images[v.Key].boundingRect.atlasnum = atnum
						row.Setx(row.X() + currentImage.boundingRect.GetSize().X())
						empty = append(empty, atlasRec{linmath.NewPSRectangle(
							currentImage.boundingRect.GetSize().X(),
							row.Y()-(currentImage.boundingRect.GetPoint().Y()+currentImage.boundingRect.GetSize().Y()),
							currentImage.boundingRect.GetPoint().X(),
							(currentImage.boundingRect.GetPoint().Y() + currentImage.boundingRect.GetSize().Y()),
						), atnum})
						break
					}
				}
			}
			if allocated {
				break
			}
			if int32(lastHeight)+currentImage.boundingRect.GetSize().Y() < i.imageSize {
				allocated = true
				rows = append(rows, atlasPoint{linmath.NewPSPoint(currentImage.boundingRect.GetSize().X(), int32(lastHeight)+currentImage.boundingRect.GetSize().Y()), atnum})
				rowheights = append(rowheights, lastHeight)
				i.images[v.Key].boundingRect.SetPoint(linmath.NewPSPoint(0, int32(lastHeight)))
				i.images[v.Key].boundingRect.atlasnum = atnum
			}

		}
		if allocated {
			continue
		} else {
			log.Error("OUT OF TEXTURE UNIT MEMORY ENGINE CLOSING")
			log.Panic("Out of TEXTURE UNIT MEMORY ENGINE CLOSING")
		}
	}

	//Finally draw the allocated images to their respective atlases
	for _, img := range i.images {
		b := img.boundingRect
		i.atlases[b.atlasnum] = file.DrawToImage(img.indexedImage, i.atlases[b.atlasnum], image.Point{int(b.GetPoint().X()), int(b.GetPoint().Y())})
	}

}

func (i *ImageAtlas) getAtlas(index int) image.Image {
	return i.atlases[index]
}
