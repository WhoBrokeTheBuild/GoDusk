package m32

import "math"

const (
	E   = float32(2.71828182845904523536028747135266249775724709369995957496696763) // https://oeis.org/A001113
	Pi  = float32(3.14159265358979323846264338327950288419716939937510582097494459) // https://oeis.org/A000796
	Phi = float32(1.61803398874989484820458683436563811772030917980576286213544862) // https://oeis.org/A001622

	Sqrt2   = float32(1.41421356237309504880168872420969807856967187537694807317667974) // https://oeis.org/A002193
	SqrtE   = float32(1.64872127070012814684865078781416357165377610071014801157507931) // https://oeis.org/A019774
	SqrtPi  = float32(1.77245385090551602729816748334114518279754945612238712821380779) // https://oeis.org/A002161
	SqrtPhi = float32(1.27201964951406896425242246173749149171560804184009624861664038) // https://oeis.org/A139339

	Ln2    = float32(0.693147180559945309417232121458176568075500134360255254120680009) // https://oeis.org/A002162
	Log2E  = float32(1 / Ln2)
	Ln10   = float32(2.30258509299404568401799145468436420760110148862877297603332790) // https://oeis.org/A002392
	Log10E = float32(1 / Ln10)

	MaxFloat32             = math.MaxFloat32
	SmallestNonzeroFloat32 = math.SmallestNonzeroFloat32
)

// Abs = math.Abs
func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

// Acos = math.Acos
func Acos(x float32) float32 {
	return float32(math.Acos(float64(x)))
}

// Acosh = math.Acosh
func Acosh(x float32) float32 {
	return float32(math.Acosh(float64(x)))
}

// Asin = math.Asin
func Asin(x float32) float32 {
	return float32(math.Asin(float64(x)))
}

// Asinh = math.Asinh
func Asinh(x float32) float32 {
	return float32(math.Asinh(float64(x)))
}

// Atan = math.Atan
func Atan(x float32) float32 {
	return float32(math.Atan(float64(x)))
}

// Atan2 = math.Atan2
func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(x), float64(y)))
}

// Atanh = math.Atanh
func Atanh(x float32) float32 {
	return float32(math.Atanh(float64(x)))
}

// Cbrt = math.Cbrt
func Cbrt(x float32) float32 {
	return float32(math.Cbrt(float64(x)))
}

// Ceil = math.Ceil
func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

// Copysign = math.Copysign
func Copysign(x, y float32) float32 {
	return float32(math.Copysign(float64(x), float64(y)))
}

// Cos = math.Cos
func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

// Cosh = math.Cosh
func Cosh(x float32) float32 {
	return float32(math.Cosh(float64(x)))
}

// Dim = math.Dim
func Dim(x, y float32) float32 {
	return float32(math.Dim(float64(x), float64(y)))
}

// Erf = math.Erf
func Erf(x float32) float32 {
	return float32(math.Erf(float64(x)))
}

// Erfc = math.Erfc
func Erfc(x float32) float32 {
	return float32(math.Erfc(float64(x)))
}

// Erfcinv = math.Erfcinv
func Erfcinv(x float32) float32 {
	return float32(math.Erfcinv(float64(x)))
}

// Erfinv = math.Erfinv
func Erfinv(x float32) float32 {
	return float32(math.Erfinv(float64(x)))
}

