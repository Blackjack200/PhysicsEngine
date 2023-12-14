package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

func NewSimpleHarmonicMotion(origin mgl64.Vec3, k physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		return obj.Location().Sub(origin).Mul(-k)
	}}
}

func NewSimpleHarmonicMotionX(origin mgl64.Vec3, k physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		return mgl64.Vec3{
			(obj.Location().X() - origin.X()) * -k,
			0,
			0,
		}
	}}
}

func NewSimpleHarmonicMotionY(origin mgl64.Vec3, k physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		return mgl64.Vec3{
			0,
			(obj.Location().Y() - origin.Y()) * -k,
			0,
		}
	}}
}

func NewSimpleHarmonicMotionZ(origin mgl64.Vec3, k physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		return mgl64.Vec3{
			0,
			0,
			(obj.Location().Z() - origin.Z()) * -k,
		}
	}}
}

func NewCyclicMotion(point mgl64.Vec3, freq, radius physics.Meter, clockwise bool) (startLoc, startVel mgl64.Vec3, f physics.Field) {
	startLoc = point.Add(mgl64.Vec3{0, radius, 0})
	startVel = mgl64.Vec3{
		2 * math.Pi * radius * freq,
		0,
		0,
	}
	if !clockwise {
		startVel = startVel.Mul(-1)
	}
	f = &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		big := 4 * math.Pi * math.Pi * radius * freq * freq
		return point.Sub(obj.Location()).Normalize().Mul(big)
	}}
	return
}

const (
	GravitationalConstant = 6.67430e-11
	CoulombConstant       = 8.99e9
)

func Universal(a physics.Object) physics.Field {
	return &physics.SpotField{
		Center: a.Location(),
		AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
			if a == obj {
				return [3]float64{}
			}
			dist := obj.Location().Sub(a.Location())
			f := dist.Normalize().Mul(-(GravitationalConstant * a.Mass() * obj.Mass()) / dist.LenSqr())
			if math.IsNaN(f[0]) {
				f[0] = 0
			}
			if math.IsNaN(f[1]) {
				f[1] = 0
			}
			if math.IsNaN(f[2]) {
				f[2] = 0
			}
			return f
		},
	}
}

func Electric(a physics.Charged) physics.Field {
	return &physics.SpotField{
		Center: a.Location(),
		AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
			dist := obj.Location().Sub(a.Location())
			if ac, ok := obj.(physics.Charged); ok {
				f := dist.Normalize().Mul((CoulombConstant * ac.Charge() * obj.(physics.Charged).Charge()) / dist.LenSqr())
				return f
			}
			return [3]float64{}
		},
	}
}
