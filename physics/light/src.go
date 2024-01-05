package light

import (
	"PhysicsEngine/physics/cube"
	"github.com/go-gl/mathgl/mgl64"
)

type Source struct {
	cube.RectangularPlane
	WaveLength float64
}

func NewLightSource(center, normal1, normal2 mgl64.Vec3, width, length, waveLength float64) *Source {
	return &Source{
		RectangularPlane: cube.RectangularPlane{
			Center:  center,
			Normal1: normal1,
			Normal2: normal2,
			Width:   width,
			Length:  length,
		},
		WaveLength: waveLength,
	}
}