// Exp = math.Exp
func Exp(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

// Exp2 = math.Exp2
func Exp2(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}

// Expm1 = math.Expm1
func Expm1(x float32) float32 {
	return float32(math.Expm1(float64(x)))
}

// Floor = math.Floor
func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// Frexp = math.Frexp
func Frexp(f float32) (float32, int) {
	frac, exp := math.Frexp(float64(f))
	return float32(frac), exp
}

// Gamma = math.Gamma
func Gamma(x float32) float32 {
	return float32(math.Gamma(float64(x)))
}

// Hypot = math.Hypot
func Hypot(p, q float32) float32 {
	return float32(math.Hypot(float64(p), float64(q)))
}

// Ilogb = math.Ilogb
func Ilogb(x float32) int {
	return math.Ilogb(float64(x))
}

// Inf = math.Inf
func Inf(sign int) float32 {
	return float32(math.Inf(sign))
}

// IsInf = math.IsInf
func IsInf(f float32, sign int) bool {
	return math.IsInf(float64(f), sign)
}

// IsNaN = math.IsNaN
func IsNaN(f float32) bool {
	return math.IsNaN(float64(f))
}

// J0 = math.J0
func J0(x float32) float32 {
	return float32(math.J0(float64(x)))
}

// J1 = math.J1
func J1(x float32) float32 {
	return float32(math.J1(float64(x)))
}

// Jn = math.Jn
func Jn(n int, x float32) float32 {
	return float32(math.Jn(n, float64(x)))
}

// Ldexp = math.Ldexp
func Ldexp(frac float32, exp int) float32 {
	return float32(math.Ldexp(float64(frac), exp))
}

// Lgamma = math.Lgamma
func Lgamma(x float32) (float32, int) {
	lgamma, sign := math.Lgamma(float64(x))
	return float32(lgamma), sign
}

// Log = math.Log
func Log(x float32) float32 {
	return float32(math.Log(float64(x)))
}

// Log10 = math.Log10
func Log10(x float32) float32 {
	return float32(math.Log10(float64(x)))
}

// Log1p = math.Log1p
func Log1p(x float32) float32 {
	return float32(math.Log1p(float64(x)))
}

// Log2 = math.Log2
func Log2(x float32) float32 {
	return float32(math.Log2(float64(x)))
}

// Logb = math.Logb
func Logb(x float32) float32 {
	return float32(math.Logb(float64(x)))
}

// Max = math.Max
func Max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

// Min = math.Min
func Min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

// Mod = math.Mod
func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

// Modf = math.Modf
func Modf(f float32) (float32, float32) {
	whole, frac := math.Modf(float64(f))
	return float32(whole), float32(frac)
}

// NaN = math.NaN
func NaN() float32 {
	return float32(math.NaN())
}

// Pow = math.Pow
func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

// Pow10 = math.Pow10
func Pow10(n int) float32 {
	return float32(math.Pow10(n))
}

// Remainder = math.Remainder
func Remainder(x, y float32) float32 {
	return float32(math.Remainder(float64(x), float64(y)))
}

// Round = math.Round
func Round(x float32) float32 {
	return float32(math.Round(float64(x)))
}

// RoundToEven = math.RoundToEven
func RoundToEven(x float32) float32 {
	return float32(math.RoundToEven(float64(x)))
}

// Signbit = math.Signbit
func Signbit(x float32) bool {
	return math.Signbit(float64(x))
}

// Sin = math.Sin
func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

// Sincos = math.Sincos
func Sincos(x float32) (float32, float32) {
	sin, cos := math.Sincos(float64(x))
	return float32(sin), float32(cos)
}

// Sinh = math.Sinh
func Sinh(x float32) float32 {
	return float32(math.Sinh(float64(x)))
}

// Sqrt = math.Sqrt
func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// Tan = math.Tan
func Tan(x float32) float32 {
	return float32(math.Tan(float64(x)))
}

// Tanh = math.Tanh
func Tanh(x float32) float32 {
	return float32(math.Tanh(float64(x)))
}

// Trunc = math.Trunc
func Trunc(x float32) float32 {
	return float32(math.Trunc(float64(x)))
}

// Y0 = math.Y0
func Y0(x float32) float32 {
	return float32(math.Y0(float64(x)))
}

// Y1 = math.Y1
func Y1(x float32) float32 {
	return float32(math.Y1(float64(x)))
}

// Yn = math.Yn
func Yn(n int, x float32) float32 {
	return float32(math.Yn(n, float64(x)))
}
