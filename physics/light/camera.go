package light

import (
	"PhysicsEngine/physics/cube"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type RayTraceResult struct {
	Photons []*Photon
}

type Camera interface {
	PixelMeter() float64
	Grid() (height, width uint64)
	RayTrace(plane []Media, maxIteration uint64) [][]*RayTraceResult
}

type WonderfulCamera struct {
	Canvas     cube.RectangularPlane
	Source     mgl64.Vec3
	PixelWidth float64
}

func (c *WonderfulCamera) PixelMeter() float64 {
	return c.PixelWidth
}

func (c *WonderfulCamera) Grid() (width, length uint64) {
	return uint64(c.Canvas.Width / c.PixelMeter()), uint64(c.Canvas.Length / c.PixelMeter())
}

func (c *WonderfulCamera) RayTrace(plane []Media, lightSources []*Source, maxIteration uint64) [][]*RayTraceResult {
	w, l := c.Grid()
	grid := make([][]*RayTraceResult, w)
	for i := uint64(0); i < w; i++ {
		grid[i] = make([]*RayTraceResult, l)
		for j := uint64(0); j < l; j++ {
			//per pixel shading

			pixelCenter := c.Canvas.Center.
				Add(c.Canvas.Normal1.Mul(float64(i) * c.PixelWidth)).
				Add(c.Canvas.Normal2.Mul(float64(j) * c.PixelWidth))

			masterPhoton := &Photon{
				//currently we don't have an light source
				WaveLength:      0.0,
				RenderingWeight: 1.0,
				Direction:       pixelCenter.Sub(c.Source).Normalize(),
			}

			result := &RayTraceResult{}
			c.trace(result, pixelCenter, masterPhoton, plane, lightSources, 1, maxIteration)
			grid[i][j] = result
		}
	}
	return grid
}

func (c *WonderfulCamera) trace(
	result *RayTraceResult,
	source mgl64.Vec3, current *Photon,
	plane []Media,
	lightSources []*Source,
	depth, maxDepth uint64,
) {
	if len(plane) == 0 {
		for _, lightSource := range lightSources {
			if _, ok := lightSource.Intersection(source, current.Direction.Normalize()); ok {
				nn := &Photon{
					WaveLength:      lightSource.WaveLength,
					RenderingWeight: current.RenderingWeight,
					Direction:       current.Direction.Normalize(),
				}
				result.Photons = append(result.Photons, nn)
				return
			}
		}
		return
	}
	if depth == maxDepth {
		return
	}
	intersect := plane[0]
	lastDist := math.MaxFloat64
	for _, p := range plane {
		if _, ok := p.Intersection(source, current.Direction.Normalize()); ok {
			dist := p.Distance(source)
			if dist < lastDist {
				intersect = p
				lastDist = dist
			}
		}
	}
	if intersectionPoint, ok := intersect.Intersection(source, current.Direction.Normalize()); ok {
		refl, refr := CalcLight(current, intersect)
		fg := false
		if refl != nil && refl.RenderingWeight > 0.0 {
			fg = true
			c.trace(result, intersectionPoint, refl, plane, lightSources, depth+1, maxDepth)
		}
		if refr != nil && refr.RenderingWeight > 0.0 {
			fg = false
			c.trace(result, intersectionPoint, refr, plane, lightSources, depth+1, maxDepth)
		}
		if !fg {
		}
	}
	for _, lightSource := range lightSources {
		if _, ok := lightSource.Intersection(source, current.Direction.Normalize()); ok {
			nn := &Photon{
				WaveLength:      lightSource.WaveLength,
				RenderingWeight: current.RenderingWeight,
				Direction:       current.Direction.Normalize(),
			}
			result.Photons = append(result.Photons, nn)
			return
		}
	}
}
