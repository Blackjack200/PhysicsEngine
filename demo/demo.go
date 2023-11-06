package physics

import (
	"PhysicsEngine/physics"
	"PhysicsEngine/realworld"
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/image/colornames"
)

func scalePosition(position mgl64.Vec3, scaleFactor float64) mgl64.Vec3 {
	return position.Mul(scaleFactor)
}

func test2(objects []*ObjectRender, timePrecision uint64, fields []physics.Field) ([]*ObjectRender, []physics.Field) {
	scaleFactor := 100.0 / 100

	// 创建地球模型
	earth := &ObjectRender{
		Obj: &physics.MovementObject{
			Object: &physics.Object{
				Location: scalePosition(mgl64.Vec3{50, 50, 0}, scaleFactor),
				Mass:     5.972e24,
			},
			RealWorldComputer: &physics.RealWorldComputer{
				TickPerSecond: timePrecision,
			},
		},
		Color: colornames.Green,
	}

	// 创建月球模型
	moon := &ObjectRender{
		Obj: &physics.MovementObject{
			Object: &physics.Object{
				Location: scalePosition(mgl64.Vec3{50 + (3.844e6), 50, 0}, scaleFactor), // 缩放月球的初始位置
				Velocity: mgl64.Vec3{0, 1022 * 1e-8, 0},                                 // 月球的线速度，根据前面的计算
				Mass:     7.342e22,                                                      // 月球的质量
			},
			RealWorldComputer: &physics.RealWorldComputer{
				TickPerSecond: timePrecision,
			},
		},
		Color: colornames.White,
	}

	// 将地球和月球模型添加到物体列表
	objects = append(objects, earth, moon)

	// 将地球和月球的引力场添加到引力场列表
	fields = append(fields, realworld.Universal(earth.Obj.Object), realworld.Universal(moon.Obj.Object))

	return objects, fields
}

func test1(objects []*ObjectRender, tickPerSecond uint64, fields []physics.Field) ([]*ObjectRender, []physics.Field) {
	baseObject1 := &physics.Object{
		Location: mgl64.Vec3{30, 50, 0},
		Velocity: mgl64.Vec3{0, 0, 0},
		Mass:     1,
	}

	objects = append(objects, &ObjectRender{
		Obj: &physics.MovementObject{
			Object: baseObject1,
			RealWorldComputer: &physics.RealWorldComputer{
				TickPerSecond: tickPerSecond,
				GlobalFields: []physics.Field{
					realworld.StartCyclicMotion(baseObject1, baseObject1.Location, 0.5, 3, true),
				},
			},
		},
		Color: colornames.Red,
	})

	baseObject2 := &physics.Object{
		Location: mgl64.Vec3{50, 50, 0},
		Velocity: mgl64.Vec3{0, 0, 0},
		Mass:     2,
	}

	objects = append(objects, &ObjectRender{
		Obj: &physics.MovementObject{
			Object: baseObject2,
			RealWorldComputer: &physics.RealWorldComputer{
				TickPerSecond: tickPerSecond,
				GlobalFields: []physics.Field{
					realworld.StartCyclicMotion(baseObject2, baseObject2.Location, 0.5, 5, false),
				},
			},
		},
		Color: colornames.Yellow,
	})

	objects = append(objects, &ObjectRender{
		Obj: &physics.MovementObject{
			Object: &physics.Object{
				Location: mgl64.Vec3{25, 35, 0},
				Velocity: mgl64.Vec3{3, 2, 0},
				Mass:     1,
			},
			RealWorldComputer: &physics.RealWorldComputer{
				TickPerSecond:  tickPerSecond,
				SubTickPerTick: tickPerSecond,
				GlobalFields: []physics.Field{
					realworld.NewSimpleHarmonicMotionY(mgl64.Vec3{45, 40, 0}, 2),
				},
			},
		},
		Color: colornames.Hotpink,
	})

	objects = append(objects, &ObjectRender{
		Obj: &physics.MovementObject{
			Object: &physics.Object{
				Location: mgl64.Vec3{15, 45, 0},
				Velocity: mgl64.Vec3{3, 2, 0},
				Mass:     1,
			},
			RealWorldComputer: &physics.RealWorldComputer{
				TickPerSecond:  tickPerSecond,
				SubTickPerTick: tickPerSecond,
				GlobalFields: []physics.Field{
					physics.NewForce(mgl64.Vec3{0, -0.98, 0}),
				},
			},
		},
		Color: colornames.Purple,
	})

	fields = append(
		fields,
		realworld.Universal(objects[0].Obj.Object),
		realworld.Universal(objects[1].Obj.Object),
		realworld.Universal(objects[2].Obj.Object),
		realworld.Universal(objects[3].Obj.Object),
	)
	return objects, fields
}
