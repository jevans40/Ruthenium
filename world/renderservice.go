package world

import (
	"errors"

	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/linmath"
)

type renderService struct {
	BaseService
	renderChan chan []float32
}

func NewRenderService(renderChan chan []float32) Service {
	newRender := &renderService{}
	newRender.renderChan = renderChan
	newRender.Name = "renderer"
	newRender.SetRunFunction(newRender.RenderRun)
	newRender.AddRequiredAccessComponent(NewComponentAccess[Transform](ReadAccess))
	newRender.AddRequiredAccessComponent(NewComponentAccess[Sprite](ReadAccess))

	return newRender
}

func (r *renderService) RenderRun(EntityCreation chan EntityCreationData, EntityDeletion chan component.EntityID) error {
	TransformRead, err1 := GetReadStorage[Transform](r)
	SpriteRead, err2 := GetReadStorage[Sprite](r)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	Entities := component.Join(TransformRead, SpriteRead)

	RenderVec := make([]float32, len(Entities)*28)
	Transforms, err1 := TransformRead.GetComponentMultiple(Entities)
	Sprites, err2 := SpriteRead.GetComponentMultiple(Entities)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	for i := range Entities {
		calculateVertices(Transforms[i], Sprites[i], RenderVec[i*28:i*28+28])
	}

	select {
	case r.renderChan <- RenderVec:
	default:
		return errors.New("render channel full, render failed")
	}
	return nil

}

func calculateVertices(transform Transform, sprite Sprite, vector []float32) {

	vert := make([]linmath.Vertice, 4)
	for i := range vert {
		vert[i] = linmath.EmptyVertice()
	}

	for i := 0; i < 4; i++ {
		vert[i].SetColor(sprite.Color)
		vert[i].SetMap(uint32(sprite.M))
		vert[i].SetTexX(sprite.X + sprite.W*float32((i)%2))
		vert[i].SetTexY(sprite.Y + sprite.H*float32(int32((i)/2)%2))
		vert[i].SetX(transform.X + transform.W*float32((i)%2))
		vert[i].SetY(transform.Y + transform.H*float32(int32((i)/2)%2))
		vert[i].SetZ(transform.Z)
	}

	for i, v := range vert {
		floats := v.ToFloats()
		copy(vector[i*7:7*i+7], floats[:])
	}
}
