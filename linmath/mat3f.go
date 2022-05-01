package linmath

type Mat3f[T float32 | float64] Mat3[T]

func Mat3fFactory[T float32 | float64]() Mat3f[T] {
	return [9]T{0, 0, 0, 0, 0, 0, 0, 0, 0}
}

//Returns a new Identiy matrix
func Identity[T float32 | float64]() Mat3f[T] {
	return [9]T{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

//Returns the mutliple of this and MatB
func (m Mat3f[T]) MatMul(MatB Mat3f[T]) Mat3f[T] {
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

func (m Mat3f[T]) Addition(MatB Mat3f[T]) Mat3f[T] {
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

func (m Mat3f[T]) Subtraction(MatB Mat3f[T]) Mat3f[T] {
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

func (m Mat3f[T]) ScalarMul(Scalar T) Mat3f[T] {
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

func (m Mat3f[T]) Det() float64 {
	return float64(m[0]*(m[4]*m[8]-m[5]*m[7]) - m[1]*(m[3]*m[8]-m[5]*m[6]) + m[2]*(m[3]*m[7]-m[4]*m[6]))
}

func (m Mat3f[T]) Inverse() Mat3f[T] {
	return Mat3f[T]{
		m[4]*m[8] - m[5]*m[7],
		m[2]*m[7] - m[1]*m[8],
		m[1]*m[5] - m[2]*m[4],
		m[5]*m[6] - m[3]*m[8],
		m[0]*m[8] - m[2]*m[6],
		m[2]*m[3] - m[0]*m[5],
		m[3]*m[7] - m[4]*m[6],
		m[1]*m[6] - m[0]*m[7],
		m[0]*m[4] - m[1]*m[3],
	}.ScalarMul(T(m.Det()))
}

func (m Mat3f[T]) VectorMul(X, Y, Z T) [3]T {
	return [3]T{
		m[0]*X + m[1]*Y + m[2]*Z,
		m[3]*X + m[4]*Y + m[5]*Z,
		m[6]*X + m[7]*Y + m[8]*Z,
	}
}
