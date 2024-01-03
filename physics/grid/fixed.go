package grid

import (
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/exp/maps"
	"math"
	"slices"
)

type Fixed[T comparable] struct {
	gridSize float64
	grid     map[int64]map[Pos][]T
	fast     map[T][]func()
}

func (g *Fixed[T]) Resize(gridSize float64) {
	g.Clear()
	g.gridSize = gridSize
}

func NewFixedGrid[T comparable](gridSize float64) *Fixed[T] {
	return &Fixed[T]{
		gridSize: gridSize,
		grid:     make(map[int64]map[Pos][]T),
		fast:     make(map[T][]func()),
	}
}

func (g *Fixed[T]) hash(p Pos) int64 {
	x, y, z := p[0], p[1], p[2]
	jHead := ((x >> 63) |
		((y >> 62) & (1 << 1)) |
		((z >> 61) & (1 << 2))) << 61
	jBody := (x ^ y ^ z) >> 3
	return jHead | jBody
}

func (g *Fixed[T]) toGridCoordinates(v mgl64.Vec3) Pos {
	return Pos{
		int64(math.Floor(v.X() / g.gridSize)),
		int64(math.Floor(v.Y() / g.gridSize)),
		int64(math.Floor(v.Z() / g.gridSize)),
	}
}

func (g *Fixed[T]) Get(v mgl64.Vec3, radius float64) []T {
	coord := g.getContainedGrid(v, radius)
	var objs []T
	for _, coord := range coord {
		hash := g.hash(coord)
		m, _ := g.grid[hash]
		objs = append(objs, m[coord]...)
	}
	return slices.Compact(objs)
}

func (g *Fixed[T]) Put(center mgl64.Vec3, scale float64, value T) {
	coord := g.getContainedGrid(center, scale)
	for _, coord := range coord {
		hash := g.hash(coord)
		f, ok := g.fast[value]
		if ok {
			for _, f := range f {
				f()
			}
		}
		m, ok := g.grid[hash]
		if !ok {
			g.grid[hash] = make(map[Pos][]T, 64)
			m = g.grid[hash]
		}
		m[coord] = append(m[coord], value)
		g.fast[value] = append(g.fast[value], func() {
			idx := slices.Index(m[coord], value)
			if idx != -1 {
				m[coord] = slices.Delete(m[coord], idx, idx)
			}
		})
	}
}

func (g *Fixed[T]) Clear() {
	maps.Clear(g.grid)
	maps.Clear(g.fast)
}

func (g *Fixed[T]) getContainedGrid(center mgl64.Vec3, radius float64) []Pos {
	centerGridCoord := g.toGridCoordinates(center)

	minX := centerGridCoord[0] - int64(math.Floor(radius/g.gridSize))
	maxX := centerGridCoord[0] + int64(math.Ceil(radius/g.gridSize))
	minY := centerGridCoord[1] - int64(math.Floor(radius/g.gridSize))
	maxY := centerGridCoord[1] + int64(math.Ceil(radius/g.gridSize))
	minZ := centerGridCoord[2] - int64(math.Floor(radius/g.gridSize))
	maxZ := centerGridCoord[2] + int64(math.Ceil(radius/g.gridSize))

	var coords []Pos
	for x := minX; x <= maxX; x++ {

		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				gridCoord := Pos{x, y, z}
				coords = append(coords, gridCoord)
			}
		}
	}

	return coords
}

func (g *Fixed[T]) GetAllGridData() map[Pos][]T {
	allData := make(map[Pos][]T)

	for _, gridMap := range g.grid {
		for coord, values := range gridMap {
			if len(values) == 0 {
				continue
			}
			allData[coord] = append(allData[coord], values...)
		}
	}

	return allData
}
