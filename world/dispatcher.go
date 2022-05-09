package world

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/constants"
	"github.com/jevans40/Ruthenium/ruthutil"
	log "github.com/sirupsen/logrus"
)

/*
This is probably the most complex object in the whole system, but if it works then it this project will be a sucess.
*/

type Dispatcher interface {

	//This method will both maintain all updates queued for storages
	//It will run all services currently on this dispatcher for one step
	Maintain() error

	//Adds a service to the list of managed services.
	//Services cannot have duplicate names, but they
	//Can otherwise be complete duplicates
	AddService(service Service) error

	//Adds a new component storage to the list of managed
	//storages. Can only accept one type per storage.
	AddStorage(storage component.ComponentStorage) error

	//IMPLEMENT: Add resource
	//Resources should be single instances of data shared at the
	//World layer. One way to implement this is to make all
	//getType functions return a string instead of a type.
	//AddResource(resource) err

	//Removes a service with the given name
	RemoveService(serviceName string) error

	//Removes a storage with the given type
	RemoveStorage(storage reflect.Type) error

	//Starts all internal services, this is called
	//At the start of each maintain loop.
	StartServices() error

	//Stops all managed and internal services.
	StopServices() error
}

var _ Dispatcher = &simpleDispatcher{}

/***************************/
/*    simpleDispatcher Methods    */

type simpleDispatcher struct {
	//
	entities  map[component.EntityID]component.Entity
	entityNum int

	running bool

	storages         []component.ComponentStorage
	services         []Service
	serviceCallbacks []chan updateSignal
	toDelete         []component.EntityID

	//Channels
	errorChannel    chan error
	entityCreations chan EntityCreationData
	entityDeletions chan component.EntityID

	//Mutex
	entityWrite     sync.Mutex
	entityProcessed sync.WaitGroup

	t1    time.Time
	t2    time.Time
	t3    time.Time
	dbgnm int
}

func NewSimpleDispatcher() Dispatcher {
	return &simpleDispatcher{entities: make(map[component.EntityID]component.Entity),
		errorChannel:    make(chan error, 100*constants.RACECHANNELSIZETEST),
		entityCreations: make(chan EntityCreationData, 100*constants.RACECHANNELSIZETEST),
		entityDeletions: make(chan component.EntityID, 100*constants.RACECHANNELSIZETEST),
		running:         false,
		t1:              time.UnixMilli(0),
		t2:              time.UnixMilli(0),
		t3:              time.UnixMilli(0),
		dbgnm:           0}
}

