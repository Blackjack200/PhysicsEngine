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

func GroundX(min, max physics.Meter) physics.Constraint {
	return physics.NewConstraint(func(obj physics.Movable) {
		radius := 0.0
		if obj, ok := obj.(physics.Collided); ok {
			radius = obj.Box().Radius
		}
		if obj.Location().X() < min+radius {
			l := obj.Location()
			l[0] = min + radius
			obj.SetLocation(l)
		}
		if obj.Location().X() > max-radius {
			l := obj.Location()
			l[0] = max - radius
			obj.SetLocation(l)
		}
	})
}

func GroundY(min, max physics.Meter) physics.Constraint {
	return physics.NewConstraint(func(obj physics.Movable) {
		radius := 0.0
		if obj, ok := obj.(physics.Collided); ok {
			radius = obj.Box().Radius
		}
		if obj.Location().Y() < min+radius {
			l := obj.Location()
			l[1] = min + radius
			obj.SetLocation(l)
		}
		if obj.Location().Y() > max-radius {
			l := obj.Location()
			l[1] = max - radius
			obj.SetLocation(l)
		}
	})
}

func GroundZ(min, max physics.Meter) physics.Constraint {
	return physics.NewConstraint(func(obj physics.Movable) {
		radius := 0.0
		if obj, ok := obj.(physics.Collided); ok {
			radius = obj.Box().Radius
		}
		if obj.Location().Z() < min+radius {
			l := obj.Location()
			l[2] = min + radius
			obj.SetLocation(l)
		}
		if obj.Location().Z() > max-radius {
			l := obj.Location()
			l[2] = max - radius
			obj.SetLocation(l)
		}
	})
}
