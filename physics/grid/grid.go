package grid

import (
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/exp/maps"
	"math"
)

type Grid[T comparable] interface {
	Get(v mgl64.Vec3, radius float64) []T
	Put(center mgl64.Vec3, scale float64, value T)
	GetAllGridData() map[Pos][]T
	Clear()
}

type Pos [3]int64

func (p Pos) Add(x, y, z int64) Pos {
	return [3]int64{
		p[0] + x,
		p[1] + y,
		p[2] + z,
	}
}

type entry[T comparable] struct {
	pos Pos
	v   []T
}

func SplitByLargeGrid[T comparable](data map[Pos][]T, largeGridSize float64) []map[Pos][]T {
	splitData := make(map[Pos][]entry[T])

	for coord, v := range data {
		largeGridCoord := Pos{
			int64(math.Floor(float64(coord[0]) / largeGridSize)),
			int64(math.Floor(float64(coord[1]) / largeGridSize)),
			int64(math.Floor(float64(coord[2]) / largeGridSize)),
		}

		splitData[largeGridCoord] = append(splitData[largeGridCoord], entry[T]{
			coord,
			v,
		})
	}
	vv := maps.Values(splitData)
	nMap := make([]map[Pos][]T, len(vv))
	for i, k := range vv {
		m := make(map[Pos][]T, len(k))
		for _, k := range k {
			m[k.pos] = k.v
		}
		nMap[i] = m
	}
	return nMap
}
