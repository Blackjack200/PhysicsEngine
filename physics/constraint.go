package physics

import "github.com/go-gl/mathgl/mgl64"

type Field interface {
	Interact(Object) bool
	Accelerate(Object, float64) mgl64.Vec3
}
