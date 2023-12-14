package physics

import "github.com/go-gl/mathgl/mgl64"

type Force struct {
	AccelerationFunc func(Object, float64) mgl64.Vec3
}

func (f *Force) Interact(Object) bool {
	return true
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

type SpotField struct {
	Center           mgl64.Vec3
	AccelerationFunc func(Object, float64) mgl64.Vec3
}

func (f *SpotField) Interact(o Object) bool {
	return f.Center.Sub(o.Location()).LenSqr() > 0.01
}

func (f *SpotField) Accelerate(obj Object, dt float64) mgl64.Vec3 {
	return f.AccelerationFunc(obj, dt)
}

type UniformField struct {
	AccelerationFunc func(Object, float64) mgl64.Vec3
}

func (f *UniformField) Interact(o Object) bool {
	return true
}

func (f *UniformField) Accelerate(obj Object, dt float64) mgl64.Vec3 {
	return f.AccelerationFunc(obj, dt)
}
