package realworld

import (
	"PhysicsEngine/physics/cube"
	"github.com/go-gl/mathgl/mgl64"
)

type MassPoint struct {
	location     mgl64.Vec3
	lastLocation mgl64.Vec3
	acceleration mgl64.Vec3
	mass         float64
	box          *cube.CollisionBox
	charge       float64
}

func NewMassPoint(location mgl64.Vec3, mass float64, box *cube.CollisionBox, charge float64) *MassPoint {
	return &MassPoint{
		location:     location,
		lastLocation: location,
		mass:         mass,
		box:          box,
		charge:       charge,
		acceleration: mgl64.Vec3{},
	}
}

func (p *MassPoint) SetVelocity(vel mgl64.Vec3, dt float64) {
	p.lastLocation = p.location.Sub(vel.Mul(dt))
}

func (p *MassPoint) NextTick() {
	p.lastLocation = p.location
}

func (p *MassPoint) Acceleration() mgl64.Vec3 {
	return p.acceleration
}

func (p *MassPoint) Accelerate(a mgl64.Vec3) {
	p.acceleration = p.acceleration.Add(a)
	p.lastLocation = p.lastLocation.Sub(a.Mul(0.5 * 1 / 400))
}

func (p *MassPoint) LastPosition() mgl64.Vec3 {
	return p.lastLocation
}

func (p *MassPoint) Location() mgl64.Vec3 {
	return p.location
}

func (p *MassPoint) SetLocation(vec3 mgl64.Vec3) {
	p.location = vec3
}

func (p *MassPoint) Mass() float64 {
	return p.mass
}

func (p *MassPoint) Charge() float64 {
	return p.charge
}

func (p *MassPoint) Box() *cube.CollisionBox {
	return p.box
}
