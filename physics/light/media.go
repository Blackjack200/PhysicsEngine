package light

import (
	"PhysicsEngine/physics/cube"
	"math"
)

type Media interface {
	cube.Plane
	RefractionFactor() float64
}

type PlaneMedia struct {
	cube.Plane
	RF float64
}

func (p *PlaneMedia) RefractionFactor() float64 {
	return p.RF
}

const VacuumRefractionFactor = 1.0

// CalcLight from vacuum to c
func CalcLight(in *Photon, c Media) (reflection, refraction *Photon) {
	reflectionPower := 0.0

	normal := c.Normal()
	projNormal := normal.Mul(in.Direction.Dot(normal))
	planeProjectVector := in.Direction.Sub(projNormal)

	reflectionVector := planeProjectVector.Add(projNormal.Mul(-1.0))

	oCos := in.Direction.Dot(normal) / in.Direction.Len()
	oRad := math.Acos(oCos)
	oSin := math.Sin(oRad)

	outRF := c.RefractionFactor()
	fq := VacuumRefractionFactor / outRF
	fp2 := fq * oSin
	fp3 := 1 - (fp2 * fp2)
	iCos := math.Sqrt(fp3)

	if fp3 <= 0.0 {
		//total internal reflection
		reflectionPower = 1.0
	} else {
		rSq1 := (VacuumRefractionFactor * iCos) - (outRF * oCos)
		rSq2 := (VacuumRefractionFactor * iCos) + (outRF * oCos)
		rS := (rSq1 / rSq2) * (rSq1 / rSq2)

		rPq1 := (VacuumRefractionFactor * oCos) - (outRF * iCos)
		rPq2 := (VacuumRefractionFactor * oCos) + (outRF * iCos)
		rP := (rPq1 / rPq2) * (rPq1 / rPq2)
		reflectionPower = (rS + rP) / 2.0
		if reflectionPower > 1 {
			reflectionPower = 1
		}
		refractionDirection := planeProjectVector.Mul(1.0 / fq).Add(normal.Mul(-iCos))
		refraction = &Photon{
			WaveLength:      in.WaveLength,
			RenderingWeight: in.RenderingWeight * (1.0 - reflectionPower),
			Direction:       refractionDirection.Normalize(),
		}
	}

	reflection = &Photon{
		WaveLength:      in.WaveLength,
		RenderingWeight: in.RenderingWeight * reflectionPower,
		Direction:       reflectionVector.Normalize(),
	}

	return
}

// CalcLight from vacuum to c
func CalcLightj(in *Photon, c Media) (reflection, refraction *Photon) {
	reflectionPower := 0.0

	normal := c.Normal()
	projNormal := normal.Mul(in.Direction.Dot(normal))
	planeProjectVector := in.Direction.Sub(projNormal)

	reflectionVector := planeProjectVector.Add(projNormal.Mul(-1.0))

	inputCos := in.Direction.Dot(normal) / in.Direction.Len()
	inputRad := math.Acos(inputCos)
	inputSin := math.Sin(inputRad)

	outRF := c.RefractionFactor()
	fq := VacuumRefractionFactor / outRF
	fp2 := fq * inputSin
	fp3 := 1 - (fp2 * fp2)
	if fp3 <= 0.0 {
		//total internal reflection
		reflectionPower = 1.0
	} else {
		outputCos := math.Sqrt(fp3)

		rSq1 := (VacuumRefractionFactor * inputCos) - (outRF * outputCos)
		rSq2 := (VacuumRefractionFactor * inputCos) + (outRF * outputCos)
		rS := (rSq1 / rSq2) * (rSq1 / rSq2)

		rPq1 := (VacuumRefractionFactor * outputCos) - (outRF * inputCos)
		rPq2 := (VacuumRefractionFactor * outputCos) + (outRF * inputCos)
		rP := (rPq1 / rPq2) * (rPq1 / rPq2)
		reflectionPower = (rS + rP) / 2.0
		if reflectionPower > 1.0 {
			reflectionPower = 1.0
		}
		refractionDirection := planeProjectVector.Mul(1.0 / fq).Add(normal.Mul(-outputCos))
		refraction = &Photon{
			WaveLength:      in.WaveLength,
			RenderingWeight: in.RenderingWeight * (1.0 - reflectionPower),
			Direction:       refractionDirection.Normalize(),
		}
	}

	reflection = &Photon{
		WaveLength:      in.WaveLength,
		RenderingWeight: in.RenderingWeight * reflectionPower,
		Direction:       reflectionVector.Normalize(),
	}

	return
}
