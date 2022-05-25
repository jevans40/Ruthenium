package linmath

//TODO:: Documentation

type Mat3[T Numeric] [9]T

//Returns a new Identiy matrix
func (m Mat3[T]) Identity() Mat3[T] {
	return [9]T{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

//Returns the mutliple of this and MatB
func (m Mat3[T]) MatMul(MatB Mat3[T]) Mat3[T] {
	return [9]T{
		m[0]*MatB[0] + m[1]*MatB[3] + m[2]*MatB[6],
		m[0]*MatB[1] + m[1]*MatB[4] + m[2]*MatB[7],
		m[0]*MatB[2] + m[1]*MatB[5] + m[2]*MatB[8],
		m[3]*MatB[0] + m[4]*MatB[3] + m[5]*MatB[6],
		m[3]*MatB[1] + m[4]*MatB[4] + m[5]*MatB[7],
		m[3]*MatB[2] + m[4]*MatB[5] + m[5]*MatB[8],
		m[6]*MatB[0] + m[7]*MatB[3] + m[8]*MatB[6],
		m[6]*MatB[1] + m[7]*MatB[4] + m[8]*MatB[7],
		m[6]*MatB[2] + m[7]*MatB[5] + m[8]*MatB[8],
	}
}

func (m Mat3[T]) Addition(MatB Mat3[T]) Mat3[T] {
	return [9]T{
		m[0] + MatB[0],
		m[1] + MatB[1],
		m[2] + MatB[2],
		m[3] + MatB[3],
		m[4] + MatB[4],
		m[5] + MatB[5],
		m[6] + MatB[6],
		m[7] + MatB[7],
		m[8] + MatB[8],
	}
}

func (m Mat3[T]) Subtraction(MatB Mat3[T]) Mat3[T] {
	return [9]T{
		m[0] - MatB[0],
		m[1] - MatB[1],
		m[2] - MatB[2],
		m[3] - MatB[3],
		m[4] - MatB[4],
		m[5] - MatB[5],
		m[6] - MatB[6],
		m[7] - MatB[7],
		m[8] - MatB[8],
	}
}

func (m Mat3[T]) ScalarMul(Scalar T) Mat3[T] {
	return [9]T{
		m[0] * Scalar,
		m[1] * Scalar,
		m[2] * Scalar,
		m[3] * Scalar,
		m[4] * Scalar,
		m[5] * Scalar,
		m[6] * Scalar,
		m[7] * Scalar,
		m[8] * Scalar,
	}
}
