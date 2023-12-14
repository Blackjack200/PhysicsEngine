package physics

import (
	"PhysicsEngine/physics/grid"
	"github.com/go-gl/mathgl/mgl64"
	"sync"
)

type CollisionBox struct {
	Radius Meter
}

type TranslatedBox struct {
	Radius Meter
	Center mgl64.Vec3
}

func (b *CollisionBox) Translate(center mgl64.Vec3) *TranslatedBox {
	return &TranslatedBox{
		Radius: b.Radius,
		Center: center,
	}
}

func (b *TranslatedBox) Collided(box *TranslatedBox) bool {
	return box.Center.Sub(b.Center).Len() <= (box.Radius + b.Radius)
}

func (r *Solver) solveCollision() {
	elem := []int64{-1, 0, 1}
	allData := r.Grid.GetAllGridData()
	wg := &sync.WaitGroup{}
	for _, grids := range grid.SplitByLargeGrid(allData, 4) {
		grids := grids
		wg.Add(1)
		go func() {
			for pos, objects := range grids {
				for _, o := range objects {
					for _, dx := range elem {
						for _, dy := range elem {
							for _, dz := range elem {
								data, ok := grids[pos.Add(dx, dy, dz)]
								if ok {
									objects = append(objects, data...)
								}
							}
						}
					}
					r.solveCollisionInternal(o, objects)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func (r *Solver) solveCollisionInternal(self MoveCollided, objects []MoveCollided) {
	sLoc := self.Location()
	sB := self.Box().Translate(sLoc)

	for _, o := range objects {
		if o == self {
			continue
		}

		oLoc := o.Location()
		oB := o.Box().Translate(oLoc)
		if oB.Collided(sB) {
			collisionNormal := sLoc.Sub(oLoc).Normalize()

			overlap := sB.Radius + oB.Radius - sLoc.Sub(oLoc).Len()
			separation := collisionNormal.Mul(overlap * 0.5)

			newLoc := sLoc.Add(separation)
			oLoc = oLoc.Sub(separation)
			self.SetLocation(newLoc)
			o.SetLocation(oLoc)
			sB = self.Box().Translate(newLoc)
		}
	}
}
