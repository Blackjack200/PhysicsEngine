package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"math/rand"
)

func ElasticGround(height float64) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object) mgl64.Vec3 {
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

func Round(center mgl64.Vec3, radius float64) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object) mgl64.Vec3 {
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

func RoundGround(center mgl64.Vec3, radius float64) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object) mgl64.Vec3 {
		direction := obj.Location().Sub(center)
		distance := direction.Len()
		if obj, ok := obj.(physics.Movable); ok {
			if distance >= radius {
				obj.SetLocation(center.Add(direction.Normalize().Mul(radius)))
				obj.SetVelocity(mgl64.Vec3{})
			}
		}
		return mgl64.Vec3{}
	}}
}

func AbsorbGroundX(absorbFactor, height float64) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object) mgl64.Vec3 {
		if obj, ok := obj.(physics.Movable); ok {
			if math.Abs(obj.Location().X()-height) <= 0.5 {
				vel := obj.Velocity().Mul(-absorbFactor)
				obj.SetVelocity(vel)
				loc := obj.Location()
				loc[0] = height
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}

func AbsorbGroundY(absorbFactor, height float64) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object) mgl64.Vec3 {
		if obj, ok := obj.(physics.Movable); ok {
			if math.Abs(obj.Location().Y()-height) <= 0.5 {
				vel := obj.Velocity().Mul(-absorbFactor)
				obj.SetVelocity(vel)
				loc := obj.Location()
				loc[1] = height
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}

func AbsorbGroundZ(absorbFactor, height float64) physics.Field {
	return &physics.Force{AccelerationFunc: func(obj physics.Object) mgl64.Vec3 {
		if obj, ok := obj.(physics.Movable); ok {
			if math.Abs(obj.Location().Z()-height) < mgl64.Epsilon*200 {
				vel := obj.Velocity().Mul(-absorbFactor)
				obj.SetVelocity(vel)
				loc := obj.Location()
				loc[2] = height
				obj.SetLocation(loc)
			}
		}
		return mgl64.Vec3{}
	}}
}
