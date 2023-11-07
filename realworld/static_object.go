package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
)

type StaticPoint struct {
	mass float64
	loc  mgl64.Vec3
	box  *physics.CollisionBox
}

func NewStaticPoint(mass float64, loc mgl64.Vec3, box *physics.CollisionBox) *StaticPoint {
	return &StaticPoint{mass: mass, loc: loc, box: box}
}

func (p *StaticPoint) Location() mgl64.Vec3 {
	return p.loc
}

func (p *StaticPoint) SetLocation(vec3 mgl64.Vec3) {
}

func (p *StaticPoint) Mass() float64 {
	return p.mass
}

func (p *StaticPoint) Box() *physics.CollisionBox {
	return p.box
}

func (p *StaticPoint) Velocity() mgl64.Vec3 {
	return mgl64.Vec3{}
}

func (p *StaticPoint) SetVelocity(vec3 mgl64.Vec3) {
}
