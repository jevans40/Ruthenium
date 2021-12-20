package component

import (
	"errors"
	"fmt"
	"reflect"
)

var _ ComponentStorage = &DenseStorage[BaseComponent]{}
var _ ReadOnlyStorage[BaseComponent] = &DenseStorage[BaseComponent]{}
var _ WriteStorage[BaseComponent] = &DenseStorage[BaseComponent]{}

//The Dense Storage struct:
//This type of storage stores all components in a dense array.
//It uses a map from EntityID's to the internal storage map for lookup.
type DenseStorage[T Component] struct {
	component   []T
	internalMap map[EntityID]int
}

//Create a new Dense Storage containing types T.
//Returns a ComponentStorage interface
func NewDenseStorage[T Component]() ComponentStorage {
	return &DenseStorage[T]{component: []T{}, internalMap: map[EntityID]int{}}
}

//Returns the type of the contained storage
func (d *DenseStorage[T]) GetType() reflect.Type {
	return reflect.TypeOf(d.component).Elem()
}

//Returns true if the entitity is stored in the internal map
func (d *DenseStorage[T]) Exists(Entity EntityID) bool {
	_, ok := d.internalMap[Entity]
	return ok
}

//Returns a mask of the entities that exist in this storage
func (d *DenseStorage[T]) ExistsMultiple(Entities []EntityID) []bool {
	toReturn := []bool{}
	for _, v := range Entities {
		toReturn = append(toReturn, d.Exists(v))
	}
	return toReturn
}

//Return all stored entities in this storage
func (d *DenseStorage[T]) GetEntities() []EntityID {
	toReturn := []EntityID{}
	for k, _ := range d.internalMap {
		toReturn = append(toReturn, k)
	}
	return toReturn
}

//Returns the number of components stored in this storage
func (d *DenseStorage[T]) GetSize() int {
	return len(d.component)
}

//Adds a new empty component to this storage
func (d *DenseStorage[T]) AddBlankComponent(Entity EntityID) error {
	if _, ok := d.internalMap[Entity]; ok {
		return OneOrMoreEntitiesAlreadyExists
	}
	var newComp T
	d.component = append(d.component, newComp)
	d.internalMap[Entity] = len(d.component)
	return nil
}

//Adds a new Entities with associated IDs to this storage
//Returns an error if one or more entitys already exist in the storage
//If an error is returned no entities are added to the storage
//This will panic if entityID is listed multiple times
func (d *DenseStorage[T]) AddBlankComponentMultiple(Entity []EntityID) error {
	for _, e := range Entity {
		if _, ok := d.internalMap[e]; ok {
			return OneOrMoreEntitiesAlreadyExists
		}
	}
	for _, e := range Entity {
		if _, ok := d.internalMap[e]; ok {
			panic("Attemped to insert a duplicate entity, this shouldnt happen")
		}
		var newComp T
		d.component = append(d.component, newComp)
		d.internalMap[e] = len(d.component)
	}
	return nil
}

//Deletes selected elements from storage
//This take O(m + n) time since it iterates over all entities and then all interal maps
//TODO: If you want single or small number deletion use the DeleteEntity function instead.
//This should only be called by the manager of the ComponentStorage.
func (d *DenseStorage[T]) DeleteEntityMultiple(Entities []EntityID) error {
	for _, e := range Entities {
		if _, ok := d.internalMap[e]; !ok {
			return EntityNotFoundError
		}
		d.internalMap[e] = -1
	}
	newStorage := []T{}
	newMap := map[EntityID]int{}
	for k, v := range d.internalMap {
		if v != -1 {
			newStorage = append(newStorage, d.component[v])
			newMap[k] = len(newStorage) - 1
		}
	}
	d.component = newStorage
	d.internalMap = newMap
	return nil
}

//The singlecase version of DeleteEntities Multiple
//TODO: Should be avoided for now until fixed
func (d *DenseStorage[T]) DeleteEntity(entity EntityID) error {
	return d.DeleteEntityMultiple([]EntityID{entity})
}

//Returns struct copies of entities from the storage
//Note: This will panic if it cannot find the given entityID
//Call GetComponent if you want to just receive an error
func (d *DenseStorage[T]) MustGetComponent(entity EntityID) T {
	if val, ok := d.internalMap[entity]; ok {
		return d.component[val]
	}
	panic(EntityNotFoundError)
}

//Calls MustGetComponent on all entities listed
func (d *DenseStorage[T]) MustGetComponentMultiple(entities []EntityID) []T {
	returnArray := []T{}
	for _, v := range entities {
		k := d.MustGetComponent(v)
		returnArray = append(returnArray, k)
	}
	return returnArray
}

//Calls GetComponent on all entities listed
func (d *DenseStorage[T]) GetComponent(entity EntityID) (T, error) {

	if val, ok := d.internalMap[entity]; ok {
		return d.component[val], nil
	}

	var errorFound T
	return errorFound, EntityNotFoundError
}

//Returns struct copies of the requested entities from storage
func (d *DenseStorage[T]) GetComponentMultiple(entities []EntityID) ([]T, error) {
	returnArray := []T{}
	err2 := error(nil)
	for _, v := range entities {
		k, err := d.GetComponent(v)
		if err != nil {
			err2 = err
		}
		returnArray = append(returnArray, k)
	}
	return returnArray, err2
}

//Writes Data to the specified entityID
func (d *DenseStorage[T]) Write(entity EntityID, data T) error {
	if val, ok := d.internalMap[entity]; ok {
		d.component[val] = data
		return nil
	}
	return EntityNotFoundError
}

//Writes Data[i] to each EntityID[i]
func (d *DenseStorage[T]) WriteMultiple(entities []EntityID, data []T) error {
	if len(entities) != len(data) {
		return errors.New("Length mismatch between entities data")
	}

	for _, e := range entities {
		if _, ok := d.internalMap[e]; !ok {
			return EntityNotFoundError
		}
	}

	for i, e := range entities {
		d.component[d.internalMap[e]] = data[i]
	}
	return nil

}

//Appends an component to the end of the list
func (d *DenseStorage[T]) AddEntity(entity EntityID, component T) error {

	if _, ok := d.internalMap[entity]; ok {
		return OneOrMoreEntitiesAlreadyExists
	}

	d.component = append(d.component, component)

	d.internalMap[entity] = len(d.component) - 1

	return nil
}

//Appends components to the end of the list
//This will panic if entityID is listed multiple times
func (d *DenseStorage[T]) AddEntityMultiple(Entitylist []EntityID, Components []T) error {
	if len(Entitylist) != len(Components) {
		return errors.New(fmt.Sprintf("Length of entity list must equal length of components %d != %d", len(Entitylist), len(Components)))
	}

	startingSize := len(d.component)
	for _, e := range Entitylist {
		if _, ok := d.internalMap[e]; ok {
			return OneOrMoreEntitiesAlreadyExists
		}
	}

	d.component = append(d.component, Components...)

	for i, v := range Entitylist {
		if _, ok := d.internalMap[v]; ok {
			panic("Attemped to insert a duplicate entity, this shouldnt happen")
		}
		d.internalMap[v] = startingSize + i
	}
	return nil
}
