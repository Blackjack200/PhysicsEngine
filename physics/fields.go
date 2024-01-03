package physics

import "github.com/go-gl/mathgl/mgl64"

type Force struct {
	AccelerationFunc func(Object, float64) mgl64.Vec3
}

func (f *Force) Accelerate(object Object, dt float64) mgl64.Vec3 {
	return f.AccelerationFunc(object, dt)
}

func NewForce(acceleration mgl64.Vec3) *Force {
	return &Force{
		AccelerationFunc: func(Object, float64) mgl64.Vec3 {
			return acceleration
		},
	}
}
