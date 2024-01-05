package motion

import (
	"PhysicsEngine/physics"
	"PhysicsEngine/physics/grid"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"sync"
)

type Solver struct {
	TickPerSecond    uint64
	CollisionPerTick uint64
	GridSize         uint64
	GlobalFields     []Field
	Constraints      []Constraint
	Grid             grid.Grid[physics.MoveCollided]
}

func (r *Solver) Compute(
	objects []physics.Object,
	forces map[physics.Object][]Field,
) {
	if r.Grid == nil {
		r.Grid = grid.NewFixedGrid[physics.MoveCollided](3)
	}
	if g, ok := r.Grid.(interface{ Resize(float64) }); ok {
		sum, sam := 0.0, 0.0
		for _, o := range objects {
			if o, ok := o.(physics.MoveCollided); ok {
				sum += o.Box().Radius
				sam += 1
			}
		}
		av := sum / sam
		g.Resize(math.Ceil(av))
	}
	wg := &sync.WaitGroup{}
	for _, o := range objects {
		if o, ok := o.(physics.MoveCollided); ok {
			r.Grid.Put(o.Location(), o.Box().Radius, o)
		}
		var f []Field
		if forces != nil {
			f, _ = forces[o]
		}
		wg.Add(1)
		o := o
		go func() {
			r.compute(o, f)
			wg.Done()
		}()
	}
	wg.Wait()

	for i := uint64(1); i < r.CollisionPerTick; i++ {
		r.solveCollision()
	}

	for _, o := range objects {
		for _, c := range r.Constraints {
			if o, ok := o.(physics.Movable); ok {
				c.Constraint(o)
			}
		}
	}
}

func (r *Solver) compute(
	self physics.Object,
	forces []Field,
) {
	//present
	dt := float64(1) / float64(r.TickPerSecond)

	if self, ok := self.(physics.Movable); ok {
		accelerationPresent := self.Acceleration()

		for _, f := range r.GlobalFields {
			accelerationPresent = accelerationPresent.Add(f.Accelerate(self, dt))
		}

		for _, f := range forces {
			accelerationPresent = accelerationPresent.Add(f.Accelerate(self, dt))
		}

		locationFuture := r.calcVerlet(self, dt, accelerationPresent)

		self.NextTick()
		//future
		self.SetLocation(locationFuture)
	}
}

func (r *Solver) calcVerlet(self physics.Movable, dt float64, accelerationPresent mgl64.Vec3) mgl64.Vec3 {
	locationPast := self.LastPosition()
	locationPresent := self.Location()
	vel := locationPresent.Sub(locationPast).Mul(-1.0 / dt)
	accel := vel.Normalize().Mul(vel.LenSqr()).Mul(1 / 2)

	accelerationPresent = accelerationPresent.Add(accel)

	locationFuture := locationPresent.Mul(2).Sub(locationPast).
		Add(accelerationPresent.Mul(dt * dt))

	return locationFuture
}
