package light

import "github.com/go-gl/mathgl/mgl64"

type Photon struct {
	WaveLength float64
	//	Intensity       float64
	RenderingWeight float64
	Direction       mgl64.Vec3
}
