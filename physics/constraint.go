package physics

import "github.com/go-gl/mathgl/mgl64"

type Field interface {
	Interact(Object) bool
	Accelerate(Object) mgl64.Vec3
}
