package main

import (
	"PhysicsEngine/physics/cube"
	"PhysicsEngine/physics/light"
	"github.com/go-gl/mathgl/mgl64"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	run()
}

func run() {
	lightSources := []*light.Source{
		{
			RectangularPlane: cube.RectangularPlane{
				Normal1: mgl64.Vec3{0, 1, 0},
				Normal2: mgl64.Vec3{1, 0, 0},
				Center:  mgl64.Vec3{0, 0, 7},
				Width:   2,
				Length:  2,
			},
			WaveLength: 630,
		}, {
			RectangularPlane: cube.RectangularPlane{
				Normal1: mgl64.Vec3{0, 1, 0}.Normalize(),
				Normal2: mgl64.Vec3{1, 0, 0}.Normalize(),
				Center:  mgl64.Vec3{-1, -2, 8},
				Width:   3,
				Length:  2,
			},
			WaveLength: 450,
		},
	}

	camera := &light.WonderfulCamera{
		Source: mgl64.Vec3{0, 0, -2},
		Canvas: cube.RectangularPlane{
			Normal1: mgl64.Vec3{0, 1, 0},
			Normal2: mgl64.Vec3{1, 0, 0},
			Center:  mgl64.Vec3{0, 0, 0},
			Width:   1,
			Length:  1,
		},
		PixelWidth: 0.004,
	}

	// 设置折射介质
	refractiveMedium := &light.PlaneMedia{
		Plane: &cube.RectangularPlane{
			Normal1: mgl64.Vec3{0, 1, 0},
			Normal2: mgl64.Vec3{1, 0, 0},
			Center:  mgl64.Vec3{0, 0, 1},
			Width:   1,
			Length:  2,
		},
		RF: 1.4,
	}

	w, l := camera.Grid()
	grid := camera.RayTrace([]light.Media{refractiveMedium}, lightSources, 20)

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(l)))
	for i := range grid {
		for j := range grid[i] {
			r := grid[i][j]
			clr := color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 255,
			}

			for _, p := range r.Photons {
				c := light.WaveLengthToRGB(p.WaveLength)
				weight := p.RenderingWeight
				clr.R += uint8(weight * c.R)
				clr.G += uint8(weight * c.G)
				clr.B += uint8(weight * c.B)
			}

			img.Set(i, j, clr)
		}
	}
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
