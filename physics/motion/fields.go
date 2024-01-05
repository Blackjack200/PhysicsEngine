package motion

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
)

type Force struct {
	AccelerationFunc func(physics.Object, float64) mgl64.Vec3
}

func (f *Force) Accelerate(object physics.Object, dt float64) mgl64.Vec3 {
	return f.AccelerationFunc(object, dt)
}

func NewForce(acceleration mgl64.Vec3) *Force {
	return &Force{
		AccelerationFunc: func(physics.Object, float64) mgl64.Vec3 {
			return acceleration
		},
	}
}
