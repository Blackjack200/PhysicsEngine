package physics

import "github.com/go-gl/mathgl/mgl64"

type CollisionBox struct {
	Radius float64
}

type TranslatedBox struct {
	Radius float64
	Center mgl64.Vec3
}

func (b *CollisionBox) Translate(center mgl64.Vec3) *TranslatedBox {
	return &TranslatedBox{
		Radius: b.Radius,
		Center: center,
	}
}

func (b *TranslatedBox) Collided(box *TranslatedBox) bool {
	return box.Center.Sub(b.Center).Len() <= (box.Radius + b.Radius)
}
