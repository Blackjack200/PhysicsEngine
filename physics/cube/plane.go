package cube

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Plane interface {
	Distance(point mgl64.Vec3) float64
	Intersection(startPoint, direction mgl64.Vec3) (mgl64.Vec3, bool)
	InBoundary(point mgl64.Vec3) bool
	Normal() mgl64.Vec3
}

type SimplePlane struct {
	NormalizedNormal mgl64.Vec3
	Center           mgl64.Vec3
	Radius           float64
}

func (p *SimplePlane) Normal() mgl64.Vec3 {
	return p.NormalizedNormal
}

func (p *SimplePlane) Distance(point mgl64.Vec3) float64 {
	d := p.Center.Sub(point)
	return d.Dot(p.NormalizedNormal)
}

func (p *SimplePlane) Intersection(startPoint, directionNormalized mgl64.Vec3) (mgl64.Vec3, bool) {
	if directionNormalized.Dot(p.NormalizedNormal) == 0.0 {
		return mgl64.Vec3{}, false
	}
	dist := p.Distance(startPoint)
	//cos dist / directionNormalized length
	cos := directionNormalized.Dot(p.NormalizedNormal)
	line := directionNormalized.Mul(dist * (1.0 / cos))
	intersectionPoint := startPoint.Add(line)
	if !p.InBoundary(intersectionPoint) {
		return mgl64.Vec3{}, false
	}
	return intersectionPoint, true
}

func (p *SimplePlane) InBoundary(point mgl64.Vec3) bool {
	d := point.Sub(p.Center)
	//orthogonal decomposition
	f := p.NormalizedNormal.Mul(d.Dot(p.NormalizedNormal))
	//another
	h := d.Sub(f)
	//center to point (projected)
	planeVect := f.Add(h)
	if planeVect.Len() <= p.Radius {
		return true
	}
	return false
}

type RectangularPlane struct {
	Normal1, Normal2 mgl64.Vec3
	Center           mgl64.Vec3
	Width, Length    float64
}

func (r *RectangularPlane) Normal() mgl64.Vec3 {
	return r.Normal1.Cross(r.Normal2).Normalize()
}

func (r *RectangularPlane) Distance(point mgl64.Vec3) float64 {
	v := point.Sub(r.Center)
	return v.Dot(r.Normal1.Cross(r.Normal2).Normalize())
}

func (r *RectangularPlane) Intersection(startPoint, directionNormalized mgl64.Vec3) (mgl64.Vec3, bool) {
	normal := r.Normal1.Cross(r.Normal2).Normalize()
	denom := directionNormalized.Dot(normal)

	if denom == 0.0 {
		return mgl64.Vec3{}, false
	}

	dist := r.Distance(startPoint)
	t := dist / denom
	intersectionPoint := startPoint.Add(directionNormalized.Mul(t))

	if r.InBoundary(intersectionPoint) {
		return intersectionPoint, true
	}

	return mgl64.Vec3{}, false
}

func (r *RectangularPlane) InBoundary(point mgl64.Vec3) bool {
	v := point.Sub(r.Center)
	halfWidth := r.Width / 2.0
	halfLength := r.Length / 2.0

	xCoord := v.Dot(r.Normal1)
	yCoord := v.Dot(r.Normal2)

	return mgl64.Abs(xCoord) <= halfWidth && mgl64.Abs(yCoord) <= halfLength
}
