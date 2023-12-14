package realworld

import (
	"PhysicsEngine/physics"
	"github.com/go-gl/mathgl/mgl64"
)

type MassPoint struct {
	location     mgl64.Vec3
	lastLocation mgl64.Vec3
	acceleration mgl64.Vec3
	mass         float64
	box          *physics.CollisionBox
	charge       float64
}

func NewMassPoint(location mgl64.Vec3, mass float64, box *physics.CollisionBox, charge float64) *MassPoint {
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

func (p *MassPoint) Box() *physics.CollisionBox {
	return p.box
}

type ChainNode struct {
	*MassPoint
	next *ChainNode
}

func NewChainNode(p *MassPoint) *ChainNode {
	return &ChainNode{
		MassPoint: p,
		next:      nil,
	}
}

func (c *ChainNode) Connect(n *ChainNode) {
	c.next = n
}

func (c *ChainNode) Next() physics.Chained {
	return c.next
}

func (c *ChainNode) HasNext() bool {
	return c.next != nil
}

func (c *ChainNode) HasPrev() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ChainNode) Prev() physics.Chained {
	//TODO implement me
	panic("implement me")
}
