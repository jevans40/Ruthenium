package component

//TODO:: Generate Tests

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var _ ComponentStorage = &VectorStorage[BaseComponent]{}
var _ ReadOnlyStorage[BaseComponent] = &VectorStorage[BaseComponent]{}
var _ WriteStorage[BaseComponent] = &VectorStorage[BaseComponent]{}

//The Vector Storage struct:
//This type of storage stores all components in a dense array.
//It uses a map from EntityID's to the internal storage map for lookup.
//TODO: Evaluate later but this storage type is litterally just a big
//memory leak, we should look into storage/entity cleaning sometime.
type VectorStorage[T Component] struct {
	internalVector []T
	//Check if the value is allocated
	allocated []bool
	numStored int
	RWLOCK    sync.RWMutex
}

//Create a new Dense Storage containing types T.
//Returns a ComponentStorage interface
func NewVectorStorage[T Component]() ComponentStorage {
	return &VectorStorage[T]{internalVector: []T{}, RWLOCK: sync.RWMutex{}}
}

//Returns the type of the contained storage
func (ve *VectorStorage[T]) GetType() reflect.Type {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	return reflect.TypeOf(ve.internalVector).Elem()
}

//func (ve *VectorStorage[T]) GetData() []Component {
//	var newArray := make()
//}

//Returns true if the entitity is stored in the internal map
func (ve *VectorStorage[T]) Exists(Entity EntityID) bool {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	if Entity < 0 || int(Entity) >= len(ve.internalVector) {
		return false
	}
	return ve.allocated[Entity]
}

//Returns true if the entitity is stored in the internal map
func (ve *VectorStorage[T]) exists(Entity EntityID) bool {
	if Entity < 0 || int(Entity) >= len(ve.internalVector) {
		return false
	}
	return ve.allocated[Entity]
}

//Returns a mask of the entities that exist in this storage
func (ve *VectorStorage[T]) ExistsMultiple(Entities []EntityID) []bool {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	toReturn := []bool{}
	for _, v := range Entities {
		toReturn = append(toReturn, ve.exists(v))
	}
	return toReturn
}

//Return all stored entities in this storage
func (ve *VectorStorage[T]) GetEntities() []EntityID {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	toReturn := make([]EntityID, ve.numStored)
	for i, k := range ve.allocated {
		if k {
			toReturn[i] = EntityID(i)
		}
	}
	return toReturn
}

//Returns the number of components stored in this storage
func (ve *VectorStorage[T]) GetSize() int {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	return ve.numStored
}

//Adds a new empty component to this storage
//TODO: Might want to use allocation maps
func (ve *VectorStorage[T]) AddBlankComponent(Entity EntityID) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()
	if ve.exists(Entity) {
		return OneOrMoreEntitiesAlreadyExists
	}

	if int(Entity) >= len(ve.internalVector) {
		togrow := (int(Entity) + 1) - len(ve.internalVector)
		ve.internalVector = append(ve.internalVector, make([]T, togrow, togrow)...)
		ve.allocated = append(ve.allocated, make([]bool, togrow, togrow)...)
	}

	var newComp T
	ve.allocated[Entity] = true
	ve.internalVector[Entity] = newComp
	ve.numStored++
	return nil
}

//Adds a new Entities with associated IDs to this storage
//Returns an error if one or more entitys already exist in the storage
//If an error is returned no entities are added to the storage
//This will panic if entityID is listed multiple times
func (ve *VectorStorage[T]) AddBlankComponentMultiple(Entity []EntityID) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()
	for _, e := range Entity {
		if ve.exists(e) {
			return OneOrMoreEntitiesAlreadyExists
		}
	}
	for _, e := range Entity {
		if ve.allocated[e] {
			panic("Attemped to insert a duplicate entity, this shouldnt happen")
		}
		var newComp T
		ve.allocated[e] = true
		ve.internalVector[e] = newComp
		ve.numStored++
	}
	return nil
}

//Deletes selected elements from storage
//This take O(m + n) time since it iterates over all entities and then all interal maps
//TODO: If you want single or small number deletion use the DeleteEntity function instead.
//This should only be called by the manager of the ComponentStorage.
//This function returns an error if any of the entities in the list do not exist in the vector.
func (ve *VectorStorage[T]) DeleteEntityMultiple(Entities []EntityID) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()
	for _, e := range Entities {
		if !ve.exists(e) {
			return EntityNotFoundError
		}
	}
	for _, e := range Entities {
		if ve.allocated[e] {
			var Deleted T
			ve.allocated[e] = false
			ve.internalVector[e] = Deleted
			ve.numStored--
		}
	}
	return nil
}

