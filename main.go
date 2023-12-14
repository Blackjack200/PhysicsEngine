package main

import (
	"PhysicsEngine/physics"
	"PhysicsEngine/realworld"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math/rand"
	"time"
)

// PixelPerMeter unit p/m
const PixelPerMeter = 10

func pixelToMeter(pixel float64) float64 {
	return pixel / PixelPerMeter
}

func meterToPixel(meter physics.Meter) float64 {
	return meter * PixelPerMeter
}

type Renderable struct {
	Obj   physics.Object
	Color color.RGBA
}

func (o *Renderable) Render(imd *imdraw.IMDraw) {
	objV := mgl64.Vec2{o.Obj.Location().X(), o.Obj.Location().Y()}.Mul(PixelPerMeter)
	imd.Color = o.Color
	imd.Push(pixel.V(objV.X(), objV.Y()))
	radius := physics.Meter(1)
	if o, ok := o.Obj.(physics.Collided); ok {
		radius = o.Box().Radius
	}
	imd.Circle(meterToPixel(radius), 0)
}

func run() {
	weight := physics.Meter(100)
	height := physics.Meter(100)
	cfg := pixelgl.WindowConfig{
		Title:     "Physics Simulation",
		Bounds:    pixel.R(0, 0, meterToPixel(weight), meterToPixel(height)),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	const tickPerSecond = 20
	const targetFPS = tickPerSecond

	var objects []*Renderable
	objects, computer := test(objects, tickPerSecond)

	var objectList []physics.Object
	for _, oo := range objects {
		objectList = append(objectList, oo.Obj)
	}

	imd := imdraw.New(nil)
	record := false
	generate := false
	fps := targetFPS
	fpsFrameCounter := 1
	fpsDur := time.Now()
	for !win.Closed() {
		imd.Clear()
		if win.Pressed(pixelgl.MouseButtonLeft) || win.Pressed(pixelgl.MouseButtonRight) {
			record = !record
		}

		if win.Pressed(pixelgl.KeySpace) {
			generate = !generate
		}

		drawStart := time.Now()
		for _, o := range objects {
			o.Render(imd)
		}
		drawDur := time.Now().Sub(drawStart)

		if !record {
			win.Clear(colornames.Black)
		}

		forces := make(map[physics.Object][]physics.Field)
		/*for _, o := range objects {
			obj := o.Obj
			if obj, ok := obj.(physics.Charged); ok {
				field := realworld.Electric(obj)
				for _, oj := range objects {
					if oj.Obj != o.Obj {
						forces[oj.Obj] = append(forces[oj.Obj], field)
					}
				}
			}
			//forces[oj.Obj] = append(forces[oj.Obj], realworld.Universal(obj))
		}*/

		stimulateStart := time.Now()

		computer.Compute(objectList, forces)

		stimulateDur := time.Now().Sub(stimulateStart)
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(40, 40), basicAtlas)

		fmt.Fprintf(basicTxt, "FPS=%v OBJ=%v \n", fps, len(objects))
		sum := drawDur.Milliseconds() + stimulateDur.Milliseconds()
		expectedFrameTime := ((time.Second) / (time.Millisecond)) / targetFPS
		fmt.Fprintf(basicTxt, "DRAW=%vms SPF=%vms TOTAL=%vms OVERLAP=%vms",
			drawDur.Milliseconds(),
			stimulateDur.Milliseconds(),
			sum,
			sum-int64(expectedFrameTime),
		)

		basicTxt.Draw(win, pixel.IM)

		imd.Draw(win)
		win.Update()
		if generate && fpsFrameCounter%30 == 0 {
			c := color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0}

			position := mgl64.Vec3{60, 40, 0}
			mul := 1.0

			velocity := mgl64.Vec3{5, 0, 0}.Mul(mul)

			object := realworld.NewMassPoint(
				position, velocity,
				1,
				&physics.CollisionBox{Radius: 1},
				-1*0.0001,
			)

			objects = append(objects, &Renderable{
				Obj:   object,
				Color: c,
			})
			objectList = append(objectList, object)

		}
		fpsFrameCounter++
		if time.Now().Sub(fpsDur) >= time.Second {
			fpsDur = time.Now()
			fps = fpsFrameCounter
			fpsFrameCounter = 0
		}
	}
}

func test(objects []*Renderable, tickPerSecond uint64) ([]*Renderable, *physics.Computer) {
	computer := &physics.Computer{
		TickPerSecond: tickPerSecond,
		GlobalFields: []physics.Field{
			physics.NewForce(mgl64.Vec3{0, -9.8, 0}),
		},
		Constraints: []physics.Constraint{
			realworld.RoundGround(mgl64.Vec3{50, 50}, 40),
			//realworld.AbsorbGroundX(20, 80),
			//realworld.AbsorbGroundY(20, 80),
		},
	}

	return objects, computer
}

func main() {
	pixelgl.Run(run)
}
