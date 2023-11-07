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

type Rendable struct {
	Obj   physics.Object
	Color color.RGBA
}

func (o *Rendable) Render(imd *imdraw.IMDraw) {
	objV := mgl64.Vec2{o.Obj.Location().X(), o.Obj.Location().Y()}.Mul(10)

	imd.Color = o.Color
	imd.Push(pixel.V(objV.X(), objV.Y()))
	imd.Circle(3, 0)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Physics Simulation",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	const simulationRate = 3
	const targetFPS = 20
	const timePrecision = simulationRate * targetFPS

	var objects []*Rendable
	var objects2 []physics.Object

	objects, computer := test(objects, timePrecision)
	for _, oo := range objects {
		objects2 = append(objects2, oo.Obj)
	}

	imd := imdraw.New(nil)
	record := false
	fps := targetFPS
	frameCnter := 0
	fpsDur := time.Now()
	for !win.Closed() {
		imd.Clear()
		if win.Pressed(pixelgl.MouseButtonLeft) || win.Pressed(pixelgl.MouseButtonRight) {
			record = !record
		}

		drawStart := time.Now()
		for _, o := range objects {
			o.Render(imd)
		}
		drawDur := time.Now().Sub(drawStart)

		if !record {
			win.Clear(colornames.Black)
		}

		stimulateStart := time.Now()
		for i := simulationRate; i > 0; i-- {
			computer.Compute(objects2, nil)
		}
		stimulateDur := time.Now().Sub(stimulateStart)
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(40, 40), basicAtlas)

		fmt.Fprintf(basicTxt, "FPS=%v OBJ=%v \n", fps, len(objects2))
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
		if record && frameCnter%10 == 0 {
			c := color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0}
			for i := float64(0); i < 10; i++ {
				position := mgl64.Vec3{20, 40 - i, 0}
				velocity := mgl64.Vec3{3, 0, 0}

				object := realworld.NewMassPoint(
					position, velocity,
					1,
					&physics.CollisionBox{Radius: 0.3},
					0,
				)

				objects = append(objects, &Rendable{
					Obj:   object,
					Color: c,
				})
				objects2 = append(objects2, object)
			}
		}
		frameCnter++
		if time.Now().Sub(fpsDur) >= time.Second {
			fpsDur = time.Now()
			fps = frameCnter
			frameCnter = 0
		}
	}
}

func test(objects []*Rendable, tickPerSecond uint64) ([]*Rendable, *physics.RealWorldComputer) {
	baseObjectA := realworld.NewMassPoint(
		mgl64.Vec3{20, 60, 0},
		mgl64.Vec3{0, 0, 0},
		0.5,
		&physics.CollisionBox{Radius: 0.2},
		0,
	)

	baseObjectB := realworld.NewMassPoint(
		mgl64.Vec3{10, 60, 0},
		mgl64.Vec3{9, 0, 0},
		4,
		&physics.CollisionBox{Radius: 0.2},
		0,
	)

	objects = append(objects, &Rendable{
		Obj:   baseObjectB,
		Color: colornames.Yellow,
	})
	objects = append(objects, &Rendable{
		Obj:   baseObjectA,
		Color: colornames.Red,
	})
	/*
		for i := 0; i < 60; i++ {
			position := mgl64.Vec3{20 + rand.Float64()*20, 40 + rand.Float64()*20, 0}
			velocity := mgl64.Vec3{20 * rand.Float64() * 10, 10 * rand.Float64() * 10, 0}

			object := realworld.NewMassPoint(
				position, velocity,
				10,
				&physics.CollisionBox{Radius: 0.1},
				0,
			)

			objects = append(objects, &Rendable{
				Obj:   object,
				Color: color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0},
			})
		}
	*/ /*
		var prevNode *realworld.ChainNode
		for i := float64(0); i < 10; i++ {
			position := mgl64.Vec3{20, 40 - i, 0}
			velocity := mgl64.Vec3{0, 0, 0}

			prevNodeN := prevNode
			prevNode = realworld.NewChainNode(realworld.NewMassPoint(
				position, velocity,
				10,
				&physics.CollisionBox{Radius: 0.1},
				0,
			))
			if i < 2 {
				prevNode.SetVelocity(mgl64.Vec3{10, 0, 0})
			}
			if prevNodeN != nil {
				prevNodeN.Connect(prevNode)
			}

			objects = append(objects, &Rendable{
				Obj:   prevNode,
				Color: color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0},
			})
		}*/
	computer := &physics.RealWorldComputer{
		TickPerSecond: tickPerSecond,
		GlobalFields: []physics.Field{
			physics.NewForce(mgl64.Vec3{0, -0.98, 0}),
			realworld.RoundGround(mgl64.Vec3{40, 40}, 30),
			//realworld.AbsorbGroundY(0, 10),
		},
	}

	return objects, computer
}

func main() {
	pixelgl.Run(run)
}
