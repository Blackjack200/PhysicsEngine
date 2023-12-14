package physics

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Computer struct {
	TickPerSecond uint64
	GlobalFields  []Field
	Constraints   []Constraint
}

func (r *Computer) Compute(
	objects []Object,
	forces map[Object][]Field,
) {

	for _, o := range objects {
		var f []Field
		if forces != nil {
			f, _ = forces[o]
		}
		if o, ok := o.(MoveCollided); ok {
			r.solveCollision(o, objects)
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

func (r *Computer) compute(
	self Object,
	forces []Field,
) {
	//present
	dt := float64(1) / float64(r.TickPerSecond)

	if self, ok := self.(Movable); ok {
		accelerationPresent := mgl64.Vec3{0, 0, 0}

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

		self.FinalizeTick()
		//future
		self.SetLocation(locationFuture)
	}
}

func (r *Computer) calcVerlet(self Movable, dt float64, accelerationPresent mgl64.Vec3) mgl64.Vec3 {
	locationPast := self.LastPosition()
	locationPresent := self.Location()

	locationFuture := locationPresent.Mul(2).Sub(locationPast).
		Add(accelerationPresent.Mul(dt * dt))

	return locationFuture
}

/*func (r *Computer) calcExplicitEuler(self Movable, dt float64, accelerationPresent mgl64.Vec3) mgl64.Vec3 {
	locationPresent := self.Location()
	velocityPresent := self.Velocity()
	locationFuture := locationPresent.Add(velocityPresent.Mul(dt))
	velocityFuture := velocityPresent.Add(accelerationPresent.Mul(dt))
	self.SetVelocity(velocityFuture)
	return locationFuture
}

func (r *Computer) calcLeapfrog(self Movable, dt float64, accelerationPresent mgl64.Vec3) mgl64.Vec3 {
	locationPresent := self.Location()
	velocityPresent := self.Velocity()

	halfVelocity := velocityPresent.Add(accelerationPresent.Mul(0.5 * dt))

	locationFuture := locationPresent.Add(halfVelocity.Mul(dt))

	velocityFuture := halfVelocity.Add(accelerationPresent.Mul(0.5 * dt))
	self.SetVelocity(velocityFuture)

	return locationFuture
}
func (r *Computer) calcRungeKutta(self Movable, dt float64, accelerationPresent mgl64.Vec3) mgl64.Vec3 {
	locationPresent := self.Location()
	velocityPresent := self.Velocity()

	// RK4 calculations
	k1 := velocityPresent.Mul(dt)
	k2 := (velocityPresent.Add(accelerationPresent.Mul(0.5 * dt))).Mul(dt)
	k3 := (velocityPresent.Add(accelerationPresent.Mul(0.5 * dt))).Mul(dt)
	k4 := (velocityPresent.Add(accelerationPresent.Mul(dt))).Mul(dt)

	// Calculate future position
	locationFuture := locationPresent.Add((k1.Add(k2).Add(k3).Add(k4)).Mul(1.0 / 6.0))

	// Update velocity using future acceleration
	velocityFuture := velocityPresent.Add(accelerationPresent.Mul(dt))
	self.SetVelocity(velocityFuture)

	return locationFuture
}*/

func (r *Computer) solveCollision(self MoveCollided, objects []Object) {
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
