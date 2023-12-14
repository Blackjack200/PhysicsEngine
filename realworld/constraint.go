package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
)

func RoundGround(center mgl64.Vec3, radius physics.Meter) physics.Constraint {
	return physics.NewConstraint(func(obj physics.Movable) {
		direction := obj.Location().Sub(center)
		rad := 0.0
		if obj, ok := obj.(physics.Collided); ok {
			rad = obj.Box().Radius
		}
		distance := direction.Len() + rad

		if distance >= radius {
			obj.SetLocation(center.Add(direction.Normalize().Mul(radius - rad)))
		}
	})
}
