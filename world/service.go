package world

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/constants"
	"github.com/jevans40/Ruthenium/ruthutil"
	log "github.com/sirupsen/logrus"
)

//TODO: Service needs to be broken into [Service Handler] -> [Service] Instead of just having a run function
//Service handler will handle syncronization and getting resources from dispatcher
//Service will only handle running the service. It will have a single callable function:
//Run(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error

//TODO: All of these need to be placed somewhere else. Putting them in the header of the service was a temporary
//work around, but now that other classes are using this we need to move them to a reasonable spot.

const (
	WriteAccess AccessType = iota
	ReadAccess
)

type AccessType int

type StorageWriteable interface {
	AddComponentWithEntityID(component.EntityID)
}

//This function is here because golang doesnt allow Generic Methods ;_;
func MakeWriteableStorage[T component.Component](NewComponent T, writeableStorage component.WriteStorage[T]) StorageWriteable {
	return &ComponentWriteDate[T]{NewComponent, writeableStorage}
}

type ComponentWriteDate[T component.Component] struct {
	component T
	storage   component.WriteStorage[T]
}

func (c *ComponentWriteDate[T]) AddComponentWithEntityID(ID component.EntityID) {
	c.storage.AddEntity(ID, c.component)
}

type EntityCreationData struct {
	NumEntities             int
	Components              []StorageWriteable
	CreatedEntitiesCallback chan component.EntityID
}

//TODO:
//REFACTOR:
//EntityCreation: signals on this channel will request entities to be made lazily
//				  The dispatcher should close this included channel and not block if the
//         		  channel buffer is too small to fit all the requested entities
//				  Form (NumEntities to make, Channel to receive entity ID's from later)
//EntityDeletion: signals on this channel will tell the dispatcher to lazily delete requested entities
type updateSignal struct {
	EntityCreation chan EntityCreationData
	EntityDeletion chan component.EntityID
}

type ComponentAccess struct {
	DataType reflect.Type
	Access   AccessType
}

func NewComponentAccess[T component.Component](access AccessType) ComponentAccess {
	return ComponentAccess{DataType: reflect.TypeOf((*T)(nil)).Elem(), Access: access}
}

