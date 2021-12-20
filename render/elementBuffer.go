package render

type ElementBuffer struct {
	elementarray []uint32
}

func (e *ElementBuffer) GetArray() []uint32 {
	return e.elementarray
}

func (e *ElementBuffer) setSize(size int) {
	e.elementarray = make([]uint32, size*6)
	for i := 0; i < size; i++ {
		e.elementarray[i*6] = uint32(i) * 4
		e.elementarray[i*6+1] = uint32(i)*4 + 1
		e.elementarray[i*6+2] = uint32(i)*4 + 2
		e.elementarray[i*6+3] = uint32(i)*4 + 1
		e.elementarray[i*6+4] = uint32(i)*4 + 2
		e.elementarray[i*6+5] = uint32(i)*4 + 3
	}
}