func (d *simpleDispatcher) Maintain() error {
	//check that all internal services are running
	if !d.running {
		d.StartServices()
	}

	d.entityProcessed.Add(2)
	if ruthutil.IsChannelClosed(d.entityCreations) {
		d.entityCreations = make(chan EntityCreationData, 100*constants.RACECHANNELSIZETEST)
	}

	if ruthutil.IsChannelClosed(d.entityDeletions) {
		d.entityDeletions = make(chan component.EntityID, 100*constants.RACECHANNELSIZETEST)
	}

	go d.startEntityCreationService()
	go d.startEntityDeletionService()

	//GetService Requirements and start them

	var comChannels []chan updateSignal
	var resReq [][]ComponentAccess
	var servReq [][]string
	var servicesNames []string

	for _, v := range d.services {
		servicesNames = append(servicesNames, v.GetName())
		comChannels = append(comChannels, v.GetChannel())
		resReq = append(resReq, v.GetStorages())
		servReq = append(servReq, v.GetServices())
	}

	newTree := greedyAllocationTree{allocated: make(map[string]bool)}
	err := newTree.AddStystems(servicesNames, resReq, servReq)
	if err != nil {
		return err
	}
	serviceOrder := newTree.GetSystemTree()

	time1 := time.Now()
	//fmt.Println(serviceOrder)
	for _, batch := range serviceOrder {
		for i, s := range d.services {
			time2 := time.Now()
			for _, c := range batch {
				if s.GetName() == c {
					var toUpdate []component.ComponentStorage
					for _, k := range resReq[i] {
						for _, j := range d.storages {
							if k.DataType == j.GetType() {
								toUpdate = append(toUpdate, j)
							}
						}
					}
					s.UpdateStoragePointers(toUpdate)
					go s.StartService(d.errorChannel, updateSignal{d.entityCreations, d.entityDeletions})
					//fmt.Printf("sent for service %s \n", s.GetName())
				}

			}
			d.t2 = d.t2.Add(time.Since(time2))
		}
		time3 := time.Now()
		for i := 0; i < len(batch); i++ {
			//fmt.Println(len(batch))
			//fmt.Println(i)
			err, ok := <-d.errorChannel
			//log.Error(err)
			if !ok {
				panic("error channel closed unexpectedly")
			}
			if err != nil {
				//log.Error(err)
			}
		}
		d.t3 = d.t3.Add(time.Since(time3))
	}
	d.t1 = d.t1.Add(time.Since(time1))

	close(d.entityCreations)
	close(d.entityDeletions)
	d.entityProcessed.Wait()

	//TODO: Optimization: So this should be reformatted to create a smarter deletion process. This is a time consuming part of the update loop,
	//But this is the first thing I thought of.
	//var toRemove []component.EntityID
	//for i, entity := range d.entities {
	//	if entity.Deleted {

	//		for _, storage := range d.storages {
	//			storage.DeleteEntity(entity.EntityNum)
	//		}
	//		toRemove = append(toRemove, i)
	//	}
	//}

	for _, i := range d.toDelete {
		delete(d.entities, i)
	}
	d.toDelete = []component.EntityID{}
	if d.dbgnm%100 == 99 {
		fmt.Printf("Dispatcher T1: %d, T2: %d, T3: %d\n", d.t1.UnixMilli()/100, d.t2.UnixMilli()/100, d.t3.UnixMilli()/100)
		d.t1 = time.UnixMilli(0)
		d.t2 = time.UnixMilli(0)
		d.t3 = time.UnixMilli(0)
	}
	d.dbgnm++
	return nil
}

func (d *simpleDispatcher) AddService(newService Service) error {
	thisName := newService.GetName()
	for _, s := range d.services {
		if thisName == s.GetName() {
			return errors.New("service already exists in this Dispatcher")
		}
	}
	d.services = append(d.services, newService)
	d.serviceCallbacks = append(d.serviceCallbacks, newService.GetChannel())
	return nil
}

func (d *simpleDispatcher) AddStorage(newStorage component.ComponentStorage) error {
	val := newStorage.GetType()
	for _, v := range d.storages {
		if v.GetType() == val {
			return errors.New("component type already exists in this Dispatcher")
		}
	}
	d.storages = append(d.storages, newStorage)
	return nil
}

func (d *simpleDispatcher) RemoveService(name string) error {
	for i, s := range d.services {
		if name == s.GetName() {
			d.services = append(d.services[:i], d.services[i+1:]...)
			return nil
		}
	}
	return errors.New("service not found in this Dispatcher")
}

func (d *simpleDispatcher) RemoveStorage(thisType reflect.Type) error {
	for i, s := range d.storages {
		if s.GetType() == thisType {
			d.storages = append(d.storages[:i], d.storages[i+1:]...)
			return nil
		}
	}
	return errors.New("storage not found in this Dispatcher")
}

func (d *simpleDispatcher) StartServices() error {
	/*
		if d.running {
			return errors.New("service already running")
		}
		d.running = true
		for _, s := range d.services {
			go s.StartService(d.errorChannel)
		}
		//Ensure that all services are ready before moving on
		for i := 0; i < len(d.services); i++ {
			<-d.errorChannel
		}
	*/
	return nil
}

func (d *simpleDispatcher) StopServices() error {
	if d.running {
		return errors.New("services already stopped")
	}
	d.running = false
	for _, s := range d.services {
		close(s.GetChannel())
	}
	return nil
}

