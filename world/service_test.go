package world

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jevans40/Ruthenium/component"
	"github.com/stretchr/testify/assert"
)

type NewService struct {
	*BaseService
}

type TestComponentHealth struct {
	Health int
}

func (t TestComponentHealth) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t TestComponentHealth) IsComponent()          {}

type TestComponentPosition struct {
	x float64
	y float64
	z float64
}

func (t TestComponentPosition) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t TestComponentPosition) IsComponent()          {}

type RequestedStruct[T any] struct {
}

func (s *NewService) Run(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) (err error) {
	var myHealthWrite component.WriteStorage[TestComponentHealth]
	var myPositionRead component.ReadOnlyStorage[TestComponentPosition]
	for i, datum := range s.dataPointers {
		if s.requiredData[i].DataType == NewComponentAccess[TestComponentHealth](ReadAccess).DataType {
			myHealthWrite, err = component.GetWriteStorage[TestComponentHealth](datum)
		}
		if s.requiredData[i].DataType == NewComponentAccess[TestComponentPosition](ReadAccess).DataType {
			myPositionRead, err = component.GetWriteStorage[TestComponentPosition](datum)
		}
	}
	myHealthWrite.Write(1, TestComponentHealth{88})
	fmt.Println(myHealthWrite.GetComponent(1))
	joined := component.Join(myHealthWrite, myPositionRead)
	fmt.Println(joined)

	return err
}

func TestThisService(t *testing.T) {
	myservice := NewService{&BaseService{Name: "NewService"}}
	myservice.SetRunFunction(myservice.Run)

	assert.Equal(t, myservice.GetName(), "NewService", "Reported service name does not match actual service name")

	//Add required AccessComponents
	myservice.AddRequiredAccessComponent(NewComponentAccess[TestComponentHealth](WriteAccess))
	myservice.AddRequiredAccessComponent(NewComponentAccess[TestComponentPosition](ReadAccess))
	err := myservice.AddRequiredAccessComponent(NewComponentAccess[TestComponentPosition](WriteAccess))
	compacc := []ComponentAccess{NewComponentAccess[TestComponentHealth](WriteAccess), NewComponentAccess[TestComponentPosition](ReadAccess)}
	assert.Equal(t, compacc, myservice.GetStorages(), "Accesscomponent array does not match storages report.")

	assert.Error(t, err, "Duplicate component access was added mistakenly")

	//Setup storages for the rest of this test
	healthStorage := component.NewDenseStorage[TestComponentHealth]()
	healthWrite, _ := component.GetWriteStorage[TestComponentHealth](healthStorage)
	EntIDs := []component.EntityID{1, 2, 4, 5, 6, 7}
	HealthObjects := []TestComponentHealth{{100}, {99}, {87}, {0}, {-1}, {99}}
	healthWrite.AddEntityMultiple(EntIDs, HealthObjects)

	//Setup storages for the rest of this test
	positionStorage := component.NewDenseStorage[TestComponentPosition]()
	positionWrite, _ := component.GetWriteStorage[TestComponentPosition](positionStorage)
	EntID2 := []component.EntityID{1, 2, 3, 4, 5, 7}
	PositionObjects := []TestComponentPosition{{0, 0, 0}, {100, 1, 1}, {9, 8, 6}, {87, 0, 1}, {0, 1, 1}, {99, -99, 99}}
	positionWrite.AddEntityMultiple(EntID2, PositionObjects)

	//Make communication channels
	call := make(chan error)
	entCreat := make(chan EntityCreationData)
	entDel := make(chan component.EntityID)

	go myservice.StartService(call)

	//Test storage updates
	err = myservice.UpdateStoragePointers([]component.ComponentStorage{healthStorage})
	assert.Error(t, err, "UpdateStoragePointers passed for an invalid configuration")
	err = myservice.UpdateStoragePointers([]component.ComponentStorage{healthStorage, positionStorage})
	assert.NoError(t, err, "UpdateStoragePointers failed for a valid configuration")

	assert.Equal(t, myservice.IsAsleep(), false, "Service falsely reported being asleep")
	myservice.SetSleepTime(10)
	assert.Equal(t, myservice.IsAsleep(), true, "Service falsely reported being awake")
	for myservice.IsAsleep() {
	}
	assert.Equal(t, myservice.IsAsleep(), false, "Service falsely reported being asleep")

	toSend := myservice.GetChannel()
	t.Log(toSend)
	close(toSend)
	toSend = myservice.GetChannel()

	go myservice.StartService(call)

	myservice.AddRequiredService("service1")
	myservice.AddRequiredService("service2")

	assert.Equal(t, []string{"service1", "service2"}, myservice.GetServices())

	toSend <- updateSignal{entCreat, entDel}
	err = <-call
	close(toSend)
	assert.NoError(t, err, "Error received from running service")
}

type comp1 struct {
}

func (t comp1) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp1) IsComponent()          {}

type comp2 struct {
}

func (t comp2) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp2) IsComponent()          {}

type comp3 struct {
}

func (t comp3) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp3) IsComponent()          {}

type comp4 struct {
}

func (t comp4) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp4) IsComponent()          {}

type comp5 struct {
}

func (t comp5) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp5) IsComponent()          {}

type comp6 struct {
}

func (t comp6) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp6) IsComponent()          {}

type comp7 struct {
}

func (t comp7) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp7) IsComponent()          {}

type comp8 struct {
}

func (t comp8) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp8) IsComponent()          {}

type comp9 struct {
}

func (t comp9) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp9) IsComponent()          {}

type comp10 struct {
}

func (t comp10) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t comp10) IsComponent()          {}