type Service interface {

	//The function that will be called with its own go routine.
	//To stop the service simply close the associated channel
	//Channel descriptions:
	//Callback: this service will report any encountered errors to Callback every iteration |Blocking|
	StartService(Callback chan error, update updateSignal)

	//Return required datatypes for this service
	GetStorages() []ComponentAccess

	//Return required Services to run before this one
	GetServices() []string

	//Return the service name
	GetName() string

	//The service will listen for messages on this channel in order to start
	//a run cycle and call Run()
	GetChannel() chan updateSignal

	//Updates the required storage pointers, returns an error if missing a storage.
	UpdateStoragePointers([]component.ComponentStorage) error

	//Should be mutex locked, returns if the service should be run at next step
	IsAsleep() bool

	//Sets the thread to sleep for sleepTime iterations
	SetSleepTime(sleepTime int)

	//The function to overload, service code should be written here
	//This should be a method that your service implements
	//It will run once per update
	//Should return any errors that need to be handled higher up
	SetRunFunction(func(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error)

	//Adds a component type that this service requires data from
	//This request will be served by the dispatcher. The only garentees
	//Are that the datatype will match the ComponentAccess type,
	//And write and read access will be atleast equivelent.
	//There is no garentee that you cannot convert a read only into
	//a write however it may panic, and could lead to a datarace.
	AddRequiredAccessComponent(newComp ComponentAccess) error

	//Adds a service that should be run before this one each update cycle
	AddRequiredService(service string)

	//Returns a componentStorage that this service has access to associated with storageType
	//Returns the accesstype of the associated storage
	//Returns nil if the storage is missing
	GetStorage(storageType reflect.Type) (component.ComponentStorage, AccessType)
}

func NewBaseService(name string) Service {
	toReturn := &BaseService{Name: name, communicationChannel: make(chan updateSignal, 100*constants.RACECHANNELSIZETEST), sleepTime: 0}
	toReturn.SetRunFunction(toReturn.Run)

	return toReturn
}

type BaseService struct {
	Name string

	requiredData         []ComponentAccess
	dataPointers         []component.ComponentStorage
	requiredServices     []string
	sleepTime            int
	communicationChannel chan updateSignal

	runFunc func(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error

	SleepLock   sync.Mutex
	StorageLock sync.Mutex
	ComLock     sync.Mutex
}

/***************************/
/*    Service Methods      */

//The function that will be called with its own go routine.
//To stop the service simply close the associated channel
//Channel descriptions:
//Callback: this service will report any encountered errors to Callback every iteration |Blocking|
//EntityCreation: signals on this channel will request entities to be made lazily
//				  The dispatcher should close this included channel and not block if the
//         		  channel buffer is too small to fit all the requested entities
//				  Form (NumEntities to make, Channel to receive entity ID's from later)
//EntityDeletion: signals on this channel will tell the dispatcher to lazily delete requested entities

//TODO:: This should loop and possibly be synced with a WorkGroup.
func (s *BaseService) StartService(Callback chan error, update updateSignal) {
	s.GetChannel()
	//fmt.Printf("got a signal for service %s\n", s.Name)
	s.SleepLock.Lock()
	s.StorageLock.Lock()
	err := s.runFunc(update.EntityCreation, update.EntityDeletion)
	Callback <- err
	s.SleepLock.Unlock()
	s.StorageLock.Unlock()

}

//TODO: This function should accept component Creation and Deletion Events
func (s *BaseService) SetRunFunction(toRun func(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error) {
	s.runFunc = toRun
}

//This function should be overloaded for your service code.
//At this point all locks should be set so you can freely use any data in the
//Sleep or Storage. If there is a datarace at this point, something went wrong.
func (s *BaseService) Run(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error {
	return nil
}

func (s *BaseService) GetStorages() []ComponentAccess {
	return s.requiredData
}

func (s *BaseService) GetServices() []string {
	return s.requiredServices
}

//Returns service name
func (s *BaseService) GetName() string {
	return s.Name
}

//Returns this services communicationChannel
//If the channel is closed, a new one will be created
func (s *BaseService) GetChannel() chan updateSignal {
	s.ComLock.Lock()
	defer s.ComLock.Unlock()
	if ruthutil.IsChannelClosed(s.communicationChannel) || s.communicationChannel == nil {
		s.communicationChannel = make(chan updateSignal, 100*constants.RACECHANNELSIZETEST)
	}
	return s.communicationChannel
}

func (s *BaseService) UpdateStoragePointers(data []component.ComponentStorage) error {
	s.StorageLock.Lock()
	defer s.StorageLock.Unlock()
	//O(n^2) but all elements should be VERY small (<16) so its okay
	for i, v := range s.requiredData {
		tofind := v.DataType
		found := false
		for _, t := range data {
			//fmt.Println(&t)
			if t.GetType() == tofind {
				s.dataPointers[i] = t
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing Required Datatype %s", v.DataType)
		}
	}
	return nil
}

//Mutex locked
//Checks wether or not the service should run, and if not decriments the counter.
//While checking will never make this counter go negative, a negative counter will
//never wake up.
func (s *BaseService) IsAsleep() bool {
	s.SleepLock.Lock()
	defer s.SleepLock.Unlock()
	if s.sleepTime != 0 {
		s.sleepTime--
	}
	return s.sleepTime != 0
}

//Mutex locked
//Tells the service to sleep for x ammount of time.
//Setting this to a negative value will cause the service to sleep until changed to positive or zero.
func (s *BaseService) SetSleepTime(sleepTime int) {
	s.SleepLock.Lock()
	defer s.SleepLock.Unlock()
	s.sleepTime = sleepTime
}

func (s *BaseService) AddRequiredAccessComponent(newComp ComponentAccess) error {
	s.StorageLock.Lock()
	defer s.StorageLock.Unlock()
	for _, v := range s.requiredData {
		if v.DataType == newComp.DataType {
			return errors.New("component type already exists for this service, unable to add duplicate type")
		}
	}
	s.requiredData = append(s.requiredData, newComp)
	s.dataPointers = append(s.dataPointers, nil)
	return nil
}

func (s *BaseService) AddRequiredService(newService string) {
	s.requiredServices = append(s.requiredServices, newService)
}

func (s *BaseService) GetStorage(storageType reflect.Type) (component.ComponentStorage, AccessType) {
	var storage component.ComponentStorage
	var access AccessType
	for i, datum := range s.dataPointers {
		if s.requiredData[i].DataType == storageType {
			storage = datum
			break
		}
	}
	for _, t := range s.requiredData {
		if t.DataType == storageType {
			access = t.Access
		}
	}

	return storage, access
}

/***************************/
/*    Typed Functions    */

func GetWriteStorage[T component.Component](this Service) (component.WriteStorage[T], error) {
	storage, access := this.GetStorage(component.ReflectType[T]())
	if access != WriteAccess {
		return nil, errors.New("component accessType is not write")
	}
	if storage == nil {
		return nil, errors.New("component storage is not found in this service")
	}
	return component.GetWriteStorage[T](storage)
}

func GetReadStorage[T component.Component](this Service) (component.ReadOnlyStorage[T], error) {
	storage, access := this.GetStorage(component.ReflectType[T]())
	if access != ReadAccess {
		log.WithFields(log.Fields{"accessType": access, "Expected": ReadAccess}).Error("incompatible accessType")
		return nil, errors.New("component accessType is not write")
	}
	if storage == nil {
		return nil, errors.New("component storage is not found in this service")
	}
	return component.GetReadOnlyStorage[T](storage)
}
