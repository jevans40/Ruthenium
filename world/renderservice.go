package world

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/linmath"
)

//TODO:: Documentation
//TODO:: Tests

type renderService struct {
	BaseService
	renderChan chan []float32
	t1         time.Time
	t2         time.Time
	t3         time.Time
	t4         time.Time
	dbgnm      int
}

func NewRenderService(renderChan chan []float32) Service {
	newRender := &renderService{t4: time.UnixMilli(0), t1: time.UnixMilli(0), t2: time.UnixMicro(0), t3: time.UnixMicro(0)}
	newRender.renderChan = renderChan
	newRender.Name = "renderer"
	newRender.SetRunFunction(newRender.RenderRun)
	newRender.AddRequiredAccessComponent(NewComponentAccess[Renderable](ReadAccess))

	return newRender
}

func (r *renderService) RenderRun(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error {
	RenderableRead, err1 := GetReadStorage[Renderable](r)

	if err1 != nil {
		return err1
	}

	Entities := RenderableRead.GetEntities()

	RenderVec := make([]float32, len(Entities)*28)
	time3 := time.Now()
	Renderables, err1 := RenderableRead.GetComponentMultiple(Entities)
	r.t3 = r.t3.Add(time.Since(time3))
	if err1 != nil {
		return err1
	}

	time4 := time.Now()
	if len(Entities) != 0 {
		var WorkerWait sync.WaitGroup
		WorkerWait.Add(6)
		for i := 0; i < 6; i++ {
			batchSize := len(Entities) / 6
			if i == 5 {
				go calculateVerticesWorker(i, len(Entities)-5*batchSize, Renderables[i*batchSize:], RenderVec[i*batchSize*28:], &WorkerWait)
			} else {
				go calculateVerticesWorker(i, batchSize, Renderables[i*batchSize:(i+1)*batchSize], RenderVec[i*batchSize*28:batchSize*28*(i+1)], &WorkerWait)
			}
		}
		//fmt.Println("Wait")
		WorkerWait.Wait()
	}
	r.t4 = r.t4.Add(time.Since(time4))

	time1 := time.Now()
	select {
	case r.renderChan <- RenderVec:
	default:
		return errors.New("render channel full, render failed")
	}
	r.t1 = r.t1.Add(time.Since(time1))
	if r.dbgnm%100 == 99 {
		fmt.Printf("Renderer T1: %d, T2: %d, T3: %d, T4: %d\n", r.t1.UnixMilli()/100, r.t2.UnixMilli()/100, r.t3.UnixMilli()/100, r.t4.UnixMilli()/100)
		r.t1 = time.UnixMilli(0)
		r.t2 = time.UnixMilli(0)
		r.t3 = time.UnixMilli(0)
		r.t4 = time.UnixMilli(0)
	}
	r.dbgnm++
	return nil

}

func calculateVerticesWorker(dbg int, num int, Renderables []*Renderable, RenderVec []float32, wait *sync.WaitGroup) {
	//TODO:: This should connect to renderer and submit to it directly.
	//No need to be calculating vertices for an already updated frame

	//fmt.Printf("Worker %d started\n", dbg)
	defer wait.Done()
	for i := 0; i < num; i++ {
		//fmt.Printf("%d Out of %d\n", dbg*num+i, num*4)
		calculateVertices(Renderables[i], RenderVec[i*28:i*28+28])
	}
	//fmt.Println("Done")
}

//These never actually
func calculateVertices(renderable *Renderable, vector []float32) {

	for i := 0; i < 4; i++ {
		vert := linmath.EmptyVertice()
		vert.SetColor(renderable.Color)
		vert.SetMap(uint32(renderable.TexM))
		vert.SetTexX(renderable.TexX + renderable.TexW*float32((i)%2))
		vert.SetTexY(renderable.TexY + renderable.TexH*float32(int32((i)/2)%2))
		//(0,0),(1,0),(0,1)(1,1)
		//(-1,-1),(1,-1),(-1,1),(1,1)
		vert.SetX(float32(renderable.X + renderable.verts[i*2]))
		vert.SetY(float32(renderable.Y + renderable.verts[i*2+1]))
		vert.SetZ(float32(renderable.Z))
		//log.WithFields(log.Fields{"Vertnum": i, "X": vert[i].GetX(), "Y": vert[i].GetY(), "Z": vert[i].GetZ(), "TexX": vert[i].GetTexX(), "TexY": vert[i].GetTexY(), "Color": vert[i].GetColor()}).Trace()
		floats := vert.ToFloats()
		copy(vector[i*7:7*i+7], floats[:])
	}
}
