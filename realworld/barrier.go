package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"math/rand"
)

func ElasticGround(height physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		if obj, ok := obj.(physics.Movable); ok {
			if obj.Location().Y()-height < mgl64.Epsilon*20 {
				vel := obj.Velocity().Mul(-1)
				obj.SetVelocity(vel)
				loc := obj.Location()
				loc[1] = height
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}

func Round(center mgl64.Vec3, radius physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		if obj, ok := obj.(physics.Movable); ok {
			direction := obj.Location().Sub(center)
			distance := direction.Len()

			if distance >= radius {
				obj.SetLocation(center.Add(direction.Normalize().Mul(radius * 0.9999)))
				randomAngle := rand.Float64() * 2 * math.Pi
				velocity := obj.Velocity()

				rotationMatrix := mgl64.HomogRotate3D(randomAngle, mgl64.Vec3{0, 0, 1})
				newVelocity := rotationMatrix.Mul4x1(velocity.Vec4(0)).Vec3()

				obj.SetVelocity(newVelocity)
			}
		}
		return mgl64.Vec3{}
	}}
}

func RoundGround(center mgl64.Vec3, radius physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		//direction := center.Sub(obj.Location())
		direction := obj.Location().Sub(center)
		rad := 0.0
		if obj, ok := obj.(physics.Collided); ok {
			rad = obj.Box().Radius
		}
		distance := direction.Len() + rad

		if distance >= radius {
			obj.SetLocation(center.Add(direction.Normalize().Mul(radius - rad)))
		}
		return mgl64.Vec3{}
	}}
}

func AbsorbGroundX(min, max physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		if obj, ok := obj.(physics.Collided); ok {
			if obj.Location().X() >= max-obj.Box().Radius {
				loc := obj.Location()
				loc[0] = max - obj.Box().Radius
				obj.SetLocation(loc)
			}
			if obj.Location().X() <= min+obj.Box().Radius {
				loc := obj.Location()
				loc[0] = min + obj.Box().Radius
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}

func AbsorbGroundY(min, max physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		if obj, ok := obj.(physics.Collided); ok {
			if obj.Location().Y() >= max-obj.Box().Radius {
				loc := obj.Location()
				loc[1] = max - obj.Box().Radius
				obj.SetLocation(loc)
			}
			if obj.Location().Y() <= min+obj.Box().Radius {
				loc := obj.Location()
				loc[1] = min + obj.Box().Radius
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}

func AbsorbGroundZ(height physics.Meter) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object, dt float64) mgl64.Vec3 {
		if obj, ok := obj.(physics.Movable); ok {
			if obj.Location().Z() < height {
				loc := obj.Location()
				loc[2] = height
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}