func TestAllocationTree(t *testing.T) {
	res1 := NewComponentAccess[comp1](ReadAccess)
	res2 := NewComponentAccess[comp2](ReadAccess)
	res3 := NewComponentAccess[comp3](WriteAccess)
	res4 := NewComponentAccess[comp4](ReadAccess)
	res5 := NewComponentAccess[comp5](WriteAccess)
	res6 := NewComponentAccess[comp6](ReadAccess)
	res7 := NewComponentAccess[comp7](ReadAccess)
	res8 := NewComponentAccess[comp8](WriteAccess)
	res9 := NewComponentAccess[comp9](ReadAccess)
	res10 := NewComponentAccess[comp10](WriteAccess)

	systema := []string{"d", "b", "c"}
	reqresa := []ComponentAccess{res1, res9, res7}
	systemb := []string{"d", "c"}
	reqresb := []ComponentAccess{res5, res2, res10}
	systemc := []string{}
	reqresc := []ComponentAccess{res3, res6, res8}
	systemd := []string{}
	reqresd := []ComponentAccess{res1, res4, res10}

	myTree := greedyAllocationTree{allocated: make(map[string]bool)}
	err := myTree.AddStystems([]string{"a", "b", "c", "d"}, [][]ComponentAccess{reqresa, reqresb, reqresc, reqresd}, [][]string{systema, systemb, systemc, systemd})
	assert.NoError(t, err, "Failed to add systems to tree")
	t.Log(myTree)
	t.Log(myTree.GetSystemTree())

}

func Run2(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) (err error) {
	EntityCreation <- EntityCreationData{10, []StorageWriteable{}, make(chan component.EntityID, 10)}
	return err
}

func TestDispatcher(t *testing.T) {
	res1 := NewComponentAccess[comp1](ReadAccess)
	res2 := NewComponentAccess[comp2](ReadAccess)
	res3 := NewComponentAccess[comp3](WriteAccess)
	res4 := NewComponentAccess[comp4](ReadAccess)
	res5 := NewComponentAccess[comp5](WriteAccess)
	res6 := NewComponentAccess[comp6](ReadAccess)
	res7 := NewComponentAccess[comp7](ReadAccess)
	res8 := NewComponentAccess[comp8](WriteAccess)
	res9 := NewComponentAccess[comp9](ReadAccess)
	res10 := NewComponentAccess[comp10](WriteAccess)
	t.Log("Created component access")

	store1 := component.NewDenseStorage[comp1]()
	store2 := component.NewDenseStorage[comp2]()
	store3 := component.NewDenseStorage[comp3]()
	store4 := component.NewDenseStorage[comp4]()
	store5 := component.NewDenseStorage[comp5]()
	store6 := component.NewDenseStorage[comp6]()
	store7 := component.NewDenseStorage[comp7]()
	store8 := component.NewDenseStorage[comp8]()
	store9 := component.NewDenseStorage[comp9]()
	store10 := component.NewDenseStorage[comp10]()
	t.Log("Created component storage")

	testingDispatcher := NewSimpleDispatcher()
	t.Log("Created simpleDispatcher")

	myservice1 := NewBaseService("a")
	myservice1.AddRequiredAccessComponent(res1)
	myservice1.AddRequiredAccessComponent(res2)
	myservice1.AddRequiredAccessComponent(res4)
	myservice1.AddRequiredAccessComponent(res3)
	myservice1.AddRequiredAccessComponent(res8)
	myservice1.AddRequiredAccessComponent(res7)
	myservice1.AddRequiredService("e")
	myservice2 := NewBaseService("b")
	myservice2.AddRequiredAccessComponent(res5)
	myservice2.AddRequiredAccessComponent(res7)
	myservice2.AddRequiredAccessComponent(res4)
	myservice2.AddRequiredService("a")
	myservice2.AddRequiredService("e")
	myservice2.SetRunFunction(Run2)
	myservice3 := NewBaseService("c")
	myservice3.AddRequiredAccessComponent(res9)
	myservice3.AddRequiredAccessComponent(res10)
	myservice2.AddRequiredService("e")
	myservice4 := NewBaseService("d")
	myservice4.AddRequiredAccessComponent(res1)
	myservice4.AddRequiredAccessComponent(res6)
	myservice4.AddRequiredAccessComponent(res9)
	myservice5 := NewBaseService("e")
	myservice5.AddRequiredAccessComponent(res1)
	myservice5.AddRequiredAccessComponent(res2)
	myservice5.AddRequiredAccessComponent(res3)
	t.Log("Created services")

	testingDispatcher.AddService(myservice1)
	testingDispatcher.RemoveService("a")
	testingDispatcher.AddService(myservice1)
	testingDispatcher.AddService(myservice2)
	testingDispatcher.AddService(myservice3)
	testingDispatcher.AddService(myservice4)
	testingDispatcher.AddService(myservice5)
	t.Log("Added Services")

	testingDispatcher.AddStorage(store1)
	testingDispatcher.RemoveStorage(store1.GetType())
	testingDispatcher.AddStorage(store1)
	testingDispatcher.AddStorage(store2)
	testingDispatcher.AddStorage(store3)
	testingDispatcher.AddStorage(store4)
	testingDispatcher.AddStorage(store5)
	testingDispatcher.AddStorage(store6)
	testingDispatcher.AddStorage(store7)
	testingDispatcher.AddStorage(store8)
	testingDispatcher.AddStorage(store9)
	testingDispatcher.AddStorage(store10)
	t.Log("Added Storages")

	testingDispatcher.StartServices()
	t.Log("Started Services")
	testingDispatcher.StopServices()
	t.Log("Stoped Services")

	testingDispatcher.Maintain()
	t.Log("Maintained")

}