//The singlecase version of DeleteEntities Multiple
//TODO: Should be avoided for now until fixed
func (ve *VectorStorage[T]) DeleteEntity(entity EntityID) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()
	return ve.DeleteEntityMultiple([]EntityID{entity})
}

//Returns struct copies of entities from the storage
//Note: This will panic if it cannot find the given entityID
//Call GetComponent if you want to just receive an error
func (ve *VectorStorage[T]) MustGetComponent(entity EntityID) T {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	if ve.exists(entity) {
		return ve.internalVector[entity]
	}
	panic(EntityNotFoundError)
}

//Returns struct copies of entities from the storage
//Note: This will panic if it cannot find the given entityID
//Call GetComponent if you want to just receive an error
func (ve *VectorStorage[T]) mustGetComponent(entity EntityID) T {
	if ve.exists(entity) {
		return ve.internalVector[entity]
	}
	panic(EntityNotFoundError)
}

//Calls MustGetComponent on all entities listed
func (ve *VectorStorage[T]) MustGetComponentMultiple(entities []EntityID) []T {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	returnArray := []T{}
	for _, v := range entities {
		k := ve.mustGetComponent(v)
		returnArray = append(returnArray, k)
	}
	return returnArray
}

//Calls GetComponent on all entities listed
func (ve *VectorStorage[T]) GetComponent(entity EntityID) (T, error) {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()

	if ve.exists(entity) {
		return ve.internalVector[entity], nil
	}

	var errorFound T
	return errorFound, EntityNotFoundError
}

//Calls GetComponent on all entities listed
func (ve *VectorStorage[T]) getComponent(entity EntityID) (*T, error) {

	if ve.exists(entity) {
		return &(ve.internalVector[entity]), nil
	}

	var errorFound T
	return &errorFound, EntityNotFoundError
}

//Returns struct copies of the requested entities from storage
func (ve *VectorStorage[T]) GetComponentMultiple(entities []EntityID) ([]*T, error) {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()
	returnArray := make([]*T, len(entities))
	err2 := error(nil)
	for i, _ := range entities {
		k, err := ve.getComponent(entities[i])
		if err != nil {
			err2 = err
		}
		returnArray[i] = k
	}
	return returnArray, err2
}

//Writes Data to the specified entityID
func (ve *VectorStorage[T]) Write(entity EntityID, data T) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()

	if ve.exists(entity) {
		ve.internalVector[entity] = data
		return nil
	}
	return EntityNotFoundError
}

//Writes Data[i] to each EntityID[i]
func (ve *VectorStorage[T]) WriteMultiple(entities []EntityID, data []*T) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()
	if len(entities) != len(data) {
		return errors.New("Length mismatch between entities data")
	}

	for _, e := range entities {
		if !ve.exists(e) {
			return EntityNotFoundError
		}
	}

	for i, e := range entities {
		ve.internalVector[e] = *data[i]
	}
	return nil

}

//Appends an component to the end of the list
func (ve *VectorStorage[T]) AddEntity(entity EntityID, component T) error {
	ve.RWLOCK.Lock()
	defer ve.RWLOCK.Unlock()

	if ve.exists(entity) {
		return OneOrMoreEntitiesAlreadyExists
	}

	if int(entity) >= len(ve.internalVector) {
		togrow := int(entity+2) - len(ve.internalVector)
		ve.internalVector = append(ve.internalVector, make([]T, togrow, togrow)...)
		ve.allocated = append(ve.allocated, make([]bool, togrow, togrow)...)
	}

	ve.allocated[entity] = true
	ve.internalVector[entity] = component
	ve.numStored++
	return nil
}

//Appends components to the end of the list
//This will panic if the same entityID is listed multiple times
func (ve *VectorStorage[T]) AddEntityMultiple(Entitylist []EntityID, Components []T) error {
	ve.RWLOCK.RLock()
	defer ve.RWLOCK.RUnlock()

	if len(Entitylist) != len(Components) {
		return errors.New(fmt.Sprintf("Length of entity list must equal length of components %d != %d", len(Entitylist), len(Components)))
	}
	max := 0

	for _, e := range Entitylist {
		if int(e) > max {
			max = int(e)
		}
		if ve.allocated[e] {
			return OneOrMoreEntitiesAlreadyExists
		}
	}
	if max >= len(ve.internalVector) {
		togrow := (max + 1) - len(ve.internalVector)
		ve.internalVector = append(ve.internalVector, make([]T, togrow, togrow)...)
		ve.allocated = append(ve.allocated, make([]bool, togrow, togrow)...)
	}

	for i, v := range Entitylist {
		if ve.allocated[v] {
			panic("Attemped to insert a duplicate entity, this shouldnt happen")
		}
		ve.internalVector[v] = Components[i]
		ve.numStored++
	}
	return nil
}
