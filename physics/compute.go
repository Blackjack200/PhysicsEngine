package physics

import (
	"PhysicsEngine/physics/grid"
	"github.com/go-gl/mathgl/mgl64"
)

type Solver struct {
	TickPerSecond    uint64
	CollisionPerTick uint64
	GridSize         uint64
	GlobalFields     []Field
	Constraints      []Constraint
	Grid             *grid.Fixed[MoveCollided]
}

func (r *Solver) Compute(
	objects []Object,
	forces map[Object][]Field,
) {
	if r.Grid == nil {
		r.Grid = grid.NewFixedGrid[MoveCollided](3)
	}
	sum, sam := 0.0, 0.0
	for _, o := range objects {
		if o, ok := o.(MoveCollided); ok {
			sum += o.Box().Radius
			sam += 1
		}
	}
	av := sum / sam
	r.Grid.Resize(av * 1.2)
	for _, o := range objects {
		if o, ok := o.(MoveCollided); ok {
			r.Grid.Put(o.Location(), o.Box().Radius, o)
		}
	}
	for i := uint64(1); i < r.CollisionPerTick; i++ {
		r.solveCollision()
	}

	for _, o := range objects {
		var f []Field
		if forces != nil {
			f, _ = forces[o]
		}
		r.compute(o, f)
	}
	for _, c := range r.Constraints {
		for _, o := range objects {
			if o, ok := o.(Movable); ok {
				c.Constraint(o)
			}
		}
	}
}

func (r *Solver) compute(
	self Object,
	forces []Field,
) {
	//present
	dt := float64(1) / float64(r.TickPerSecond)

	if self, ok := self.(Movable); ok {
		accelerationPresent := self.Acceleration()

		for _, f := range r.GlobalFields {
			if f.Interact(self) {
				accelerationPresent = accelerationPresent.Add(f.Accelerate(self, dt))
			}
		}

		for _, f := range forces {
			if f.Interact(self) {
				accelerationPresent = accelerationPresent.Add(f.Accelerate(self, dt))
			}
		}

		locationFuture := r.calcVerlet(self, dt, accelerationPresent)

		self.NextTick()
		//future
		self.SetLocation(locationFuture)
	}
}

func (r *Solver) calcVerlet(self Movable, dt float64, accelerationPresent mgl64.Vec3) mgl64.Vec3 {
	locationPast := self.LastPosition()
	locationPresent := self.Location()

	locationFuture := locationPresent.Mul(2).Sub(locationPast).
		Add(accelerationPresent.Mul(dt * dt))

	if locationFuture.Sub(locationPresent).Len() < 0.01 {
		locationFuture = locationPresent
	}

	return locationFuture
}
