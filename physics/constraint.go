package physics

import "github.com/go-gl/mathgl/mgl64"

type Field interface {
	Interact(Object) bool
	Accelerate(Object, float64) mgl64.Vec3
}

type Constraint interface {
	Constraint(Movable)
}

type SimpleConstraint struct {
	ConstraintFunc func(Movable)
}

func NewConstraint(constraintFunc func(Movable)) *SimpleConstraint {
	return &SimpleConstraint{ConstraintFunc: constraintFunc}
}

func (s *SimpleConstraint) Constraint(o Movable) {
	s.ConstraintFunc(o)
}
