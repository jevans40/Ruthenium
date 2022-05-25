package linmath

//TODO:: Documentation

//TODO:: Update with generics
//Vector4 Default interface to export struct information
type Vector4 interface {
	GetX() float32
	GetY() float32
	GetZ() float32
	GetW() float32

	SetX(float32)
	SetY(float32)
	SetZ(float32)
	SetW(float32)
}

//vector4 implements Vector4 and Vector3
type vector4 struct {
	x float32
	y float32
	z float32
	w float32
}

//NewVector4 Vector factory, creates a new vector with given parameters
//x float32
//y float32
//z float32
//w float32
func NewVector4(x float32, y float32, z float32, w float32) Vector4 {
	newVec := vector4{x, y, z, w}
	return &newVec
}

//AddVector4 Adds two vectors together and returns a new Vector4
func AddVector4(vec1 Vector4, vec2 Vector4) Vector4 {
	newVec := NewVector4(vec1.GetX()+vec2.GetX(), vec1.GetY()+vec2.GetY(), vec1.GetZ()+vec2.GetZ(), vec1.GetW()+vec2.GetW())
	return newVec
}

//SubVector4 subtracts two vectors and returns a new Vector4
func SubVector4(vec1 Vector4, vec2 Vector4) Vector4 {
	newVec := NewVector4(vec1.GetX()-vec2.GetX(), vec1.GetY()-vec2.GetY(), vec1.GetZ()-vec2.GetZ(), vec1.GetW()-vec2.GetW())
	return newVec
}

//ScalarMult4 Multiplies given Vector by scalar multiple s
func ScalarMult4(thisVec *Vector4, s float32) {
	(*thisVec).SetX((*thisVec).GetX() * s)
	(*thisVec).SetY((*thisVec).GetY() * s)
	(*thisVec).SetZ((*thisVec).GetZ() * s)
	(*thisVec).SetW((*thisVec).GetW() * s)
}

func (thisVec vector4) GetX() float32 {
	return thisVec.x
}

func (thisVec vector4) GetY() float32 {
	return thisVec.y
}

func (thisVec vector4) GetZ() float32 {
	return thisVec.z
}

func (thisVec vector4) GetW() float32 {
	return thisVec.w
}

func (thisVec *vector4) SetX(x float32) {
	thisVec.x = x
}

func (thisVec *vector4) SetY(y float32) {
	thisVec.y = y
}

func (thisVec *vector4) SetZ(z float32) {
	thisVec.z = z
}

func (thisVec *vector4) SetW(w float32) {
	thisVec.w = w
}
