package physics

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Computer interface {
	ComputeLocation(obj Object, fields []Field)
}

type NewtonComputer struct {
	TickPerSecond    uint64
	GlobalFields     []Field
	ConstraintFields []Field
}

func (r *NewtonComputer) Compute(
	objects []Object,
	forces map[Object][]Field,
) {
	secondPerTick := float64(1) / float64(r.TickPerSecond)
	for _, o := range objects {
		//collision
		if o, ok := o.(MoveCollided); ok {
			for i := 0; i < 100; i++ {
				r.solveCollision(o, objects)
			}
		}
		if o, ok := o.(Movable); ok {
			r.constraintLocation(o, secondPerTick)
		}
	}
	for _, o := range objects {
		var f []Field
		if forces != nil {
			f, _ = forces[o]
		}
		r.compute(o, objects, f)
	}
}

func (r *NewtonComputer) compute(
	self Object,
	objects []Object,
	objectFields []Field,
) {
	//present
	secondPerTick := float64(1) / float64(r.TickPerSecond)

	if self, ok := self.(Movable); ok {
		acceleration := mgl64.Vec3{0, 0, 0}

		for _, f := range r.GlobalFields {
			if f.Interact(self) {
				acceleration = acceleration.Add(f.Accelerate(self, secondPerTick))
			}
		}

		for _, f := range objectFields {
			if f.Interact(self) {
				acceleration = acceleration.Add(f.Accelerate(self, secondPerTick))
			}
		}

		velocity := self.Velocity()
		oldLoc := self.Location()
		self.SetLocation(oldLoc.Add(velocity.Mul(secondPerTick)))

		deltaVelocity := acceleration.Mul(secondPerTick)

		//chain constraints
		if self, ok := self.(Chained); ok {
			if next := self.Next(); self.HasNext() {
				delta := self.Location().Sub(oldLoc)
				next.SetLocation(next.Location().Add(delta))
			}
		}
		//future
		self.SetVelocity(velocity.Add(deltaVelocity))
	}
}

func (r *NewtonComputer) constraintLocation(self Movable, secondPerTick float64) {
	constraintAcc := mgl64.Vec3{0, 0, 0}

	for _, f := range r.ConstraintFields {
		if f.Interact(self) {
			constraintAcc = constraintAcc.Add(f.Accelerate(self, secondPerTick))
		}
	}
	self.SetLocation(self.Location().Add(constraintAcc.Mul(secondPerTick * secondPerTick)))
}

func (r *NewtonComputer) solveCollision(self MoveCollided, objects []Object) {
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

				relativeVelocity := self.Velocity().Sub(o.Velocity())

				mass1 := self.Mass()
				mass2 := o.Mass()
				impulseMolecule := 2 * relativeVelocity.Dot(collisionNormal)
				impulseDenominator := (1/mass1 + 1/mass2) * collisionNormal.LenSqr()

				impulse := impulseMolecule / impulseDenominator

				selfVel := self.Velocity().Sub(collisionNormal.Mul(impulse / mass1))
				oVel := o.Velocity().Add(collisionNormal.Mul(impulse / mass2))
				self.SetVelocity(selfVel)
				o.SetVelocity(oVel)

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
