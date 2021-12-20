package linmath

//Matrix4f Base Matrix object, can just print what it holds
type Matrix4f interface {
	ToFloats() [16]float32
}

type orthoMat4f [16]float32

//NewOrthoMat4f Makes an orthoganal projection matrix
func NewOrthoMat4f(bottom, top, left, right, near, far float32) Matrix4f {
	if right-left == 0 || top-bottom == 0 || far-near == 0 {
		return &orthoMat4f{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	return &orthoMat4f{2 / (right - left), 0, 0, 0,
		0, 2 / (top - bottom), 0, 0,
		0, 0, -2 / (far - near), 0,
		-(right + left) / (right - left), -(top + bottom) / (top - bottom), -(far + near) / (far - near), 1}
}

func (thisMat *orthoMat4f) ToFloats() [16]float32 {
	return *thisMat
}
