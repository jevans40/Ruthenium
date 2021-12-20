package linmath

import "errors"

//A simple distance matrix with custom id rows and columns.
//The matrix is reflexive so GetDistance(i,j) == GetDistance(j,i)
type DistanceMatrixInt interface {
	//Returns the distance between two points
	//Returns an error if out of bounds
	GetDistance(i, j int) (int, error)
	//Changes the distancebetween two points
	//Returns an error if out of bounds
	ChangeDistance(i, j, dist int) error
	//Appends a new array of distances to the matrix
	//Distance to self should not be included
	Append(distances []int)
	//Removes an array of distances from the matrix
	//Returns an error if out of bounds
	Remove(item int) error

	//Returns the maximum distance, and the id of the point its farthest from
	Max(item int) (int, int)
	//Returns the minimum distance, and the id of the point its closest to
	Min(item int) (int, int)
}

type distanceMatrixInt struct {
	matrix [][]int
	max    []int
	min    []int
}

func DistanceMatrixIntFactory() DistanceMatrixInt {
	matrix := [][]int{}
	max := []int{}
	min := []int{}
	return &distanceMatrixInt{matrix: matrix, max: max, min: min}
}
func DistanceMatrixIntFromMatrixFactory(mat [][]int) DistanceMatrixInt {
	max := make([]int, len(mat)*2)
	min := make([]int, len(mat)*2)
	for i, v := range mat {
		ma := v[i+(len(mat)-1)%len(mat)]
		maid := 0
		mi := v[i+(len(mat)-1)%len(mat)]
		miid := 0
		for j, k := range v {
			if i == j {
				continue
			}
			if k > ma {
				ma = k
				maid = j
			} else if k < mi {
				ma = mi
				miid = k
			}
		}
		max[i*2] = ma
		max[(i*2)+1] = maid
		min[i*2] = mi
		min[(i*2)+1] = miid
	}
	return &distanceMatrixInt{matrix: mat, max: max, min: min}
}

func (d *distanceMatrixInt) GetDistance(i, j int) (int, error) {
	if i >= len(d.matrix) || j >= len(d.matrix) {
		return 0, errors.New("Index Out of Bounds")
	}
	return d.matrix[i][j], nil
}

func (d *distanceMatrixInt) ChangeDistance(i, j, dist int) error {
	if i >= len(d.matrix) || j >= len(d.matrix) {
		return errors.New("Index Out of Bounds")
	}
	if min, _ := d.Min(i); dist < min {
		d.min[i*2] = dist
		d.min[(i*2)+1] = j
	}
	if min, _ := d.Min(j); dist < min {
		d.min[j*2] = dist
		d.min[(j*2)+1] = i
	}
	if max, _ := d.Max(i); dist < max {
		d.max[i*2] = dist
		d.max[(i*2)+1] = j
	}
	if max, _ := d.Max(j); dist < max {
		d.max[j*2] = dist
		d.max[(j*2)+1] = i
	}
	d.matrix[i][j] = dist
	d.matrix[j][i] = dist
	return nil
}

func (d *distanceMatrixInt) Min(index int) (int, int) {
	return d.min[index*2], d.min[index*2+1]
}
func (d *distanceMatrixInt) Max(index int) (int, int) {
	return d.max[index*2], d.max[index*2+1]
}

//Distance to self not included
func (d *distanceMatrixInt) Append(array []int) {
	for i, _ := range d.matrix {
		d.matrix[i] = append(d.matrix[i], array[i])
	}
	array = append(array, 0)
	d.matrix = append(d.matrix, array)
}

func (d *distanceMatrixInt) Remove(index int) error {
	if index >= len(d.matrix) {
		return errors.New("Index Out of Bounds")
	}
	d.matrix = append(d.matrix[:index-1], d.matrix[index+1:]...)
	for i, _ := range d.matrix {
		d.matrix[i] = append(d.matrix[i][0:index-1], d.matrix[i][index+1:]...)
	}
	return nil
}
