package linmath

import "testing"

//Tests for mat4f.go
//TODO:: More tests
func TestNewOrthoMat4f(t *testing.T) {
	TestOrthoCon := NewOrthoMat4f(1, 0, 0, 1, 1, 0)
	TestOrthoNeg := NewOrthoMat4f(-1, -0, -0, -1, -1, -0)
	TestOrthoZed := NewOrthoMat4f(0, 0, 0, 0, 0, 0)
	if TestOrthoCon.ToFloats() != [16]float32{2, 0, 0, 0, 0, -2, 0, 0, 0, 0, 2, 0, -1, 1, 1, 1} {
		t.Error()
	}
	if TestOrthoNeg.ToFloats() != [16]float32{-2, 0, 0, 0, 0, 2, 0, 0, 0, 0, -2, 0, -1, 1, 1, 1} {
		t.Error()
	}
	if TestOrthoZed.ToFloats() != [16]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Error()
	}

	triangle := make([]Vector3, 3)

	triangle[0] = NewVector3(-1, -1, 0)
	triangle[1] = NewVector3(1, -1, 0)
	triangle[2] = NewVector3(0, 1, 0)

	searialTriangle := []float32{}
	for _, v := range triangle {
		x := v.ToFloats()
		searialTriangle = append(searialTriangle, x[:]...)
	}
}