//Start this function in a seperate goroutine to handle entity deletion/creation requests.
//Will set a writeEntity mutex lock.
func (d *simpleDispatcher) startEntityCreationService() {
	var toAdd []int
	var toAddComp [][]StorageWriteable
	var channels []chan component.EntityID
	for {
		ent, err := ruthutil.WaitChannel(d.entityCreations)
		if err != nil {
			break
		}
		if ent.NumEntities < 0 {
			continue
		}
		toAdd = append(toAdd, ent.NumEntities)
		toAddComp = append(toAddComp, ent.Components)
		channels = append(channels, ent.CreatedEntitiesCallback)
	}
	d.entityWrite.Lock()
	for i, v := range channels {
		for j := 0; j < toAdd[i]; j++ {
			newID := component.EntityID(d.entityNum)
			d.entities[newID] = component.Entity{newID, false}
			for _, k := range toAddComp[i] {
				k.AddComponentWithEntityID(newID)
			}
			select {

			case v <- newID:

			default:
				log.Info("attempted to send an entity to a full channel")
			}
			d.entityNum++
		}
	}
	d.entityWrite.Unlock()
	d.entityProcessed.Done()
}

//TODO: Add a lazy Componnent Addition/Deletion Service

func (d *simpleDispatcher) startEntityDeletionService() {
	var toDelete []component.EntityID
	for {
		delete, err := ruthutil.WaitChannel(d.entityDeletions)
		if err != nil {
			break
		}

		if delete < 0 {
			continue
		}
		toDelete = append(toDelete, delete)
	}
	d.entityWrite.Lock()
	for _, v := range toDelete {
		d.entities[v] = component.Entity{d.entities[v].EntityNum, true}
		d.toDelete = append(d.toDelete, v)
	}
	d.entityWrite.Unlock()
	d.entityProcessed.Done()
}

/***************************/
/* Greedy Allocation Tree  */
type branchLevel struct {
	WrittenTypes     []reflect.Type
	ReadTypes        []reflect.Type
	AllocatedSystems []string
	NextLevel        *branchLevel
}

func (b *branchLevel) AddLevel() {
	b.NextLevel = &branchLevel{}
}

func (b *branchLevel) AddSystem(system string, res []ComponentAccess, req []string) bool {
	for _, s := range req {
		for _, k := range b.AllocatedSystems {
			if s == k {
				if b.NextLevel == nil {
					b.NextLevel = &branchLevel{}
				}
				return b.NextLevel.AddSystem(system, res, req)
			}
		}
	}
	for _, r := range res {
		//Check if write and read are already allocated
		for _, v := range b.WrittenTypes {
			if v == r.DataType {
				if b.NextLevel == nil {
					b.NextLevel = &branchLevel{}
				}
				return b.NextLevel.AddSystem(system, res, req)
			}
		}
		if r.Access == WriteAccess {
			for _, v := range b.ReadTypes {
				if v == r.DataType {
					if b.NextLevel == nil {
						b.NextLevel = &branchLevel{}
					}
					return b.NextLevel.AddSystem(system, res, req)
				}
			}
		}

	}
	//This level is available for allocation
	b.AllocatedSystems = append(b.AllocatedSystems, system)
	for _, r := range res {
		if r.Access == WriteAccess {
			b.WrittenTypes = append(b.WrittenTypes, r.DataType)
		} else {
			for _, v := range b.ReadTypes {
				if r.DataType == v {
					goto NextRes
				}
			}
			b.ReadTypes = append(b.ReadTypes, r.DataType)
		}
	NextRes:
	}
	return true
}

func (b *branchLevel) GetSystemTree(tree [][]string) [][]string {
	tree = append(tree, b.AllocatedSystems)
	if b.NextLevel != nil {
		return b.NextLevel.GetSystemTree(tree)
	}
	return tree
}

type greedyAllocationTree struct {
	allocated  map[string]bool
	firstLevel branchLevel
}

func (g *greedyAllocationTree) AddStystems(systems []string, resources [][]ComponentAccess, sysReq [][]string) error {
	//OuterMost can only run at most len(system) times
	for k := 0; k < len(systems); k++ {
		for i, s := range systems {
			if _, ok := g.allocated[s]; ok {
				continue
			}
			for _, r := range sysReq[i] {
				if _, ok := g.allocated[r]; !ok {
					goto failedToFind
				}
			}
			g.allocated[s] = true
			g.firstLevel.AddSystem(s, resources[i], sysReq[i])
		failedToFind:
		}
		if len(systems) == len(g.allocated) {
			return nil
		}
	}
	return errors.New("found a cycle within the allocation tree allocation failed")

}

func (g *greedyAllocationTree) GetSystemTree() [][]string {
	var tree [][]string
	return g.firstLevel.GetSystemTree(tree)
}

/***************************/
/*    Generic Functions    */
