package component

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testComponent struct {
	value int
}

func (t testComponent) IsComponent()          {}
func (t testComponent) GetType() reflect.Type { return reflect.TypeOf(t) }

func TestDenseStorage(t *testing.T) {
	comp := testComponent{}
	testStorage := NewDenseStorage[testComponent]()
	readStorage, err := GetReadOnlyStorage[testComponent](testStorage)
	assert.NoError(t, err, "Test1.A ReadOnlyStorage failed to initalize")
	writeStorage, err := GetWriteStorage[testComponent](testStorage)
	assert.NoError(t, err, "Test1.B WriteStorage failed to initalize")

	//Test1
	assert.Equal(t, testStorage.GetType(), reflect.TypeOf(comp), "TEST1.C: Type of stored object does not match component")

	test2, err := readStorage.GetComponent(1)
	var empty testComponent
	if err == nil || test2 != empty {
		t.Error("TEST2: GetComponent did not return an error for an invalid component")
	}

	//Test
	test3Err1 := writeStorage.AddEntityMultiple([]EntityID{65}, []testComponent{{1}, {2}})
	writeStorage.AddEntityMultiple([]EntityID{65, 32}, []testComponent{{1}, {2}})
	test3Err2 := writeStorage.AddEntityMultiple([]EntityID{65}, []testComponent{{1}})
	test3, err := readStorage.GetComponent(65)
	totest := testComponent{1}
	if test3 != totest {
		t.Error("TEST3.A: Test3 was nil")
	}
	if err != nil {
		t.Error("TEST3.B: err was not nil")
	}
	assert.Error(t,
		test3Err1,
		"TEST3.C test3Err1 did not throw an error")

	assert.Error(t,
		test3Err2,
		"Test3.D test3Err2 did not throw an error")

	assert.Panics(t,
		func() { readStorage.MustGetComponent(100) },
		"Test3.E mustGetComponent did not panic when trying to get an invalid component")

	assert.NotPanics(t,
		func() { readStorage.MustGetComponent(65) },
		"Test3.F mustGetComponent encountered a panic for a valid element")

	assert.Panics(t,
		func() { writeStorage.AddEntityMultiple([]EntityID{13, 13}, []testComponent{{1}, {2}}) },
		"Test3.G Panic was not thrown when trying to add a duplicate element")

	testStorage = NewDenseStorage[testComponent]()
	entities := []EntityID{1, 2, 3, 1020, 27, 23, 5, 28, 12, 35, 55, 31, 21}
	comps := []testComponent{{12}, {1}, {14}, {1}, {17}, {1}, {13}, {10}, {14}, {1}, {17}, {1}, {13}}

	err = writeStorage.AddEntityMultiple(entities, comps)
	assert.NoError(t, err, "Test4.A: adding entities failed for bulk test")
	_, err = readStorage.GetComponentMultiple([]EntityID{1020, 27, 23})
	assert.NoError(t, err, "Test4.B: Getting entities failed for bulk test")
	_, err = writeStorage.GetComponentMultiple([]EntityID{1020, 27, 2233})
	assert.Error(t, err, "Test4.C: Getting entities did not fail for invalid bulk test")
	test4a, _ := readStorage.GetComponentMultiple(entities[:3])

	assert.Equal(t, comps[:3], test4a)

	assert.NotPanics(t,
		func() { readStorage.MustGetComponentMultiple(entities[:3]) },
		"Test4.D mustGetComponent encountered a panic for a valid element")
	assert.Panics(t,
		func() { readStorage.MustGetComponentMultiple(append(entities[:3], 2233)) },
		"Test4.E mustGetComponent did not encounter a panic for a invalid element")

	writeStorage.DeleteEntityMultiple(entities[:3])

	_, err = readStorage.GetComponentMultiple(entities[:3])
	assert.Error(t, err, "Test5.A found deleted components")

	err = writeStorage.DeleteEntity(EntityID(2233))
	assert.Error(t, err, "Test5.B no error for invalid entity ID")
}
