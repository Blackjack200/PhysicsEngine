package main

import (
	"PhysicsEngine/physics"
	"PhysicsEngine/physics/cube"
	"PhysicsEngine/physics/motion"
	"PhysicsEngine/physics/unit"
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
	"math"
	"math/rand"
	"time"
)

// PixelPerMeter unit p/m
const PixelPerMeter = 10

func pixelToMeter(pixel float64) float64 {
	return pixel / PixelPerMeter
}

func meterToPixel(meter unit.Meter) float64 {
	return meter * PixelPerMeter
}

// 定义环境光强度
const ambientIntensity = 0.2

type Light struct {
	Direction mgl64.Vec3
	Color     color.RGBA
	Intensity float64
}

var (
	cameraPos = mgl64.Vec3{50, 33, 75}
	fov       = 110.0
)

type Renderable3D struct {
	Obj   physics.Object
	Color color.RGBA
}

func (o *Renderable3D) Render(imd *imdraw.IMDraw, win *pixelgl.Window) {
	objV := o.Obj.Location()
	imd.Color = o.Color

	projected := perspectiveProjection(objV, win.Bounds().Center(), fov, win.Bounds().Size())
	imd.Push(pixel.V(projected.X(), projected.Y()))

	imd.Color = o.Color
	distanceToCamera := objV.Sub(cameraPos).Len()

	radius := unit.Meter(1) / distanceToCamera
	if o, ok := o.Obj.(physics.Collided); ok {
		radius = o.Box().Radius / distanceToCamera
	}
	imd.Circle(meterToPixel(radius), 0)
}

func perspectiveProjection(point mgl64.Vec3, screenCenter pixel.Vec, fov float64, screenSize pixel.Vec) mgl64.Vec2 {
	relativePos := point.Sub(cameraPos)

	scale := screenSize.X / (math.Tan(fov*math.Pi/360.0) * relativePos.Z())
	projectedX := screenCenter.X + relativePos.X()*scale
	projectedY := screenCenter.Y + relativePos.Y()*scale

	return mgl64.Vec2{projectedX, projectedY}
}

func run() {
	weight := unit.Meter(100)
	height := unit.Meter(100)
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
	const secondPerTick = 1.0 / tickPerSecond
	const targetFPS = tickPerSecond

	var objects []*Renderable3D
	objects, computer := test3D(objects, tickPerSecond)

	var objectList []physics.Object
	for _, oo := range objects {
		objectList = append(objectList, oo.Obj)
	}
	win.SetColorMask(color.White)
	imd := imdraw.New(nil)
	record := false
	fps := targetFPS
	fpsFrameCounter := 1
	cameraSpeed := 5.0

	fpsDur := time.Now()
	for !win.Closed() {
		imd.Clear()
		if win.Pressed(pixelgl.KeyW) {
			cameraPos = cameraPos.Add(mgl64.Vec3{0, 0, -cameraSpeed * secondPerTick})
		}
		if win.Pressed(pixelgl.KeyS) {
			cameraPos = cameraPos.Add(mgl64.Vec3{0, 0, cameraSpeed * secondPerTick})
		}
		if win.Pressed(pixelgl.KeyA) {
			cameraPos = cameraPos.Add(mgl64.Vec3{-cameraSpeed * secondPerTick, 0, 0})
		}
		if win.Pressed(pixelgl.KeyD) {
			cameraPos = cameraPos.Add(mgl64.Vec3{cameraSpeed * secondPerTick, 0, 0})
		}
		if win.Pressed(pixelgl.KeyUp) {
			cameraPos = cameraPos.Add(mgl64.Vec3{0, cameraSpeed * secondPerTick, 0})
		}
		if win.Pressed(pixelgl.KeyDown) {
			cameraPos = cameraPos.Add(mgl64.Vec3{0, -cameraSpeed * secondPerTick, 0})
		}
		if win.Pressed(pixelgl.MouseButtonLeft) || win.Pressed(pixelgl.MouseButtonRight) {
			record = !record
		}

		generate := win.Pressed(pixelgl.KeySpace)

		drawStart := time.Now()
		for _, o := range objects {
			o.Render(imd, win)
		}
		drawDur := time.Now().Sub(drawStart)

		if !record {
			win.Clear(colornames.Black)
		}

		forces := make(map[physics.Object][]motion.Field)
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
		if generate && fpsFrameCounter%5 == 0 {
			for i := 1.0; i < 2; i++ {
				c := color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0}

				position := mgl64.Vec3{60, 40 - i, 0}

				object := realworld.NewMassPoint(
					position,
					1,
					&cube.CollisionBox{Radius: 5},
					-1*0.0001,
				)
				object.SetVelocity(mgl64.Vec3{0, 0, 20}, secondPerTick)
				objects = append(objects, &Renderable3D{
					Obj:   object,
					Color: c,
				})
				objectList = append(objectList, object)
			}
		}
		fpsFrameCounter++
		if time.Now().Sub(fpsDur) >= time.Second {
			fpsDur = time.Now()
			fps = fpsFrameCounter
			fpsFrameCounter = 0
		}
	}
}

func test3D(objects []*Renderable3D, tickPerSecond uint64) ([]*Renderable3D, *motion.Solver) {
	computer := &motion.Solver{
		TickPerSecond:    tickPerSecond,
		CollisionPerTick: 2,
		GlobalFields: []motion.Field{
			motion.NewForce(mgl64.Vec3{0, -9.8, 0}),
		},
		Constraints: []motion.Constraint{
			realworld.RoundGround(mgl64.Vec3{50, 50, 0}, 50),
		},
	}

	return objects, computer
}

func main() {
	pixelgl.Run(run)
}
