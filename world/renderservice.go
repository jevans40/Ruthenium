package world

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/linmath"
)

type renderService struct {
	BaseService
	renderChan chan []float32
	t1         time.Time
	t2         time.Time
	t3         time.Time
	dbgnm      int
}

func NewRenderService(renderChan chan []float32) Service {
	newRender := &renderService{t1: time.UnixMilli(0), t2: time.UnixMicro(0), t3: time.UnixMicro(0)}
	newRender.renderChan = renderChan
	newRender.Name = "renderer"
	newRender.SetRunFunction(newRender.RenderRun)
	newRender.AddRequiredAccessComponent(NewComponentAccess[Transform](ReadAccess))
	newRender.AddRequiredAccessComponent(NewComponentAccess[Coordinate](ReadAccess))
	newRender.AddRequiredAccessComponent(NewComponentAccess[Sprite](ReadAccess))

	return newRender
}

func (r *renderService) RenderRun(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error {
	time1 := time.Now()
	TransformRead, err1 := GetReadStorage[Transform](r)
	SpriteRead, err2 := GetReadStorage[Sprite](r)
	CoordRead, err3 := GetReadStorage[Coordinate](r)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	r.t1 = r.t1.Add(time.Since(time1))

	time2 := time.Now()
	Entities := component.Join(TransformRead, SpriteRead, CoordRead)
	r.t2 = r.t2.Add(time.Since(time2))

	RenderVec := make([]float32, len(Entities)*28)
	time3 := time.Now()
	Transforms, err1 := TransformRead.GetComponentMultiple(Entities)
	Sprites, err2 := SpriteRead.GetComponentMultiple(Entities)
	Coordinates, err3 := CoordRead.GetComponentMultiple(Entities)
	r.t3 = r.t3.Add(time.Since(time3))
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}

	if len(Entities) != 0 {
		var WorkerWait sync.WaitGroup
		WorkerWait.Add(4)
		for i := 0; i < 4; i++ {
			batchSize := len(Entities) / 4
			if i == 3 {
				go calculateVerticesWorker(i, len(Entities)-3*batchSize, Transforms[i*batchSize:], Sprites[i*batchSize:], Coordinates[i*batchSize:], RenderVec[i*batchSize*28:], &WorkerWait)
			} else {
				go calculateVerticesWorker(i, batchSize, Transforms[i*batchSize:(i+1)*batchSize], Sprites[i*batchSize:(i+1)*batchSize], Coordinates[i*batchSize:(i+1)*batchSize], RenderVec[i*batchSize*28:batchSize*28*(i+1)], &WorkerWait)
			}
		}
		//fmt.Println("Wait")
		WorkerWait.Wait()
	}
	select {
	case r.renderChan <- RenderVec:
	default:
		return errors.New("render channel full, render failed")
	}
	if r.dbgnm%100 == 99 {
		fmt.Printf("Renderer T1: %d, T2: %d, T3: %d\n", r.t1.UnixMilli()/100, r.t2.UnixMilli()/100, r.t3.UnixMilli()/100)
		r.t1 = time.UnixMilli(0)
		r.t2 = time.UnixMilli(0)
		r.t3 = time.UnixMilli(0)
	}
	r.dbgnm++
	return nil

}

func calculateVerticesWorker(dbg int, num int, Transforms []Transform, Sprites []Sprite, Coordinates []Coordinate, RenderVec []float32, wait *sync.WaitGroup) {
	//fmt.Printf("Worker %d started\n", dbg)
	defer wait.Done()
	for i := 0; i < num; i++ {
		//fmt.Printf("%d Out of %d\n", dbg*num+i, num*4)
		calculateVertices(Transforms[i], Sprites[i], Coordinates[i], RenderVec[i*28:i*28+28])
	}
	//fmt.Println("Done")
}

func calculateVertices(transform Transform, sprite Sprite, coordinates Coordinate, vector []float32) {

	vert := make([]linmath.Vertice, 4)
	for i := range vert {
		vert[i] = linmath.EmptyVertice()
	}
	Xindex := [4]float64{-0.5, 0.5, -0.5, 0.5}
	Yindex := [4]float64{-0.5, -0.5, 0.5, 0.5}

	for i := 0; i < 4; i++ {
		vert[i].SetColor(sprite.Color)
		vert[i].SetMap(uint32(sprite.M))
		vert[i].SetTexX(sprite.X + sprite.W*float32((i)%2))
		vert[i].SetTexY(sprite.Y + sprite.H*float32(int32((i)/2)%2))
		transforms := linmath.Mat3f[float64](transform).VectorMul(Xindex[i], Yindex[i], 1)
		//(0,0),(1,0),(0,1)(1,1)
		//(-1,-1),(1,-1),(-1,1),(1,1)
		vert[i].SetX(float32(coordinates[0] + transforms[0]))
		vert[i].SetY(float32(coordinates[1] + transforms[1]))
		vert[i].SetZ(float32(coordinates[2]))
		//log.WithFields(log.Fields{"Vertnum": i, "X": vert[i].GetX(), "Y": vert[i].GetY(), "Z": vert[i].GetZ(), "TexX": vert[i].GetTexX(), "TexY": vert[i].GetTexY(), "Color": vert[i].GetColor()}).Trace()
	}

	for i, v := range vert {
		floats := v.ToFloats()
		copy(vector[i*7:7*i+7], floats[:])
	}
}
