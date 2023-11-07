package physics

import (
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"math"
)

type Computer interface {
	ComputeLocation(obj Object, fields []Field)
}

type pos [3]int64
type gridMap struct {
	gridSize float64
	grid     map[int64]map[pos][]Object
	fast     map[Object][]func()
}

func newGridMap(gridSize float64) *gridMap {
	return &gridMap{
		gridSize: gridSize,
		grid:     make(map[int64]map[pos][]Object),
		fast:     make(map[Object][]func()),
	}
}

func (g *gridMap) hash(p pos) int64 {
	x, y, z := p[0], p[1], p[2]
	jHead := ((x >> 63) |
		((y >> 62) & (1 << 1)) |
		((z >> 61) & (1 << 2))) << 61
	jBody := (x ^ y ^ z) >> 3
	return jHead | jBody
}

func (g *gridMap) toGridCoordinates(v mgl64.Vec3) pos {
	return pos{
		int64(math.Floor(v.X() / g.gridSize)),
		int64(math.Floor(v.Y() / g.gridSize)),
		int64(math.Floor(v.Z() / g.gridSize)),
	}
}

func (g *gridMap) Get(v mgl64.Vec3, scale float64) []Object {
	coord := g.circleCoordinates(v, scale)
	var objs []Object
	for _, coord := range coord {
		hash := g.hash(coord)
		m, _ := g.grid[hash]
		objs = append(objs, m[coord]...)
	}
	return slices.Compact(objs)
}

func (g *gridMap) Put(v mgl64.Vec3, scale float64, o Object) {
	coord := g.circleCoordinates(v, scale)
	for _, coord := range coord {
		hash := g.hash(coord)
		f, ok := g.fast[o]
		if ok {
			for _, f := range f {
				f()
			}
		}
		m, ok := g.grid[hash]
		if !ok {
			g.grid[hash] = make(map[pos][]Object, 64)
			m = g.grid[hash]
		}
		m[coord] = append(m[coord], o)
		g.fast[o] = append(g.fast[o], func() {
			idx := slices.Index(m[coord], o)
			if idx != -1 {
				m[coord] = slices.Delete(m[coord], idx, idx)
			}
		})
	}
}

func (g *gridMap) Clear() {
	maps.Clear(g.grid)
	maps.Clear(g.fast)
}

func (g *gridMap) circleCoordinates(v mgl64.Vec3, r float64) []pos {
	centerPos := g.toGridCoordinates(v)

	minX := centerPos[0] - int64(math.Floor(r/g.gridSize))
	maxX := centerPos[0] + int64(math.Ceil(r/g.gridSize))
	minY := centerPos[1] - int64(math.Floor(r/g.gridSize))
	maxY := centerPos[1] + int64(math.Ceil(r/g.gridSize))
	minZ := centerPos[2] - int64(math.Floor(r/g.gridSize))
	maxZ := centerPos[2] + int64(math.Ceil(r/g.gridSize))

	var coords []pos

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				gridCoord := pos{x, y, z}
				gridCenter := mgl64.Vec3{
					float64(gridCoord[0]) * g.gridSize,
					float64(gridCoord[1]) * g.gridSize,
					float64(gridCoord[2]) * g.gridSize,
				}
				distance := v.Sub(gridCenter).Len()
				if distance <= r {
					coords = append(coords, gridCoord)
				}
			}
		}
	}

	return coords
}

type RealWorldComputer struct {
	TickPerSecond uint64
	GlobalFields  []Field
	grid          *gridMap
}

func (r *RealWorldComputer) Compute(objects []Object, fields map[Object][]Field) {
	r.grid = newGridMap(0.3)
	for _, o := range objects {
		scale := 1.0
		if o, ok := o.(Collided); ok {
			scale = o.Box().Radius
		}
		r.grid.Put(o.Location(), scale, o)
	}
	for _, o := range objects {
		var f []Field
		if fields != nil {
			f, _ = fields[o]
		}
		r.compute(o, objects, f)
	}
	r.grid.Clear()
}

func (r *RealWorldComputer) compute(
	self Object,
	objects []Object,
	objectFields []Field,
) {
	//present
	secondPerTick := float64(1) / float64(r.TickPerSecond)

	if self, ok := self.(Movable); ok {
		acceleration := mgl64.Vec3{0, 0, 0}

		for _, f := range r.GlobalFields {
			if f.Interact(self) {
				acceleration = acceleration.Add(f.Accelerate(self))
			}
		}

		for _, f := range objectFields {
			if f.Interact(self) {
				acceleration = acceleration.Add(f.Accelerate(self))
			}
		}

		velocity := self.Velocity()
		oldLoc := self.Location()
		self.SetLocation(oldLoc.Add(velocity.Mul(secondPerTick)))
		deltaVelocity := acceleration.Mul(secondPerTick)

		//collision
		if self, ok := self.(MoveCollided); ok {
			r.solveCollision(self, objects)
		}

		//chain constraints
		if self, ok := self.(Chained); ok {
			if next := self.Next(); self.HasNext() {
				delta := self.Location().Sub(oldLoc)
				next.SetLocation(next.Location().Add(delta))
			}
		}

		//future
		self.SetVelocity(velocity.Add(deltaVelocity))
	}
}

func (r *RealWorldComputer) solveCollision(self MoveCollided, objects []Object) {
	sLoc := self.Location()
	sB := self.Box().Translate(sLoc)

	for _, o := range r.grid.Get(sLoc, self.Box().Radius) {
		if o == self {
			continue
		}

		if o, ok := o.(MoveCollided); ok {
			oLoc := o.Location()
			oB := o.Box().Translate(oLoc)
			if oB.Collided(sB) {
				collisionNormal := sLoc.Sub(oLoc).Normalize()

				relativeVelocity := self.Velocity().Sub(o.Velocity())

				mass1 := self.Mass()
				mass2 := o.Mass()
				impulseMolecule := 2 * relativeVelocity.Dot(collisionNormal)
				impulseDenominator := (1/mass1 + 1/mass2) * collisionNormal.LenSqr()

				impulse := impulseMolecule / impulseDenominator

				selfVel := self.Velocity().Sub(collisionNormal.Mul(impulse / mass1))
				oVel := o.Velocity().Add(collisionNormal.Mul(impulse / mass2))
				self.SetVelocity(selfVel)
				o.SetVelocity(oVel)

				overlap := sB.Radius + oB.Radius - sLoc.Sub(oLoc).Len()
				separation := collisionNormal.Mul(overlap * 0.5)

				newLoc := sLoc.Add(separation)
				oLoc = oLoc.Sub(separation)
				self.SetLocation(newLoc)
				o.SetLocation(oLoc)

				sB = self.Box().Translate(newLoc)
			}
		}
	}
}
