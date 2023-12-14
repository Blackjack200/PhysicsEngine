package physics

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Solver struct {
	TickPerSecond    uint64
	CollisionPerTick uint64
	GlobalFields     []Field
	Constraints      []Constraint
}

func (r *Solver) Compute(
	objects []Object,
	forces map[Object][]Field,
) {

	for i := uint64(1); i < r.CollisionPerTick; i++ {
		for _, o := range objects {
			if o, ok := o.(MoveCollided); ok {
				r.solveCollision(o, objects)
			}
		}
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

	return locationFuture
}

func (r *Solver) solveCollision(self MoveCollided, objects []Object) {
	sLoc := self.Location()
	sB := self.Box().Translate(sLoc)

	for _, o := range objects {
		if o == self {
			continue
		}

		if o, ok := o.(MoveCollided); ok {
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
}
