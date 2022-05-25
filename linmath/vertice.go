package linmath

//TODO:: Documentation

//This file contains the default interface for vertice

//Vertice interface, all things that are rendered should go through this
type FourVertice interface {
	SetX(float32, float32)
	SetY(float32, float32)
	SetZ(float32)

	SetTexX(float32, float32)
	SetTexY(float32, float32)
	SetColor([4]uint8)
	SetMap(uint32)

	ToFloats() *[28]float32
}

type fourvertice struct {
	val [28]float32
}

//NewVertice Factory for creating a new vertice

func EmptyFourVertice() FourVertice {
	return &fourvertice{[28]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
}

func (thisVert *fourvertice) SetX(x, width float32) {
	thisVert.val[0] = x
	thisVert.val[7] = x + width
	thisVert.val[14] = x
	thisVert.val[21] = x + width
}

func (thisVert *fourvertice) SetY(y, height float32) {
	thisVert.val[1] = y
	thisVert.val[8] = y
	thisVert.val[15] = y + height
	thisVert.val[22] = y + height
}

func (thisVert *fourvertice) SetZ(z float32) {
	thisVert.val[2] = z
	thisVert.val[9] = z
	thisVert.val[16] = z
	thisVert.val[23] = z
}

func (thisVert *fourvertice) SetTexX(x, width float32) {
	thisVert.val[3] = x
	thisVert.val[10] = x + width
	thisVert.val[17] = x
	thisVert.val[24] = x + width
}

func (thisVert *fourvertice) SetTexY(y, height float32) {
	thisVert.val[4] = y
	thisVert.val[11] = y
	thisVert.val[18] = y + height
	thisVert.val[25] = y + height
}
func (thisVert *fourvertice) SetColor(color [4]uint8) {
	ccolor := intToFloat(charToInt(color))
	thisVert.val[5] = ccolor
	thisVert.val[12] = ccolor
	thisVert.val[19] = ccolor
	thisVert.val[26] = ccolor
}
func (thisVert *fourvertice) SetMap(cmap uint32) {
	thisVert.val[6] = intToFloat(cmap)
	thisVert.val[13] = intToFloat(cmap)
	thisVert.val[20] = intToFloat(cmap)
	thisVert.val[27] = intToFloat(cmap)

}

func (thisVert *fourvertice) ToFloats() *[28]float32 {
	return &thisVert.val
}
