package linmath

//Vector3 Default interface to export struct information
type Vector3 interface {
	GetX() float32
	GetY() float32
	GetZ() float32

	SetX(float32)
	SetY(float32)
	SetZ(float32)

	ToFloats() [3]float32
}

//vector3 implements Vector3
type vector3 [3]float32

//NewVector3 Vector factory, creates a new vector with given parameters
//x float32
//y float32
//z float32
func NewVector3(x float32, y float32, z float32) Vector3 {
	newVec := vector3{x, y, z}
	return &newVec
}

//AddVector3 Adds two vectors together and returns a new Vector3
func AddVector3(vec1 Vector3, vec2 Vector3) Vector3 {
	newVec := NewVector3(vec1.GetX()+vec2.GetX(), vec1.GetY()+vec2.GetY(), vec1.GetZ()+vec2.GetZ())
	return newVec
}

//SubVector3 subtracts two vectors and returns a new Vector3
func SubVector3(vec1 Vector3, vec2 Vector3) Vector3 {
	newVec := NewVector3(vec1.GetX()-vec2.GetX(), vec1.GetY()-vec2.GetY(), vec1.GetZ()-vec2.GetZ())
	return newVec
}

//ScalarMult3 Multiplies given Vector by scalar multiple s
func ScalarMult3(thisVec *Vector3, s float32) {
	(*thisVec).SetX((*thisVec).GetX() * s)
	(*thisVec).SetY((*thisVec).GetY() * s)
	(*thisVec).SetZ((*thisVec).GetZ() * s)
}

func (thisVec vector3) GetX() float32 {
	return thisVec[0]
}

func (thisVec vector3) GetY() float32 {
	return thisVec[1]
}

func (thisVec vector3) GetZ() float32 {
	return thisVec[2]
}

func (thisVec *vector3) SetX(x float32) {
	thisVec[0] = x
}

func (thisVec *vector3) SetY(y float32) {
	thisVec[0] = y
}

func (thisVec *vector3) SetZ(z float32) {
	thisVec[0] = z
}

func (thisVec *vector3) ToFloats() [3]float32 {
	return *thisVec
}
