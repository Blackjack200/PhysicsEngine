package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	GravitationalConstant = 6.67430e-11
	CoulombConstant       = 8.99e9
)

func Universal(a physics.Object) physics.Field {
	return &physics.Force{
		AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
			if a == obj {
				return [3]float64{}
			}
			dist := obj.Location().Sub(a.Location())
			f := dist.Normalize().Mul(-(GravitationalConstant * a.Mass() * obj.Mass()) / dist.LenSqr())
			return f
		},
	}
}

func Electric(a physics.Charged) physics.Field {
	return &physics.Force{
		AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
			dist := obj.Location().Sub(a.Location())
			if ac, ok := obj.(physics.Charged); ok {
				f := dist.Normalize().Mul((CoulombConstant * ac.Charge() * obj.(physics.Charged).Charge()) / dist.LenSqr())
				return f
			}
			return mgl64.Vec3{}
		},
	}
}
