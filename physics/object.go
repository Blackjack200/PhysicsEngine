package physics

import (
	"PhysicsEngine/physics/cube"
	"github.com/go-gl/mathgl/mgl64"
)

type Object interface {
	Location() mgl64.Vec3
	Mass() float64
}

type Movable interface {
	Object
	LastPosition() mgl64.Vec3
	NextTick()
	SetLocation(mgl64.Vec3)
	Acceleration() mgl64.Vec3
	Accelerate(mgl64.Vec3)
}

type Collided interface {
	Object
	Box() *cube.CollisionBox
}

type Charged interface {
	Object
	Charge() float64
}

type MoveCollided interface {
	Movable
	Collided
}
