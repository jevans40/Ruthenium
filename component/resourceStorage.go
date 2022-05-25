package component

import "reflect"

var _ ComponentStorage = &ResourceStorage[BaseComponent]{}
var _ ReadOnlyStorage[BaseComponent] = &ResourceStorage[BaseComponent]{}
var _ WriteStorage[BaseComponent] = &ResourceStorage[BaseComponent]{}

//Resource Storage struct, stores a single object at entity ID 0
//This is used for storing global varibles
type ResourceStorage[T Component] struct {
	resource T
}

func NewResourceStorage[T Component](resource T) ComponentStorage {
	return &ResourceStorage[T]{resource: resource}
}

//Get the type of the resource
func (r *ResourceStorage[T]) GetType() reflect.Type {
	return reflect.TypeOf(r.resource)
}

//Returns false, this is a resource it only stores the resource, no entities
func (r *ResourceStorage[T]) Exists(EntityID) bool {
	return false
}

//Returns [false], this is a resource it only stores the resource, no entitites
func (r *ResourceStorage[T]) ExistsMultiple(Entities []EntityID) []bool {
	return []bool{false}
}

//Returns [-1], the resource contains no entities.
func (r *ResourceStorage[T]) GetEntities() []EntityID {
	return []EntityID{-1}
}

//Returns 1, only contains the resource
func (r *ResourceStorage[T]) GetSize() int {
	return 1
}

//Returns NotEntityStorage error
func (r *ResourceStorage[T]) AddBlankComponent(e EntityID) error {
	return NotEntityStorageError
}

//Returns NotEntityStorage error
func (r *ResourceStorage[T]) AddBlankComponentMultiple(e []EntityID) error {
	return NotEntityStorageError
}

//Retursn NotEntityStorage error
func (r *ResourceStorage[T]) DeleteEntityMultiple(e []EntityID) error {
	return NotEntityStorageError
}

//Returns NotEntityStorage error
func (r *ResourceStorage[T]) DeleteEntity(entity EntityID) error {
	return NotEntityStorageError
}

//Returns the resource
func (r *ResourceStorage[T]) MustGetComponent(e EntityID) T {
	return r.resource
}

func (r *ResourceStorage[T]) MustGetComponentMultiple(e []EntityID) []T {
	return []T{r.resource}
}

func (r *ResourceStorage[T]) GetComponent(entity EntityID) (T, error) {
	return r.resource, nil
}

func (r *ResourceStorage[T]) GetComponentMultiple(entities []EntityID) ([]*T, error) {
	return []*T{&r.resource}, nil
}

func (r *ResourceStorage[T]) Write(entity EntityID, data T) error {
	r.resource = data
	return nil
}

func (r *ResourceStorage[T]) WriteMultiple(entities []EntityID, data []*T) error {
	return NotEntityStorageError
}

func (r *ResourceStorage[T]) AddEntity(entity EntityID, component T) error {
	return NotEntityStorageError
}

func (r *ResourceStorage[T]) AddEntityMultiple(Entitylist []EntityID, Components []T) error {
	return NotEntityStorageError
}
