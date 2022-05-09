package component

import (
	"errors"
	"math"
	"reflect"

	"github.com/jinzhu/copier"
)

//Custom standard errors for future error handling
type entityNotFound string
type oneOrMoreEntitiesAlreadyExists string

func (e entityNotFound) Error() string                 { return string(e) }
func (e oneOrMoreEntitiesAlreadyExists) Error() string { return string(e) }

const EntityNotFoundError = entityNotFound("EntityNotFound")
const OneOrMoreEntitiesAlreadyExists = oneOrMoreEntitiesAlreadyExists("OneOrMoreEntitiesAlreadyExists")

//A component storage should only store one type and should always
//Respect basic database ideas. Entites should be immutable without a
//Write reference, read references should be prefered, and all mutable
//interactions should always do either all changes or none.
type ComponentStorage interface {

	//Returns the reflect type of the stored components
	GetType() reflect.Type

	//Return true if the entityID is associated with a component in the storage
	Exists(EntityID) bool

	//Returns a mask of all components found in the storage
	ExistsMultiple([]EntityID) []bool

	//Returns a list of all entities found in the storage.
	GetEntities() []EntityID

	//Returns the number of components stored here.
	GetSize() int

	//Adds a new empty component to the storage immediately.
	//Optimization: This should ideally not be called if possible from services
	//WriteStoreage[T] allows to add a preInitalized element and wont slow it down
	//And dispatcher should allow for lazy creation
	AddBlankComponent(entity EntityID) error

	//calls AddComponent on all the entityIDs in the list
	AddBlankComponentMultiple(entities []EntityID) error

	//Allows for deletion of element in the storage immediately
	//Optimization: This should ideally not be called if possible from services
	//WriteStoreage[T]: safe to call this without slowdown
	//And dispatcher should allow for lazy deletion using this method
	DeleteEntity(entitity EntityID) error

	//calls DeleteEntity on all the entityIDs in the list immediately.
	DeleteEntityMultiple(entitities []EntityID) error
}

type ReadOnlyStorage[T Component] interface {
	//Return true if the entityID is associated with a component in the storage
	Exists(EntityID) bool

	//Returns a mask of all components found in the storage
	ExistsMultiple([]EntityID) []bool

	//Returns a list of all entities found in the storage.
	GetEntities() []EntityID

	//Returns the number of components stored here.
	GetSize() int

	//Returns struct copies of entities from the storage
	//Note: This will panic if it cannot find the given entityID
	//Call GetComponent if you want to just receive an error
	MustGetComponent(entitityID EntityID) T

	//Calls MustGetComponent on all entities listed
	MustGetComponentMultiple(entities []EntityID) []T

	//Returns struct copies of the requested entities from storage
	GetComponent(entitityID EntityID) (T, error)

	//Calls GetComponent on all entities listed
	GetComponentMultiple(entities []EntityID) ([]*T, error)
}

type WriteStorage[T Component] interface {
	//Return true if the entityID is associated with a component in the storage
	Exists(EntityID) bool

	//Returns a mask of all components found in the storage
	ExistsMultiple([]EntityID) []bool

	//Returns a list of all entities found in the storage.
	GetEntities() []EntityID

	//Returns the number of components stored here.
	GetSize() int

	//Writes Data to the specified entityID
	Write(entityID EntityID, data T) error

	//Writes Data[i] to each EntityID[i]
	WriteMultiple(entityIDs []EntityID, data []*T) error

	//Returns struct copies of entities from the storage
	//Note: This will panic if it cannot find the given entityID
	//Call GetComponent if you want to just receive an error
	MustGetComponent(entitityID EntityID) T

	//Calls MustGetComponent on all entities listed
	MustGetComponentMultiple(entities []EntityID) []T

	//Returns struct copies of the requested entities from storage
	GetComponent(entitityID EntityID) (T, error)

	//Calls GetComponent on all entities listed
	GetComponentMultiple(entities []EntityID) ([]*T, error)

	//Adds a new empty component to the storage immediately.
	//Optimization: This should ideally not be called if possible from services
	//WriteStoreage[T] allows to add a preInitalized element and wont slow it down
	//And dispatcher should allow for lazy creation
	AddBlankComponent(entity EntityID) error

	//calls AddComponent on all the entityIDs in the list
	AddBlankComponentMultiple(entities []EntityID) error

	//Allows for deletion of element in the storage immediately
	//Optimization: This should ideally not be called if possible from services
	//WriteStoreage[T]: safe to call this without slowdown
	//And dispatcher should allow for lazy deletion using this method
	DeleteEntity(entitity EntityID) error

	//calls DeleteEntity on all the entityIDs in the list immediately.
	DeleteEntityMultiple(entitities []EntityID) error

	//Appends an component to the end of the list
	AddEntity(entity EntityID, component T) error

	//Appends components to the end of the list
	//This will panic if entityID is listed multiple times
	AddEntityMultiple(Entitylist []EntityID, Components []T) error
}

