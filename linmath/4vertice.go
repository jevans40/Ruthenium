package linmath

import (
	"encoding/binary"
	"unsafe"
)

//TODO:: Documentation
//TODO:: Tests

//This file contains the default interface for vertice

//Vertice interface, all things that are rendered should go through this
type Vertice interface {
	GetX() float32
	GetY() float32
	GetZ() float32

	GetTexX() float32
	GetTexY() float32
	GetColor() [4]uint8
	GetMap() uint32

	SetX(float32)
	SetY(float32)
	SetZ(float32)

	SetTexX(float32)
	SetTexY(float32)
	SetColor([4]uint8)
	SetMap(uint32)

	ToFloats() [7]float32
}

type vertice [7]float32

//NewVertice Factory for creating a new vertice
func VerticeFactory(x, y, z, tx, ty float32, R, G, B, A uint8, texmap uint32) Vertice {
	return &vertice{x, y, z, tx, ty, intToFloat(charToInt(([4]uint8{R, G, B, A}))), intToFloat(texmap)}
}

func EmptyVertice() Vertice {
	return &vertice{0, 0, 0, 0, 0, 0, 0}
}

//charToInt converts a char array to a single int
func charToInt(input [4]uint8) uint32 {
	return binary.BigEndian.Uint32(([]byte)(input[:]))
}

func intToChar(input uint32) [4]uint8 {
	var bs [4]byte
	binary.LittleEndian.PutUint32(bs[:], input)
	return bs
}

//intToFloat THIS FUNCTION IS STUPID
//WARINING this is an unsafe conversion
//that is only used to be converted back at a later date
//For the love of all that is holy dont use this for anything else!!!
//If you dont know why this function is stupid, then please dont use it
func intToFloat(input uint32) float32 {
	return *(*float32)(unsafe.Pointer(&input))
}

func floatToInt(input float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&input))
}

func (thisVert *vertice) GetX() float32 {
	return thisVert[0]
}

func (thisVert *vertice) GetY() float32 {
	return thisVert[1]
}

func (thisVert *vertice) GetZ() float32 {
	return thisVert[2]
}

func (thisVert *vertice) GetTexX() float32 {
	return thisVert[3]
}

func (thisVert *vertice) GetTexY() float32 {
	return thisVert[4]
}

func (thisVert *vertice) GetColor() [4]uint8 {
	return intToChar(floatToInt(thisVert[5]))
}

func (thisVert *vertice) GetMap() uint32 {
	return floatToInt(thisVert[6])
}

func (thisVert *vertice) SetX(x float32) {
	thisVert[0] = x
}

func (thisVert *vertice) SetY(y float32) {
	thisVert[1] = y
}

func (thisVert *vertice) SetZ(z float32) {
	thisVert[2] = z
}

func (thisVert *vertice) SetTexX(x float32) {
	thisVert[3] = x
}

func (thisVert *vertice) SetTexY(y float32) {
	thisVert[4] = y
}
func (thisVert *vertice) SetColor(color [4]uint8) {
	thisVert[5] = intToFloat(charToInt(color))
}
func (thisVert *vertice) SetMap(cmap uint32) {
	thisVert[6] = intToFloat(cmap)
}

func (thisVert *vertice) ToFloats() [7]float32 {
	return *thisVert
}
