package light

import "math"

type RGB struct {
	R, G, B float64
}

func WaveLengthToRGB(waveLengthNm float64) RGB {
	f := func(lambda float64, a, b float64) float64 {
		//return 256.0 * math.Exp(-((lambda)-(b+a)/2)/(b-a)/2)
		//return 256.0 * math.Exp(−((lambda−(b+a)/2)/((b−a)/2)^2))
		k := (lambda - ((b + a) / 2)) / ((b - a) / 2)
		return 256.0 * math.Exp(-k*k)
	}
	return RGB{
		R: f(waveLengthNm, 520, 630),
		G: f(waveLengthNm, 500, 590),
		B: f(waveLengthNm, 410, 480),
	}
}