type Joinable interface {
	//Returns a mask of all components found in the storage
	ExistsMultiple([]EntityID) []bool

	//Returns a list of all entities found in the storage.
	GetEntities() []EntityID

	//Returns the number of components stored here.
	GetSize() int
}

//Joins multiple sets of ComponentStorages to get a set of common entities.
func Join(toJoin ...Joinable) (joinedSet []EntityID) {
	min := 0
	val := math.MaxInt
	for i, join := range toJoin {
		if join.GetSize() < val {
			min = i
			val = join.GetSize()
		}
	}
	toCheck := toJoin[min].GetEntities()
	for _, v := range toJoin {
		toPrune := v.ExistsMultiple(toCheck)
		for i, v := range toPrune {
			if !v {
				toCheck[i] = -1
			}
		}
	}
	for _, v := range toCheck {
		if v != -1 {
			joinedSet = append(joinedSet, v)
		}
	}
	return joinedSet
}

//Asserts a component storage into a generic typed read-only representation
//This is not garenteed to work, and will panic if storage type does not implement ReadOnlyStorage
func GetReadOnlyStorage[T Component](storage ComponentStorage) (ReadOnlyStorage[T], error) {
	var toTest T
	var newVStorage = VectorStorage[T]{}
	var newDStorage = DenseStorage[T]{}
	if storage.GetType() == reflect.TypeOf(toTest) {
		switch reflect.TypeOf(storage) {
		case reflect.TypeOf(&newVStorage):
			copier.Copy(&newVStorage, storage)
			return ComponentStorage(&newVStorage).(ReadOnlyStorage[T]), nil
		case reflect.TypeOf(&newDStorage):
			copier.Copy(&newDStorage, storage)
			return ComponentStorage(storage).(ReadOnlyStorage[T]), nil
		default:

			//fmt.Printf("%T\n", Type.GetTypeCode((*VectorStorage[T])))
		}
	}
	return nil, errors.New("Type mismatch for given storage and function generic")
}

//Asserts a component storage into a generic typed writeable representation
//This is not garenteed to work, and will panic if storage type does not implement WriteStorage
func GetWriteStorage[T Component](storage ComponentStorage) (WriteStorage[T], error) {
	var toTest T
	var newVStorage = VectorStorage[T]{}
	var NewDenseStorage = DenseStorage[T]{}
	if storage.GetType() == reflect.TypeOf(toTest) {
		switch reflect.TypeOf(storage) {
		case reflect.TypeOf(&newVStorage):
			copier.Copy(&newVStorage, storage)
			return ComponentStorage(storage).(WriteStorage[T]), nil
		case reflect.TypeOf(&NewDenseStorage):
			copier.Copy(&NewDenseStorage, storage)
			return ComponentStorage(storage).(WriteStorage[T]), nil
		}
		return (storage).(WriteStorage[T]), nil
	}
	return nil, errors.New("Type mismatch for given storage and function generic")
}
