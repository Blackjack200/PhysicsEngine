package motion

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
)

type Field interface {
	Accelerate(physics.Object, float64) mgl64.Vec3
}

type Constraint interface {
	Constraint(physics.Movable)
}

type SimpleConstraint struct {
	ConstraintFunc func(physics.Movable)
}

func NewConstraint(constraintFunc func(physics.Movable)) *SimpleConstraint {
	return &SimpleConstraint{ConstraintFunc: constraintFunc}
}

func (s *SimpleConstraint) Constraint(o physics.Movable) {
	s.ConstraintFunc(o)
}
