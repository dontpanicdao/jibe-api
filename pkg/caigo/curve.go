// golang transcription for https://github.com/starkware-libs/cairo-lang/tree/master/src/starkware/crypto/starkware/crypto
// use at your own risk
package caigo

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
)

var sc StarkCurve

type StarkCurve struct {
	*elliptic.CurveParams
	EcGenX           *big.Int
	EcGenY           *big.Int
	MinusShiftPointX *big.Int
	MinusShiftPointY *big.Int
	Alpha            *big.Int
}

func SC() StarkCurve {
	InitCurve()
	return sc
}

func InitCurve() {
	sc.CurveParams = &elliptic.CurveParams{Name: "stark-curve"}
	sc.P, _ = new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)  // Field Prime ./pedersen_json
	sc.N, _ = new(big.Int).SetString("3618502788666131213697322783095070105526743751716087489154079457884512865583", 10)  // Order of base point ./pedersen_json
	sc.B, _ = new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406665", 10)  // Constant of curve equation ./pedersen_json
	sc.Gx, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // (x, _) of basepoint ./pedersen_json
	sc.Gy, _ = new(big.Int).SetString("1713931329540660377023406109199410414810705867260802078187082345529207694986", 10) // (_, y) of basepoint ./pedersen_json
	sc.EcGenX, _ = new(big.Int).SetString("874739451078007766457464989774322083649278607533249481151382481072868806602", 10)
	sc.EcGenY, _ = new(big.Int).SetString("152666792071518830868575557812948353041420400780739481342941381225525861407", 10)
	sc.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	sc.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	sc.Alpha = big.NewInt(1)
	sc.BitSize = 251
}

func (sc StarkCurve) Params() *elliptic.CurveParams {
	return sc.CurveParams
}

func (sc StarkCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	yDelta := new(big.Int)
	xDelta := new(big.Int)
	yDelta.Sub(y1, y2)
	xDelta.Sub(x1, x2)

	m := DivMod(yDelta, xDelta, sc.P)

	xm := new(big.Int)
	xm = xm.Mul(m, m)

	x = new(big.Int)
	x = x.Sub(xm, x1)
	x = x.Sub(x, x2)
	x = x.Mod(x, sc.P)

	y = new(big.Int)
	y = y.Sub(x1, x)
	y = y.Mul(m, y)
	y = y.Sub(y, y1)
	y = y.Mod(y, sc.P)

	return x, y
}

func (sc StarkCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	xin := new(big.Int)
	xin = xin.Mul(big.NewInt(3), x1)
	xin = xin.Mul(xin, x1)
	xin = xin.Add(xin, sc.Alpha)

	yin := new(big.Int)
	yin = yin.Mul(y1, big.NewInt(2))

	m := DivMod(xin, yin, sc.P)

	xout := new(big.Int)
	xout = xout.Mul(m, m)
	xmed := new(big.Int)
	xmed = xmed.Mul(big.NewInt(2), x1)
	xout = xout.Sub(xout, xmed)
	xout = xout.Mod(xout, sc.P)

	yout := new(big.Int)
	yout = yout.Sub(x1, xout)
	yout = yout.Mul(m, yout)
	yout = yout.Sub(yout, y1)
	yout = yout.Mod(yout, sc.P)

	return xout, yout
}

func (sc StarkCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	var _scalarMult func(x1, y1 *big.Int, k []byte) (x, y *big.Int)
	var _add func(x1, y1, x2, y2 *big.Int) (x, y *big.Int)

	_add = func(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
		yDelta := new(big.Int)
		xDelta := new(big.Int)
		yDelta.Sub(y1, y2)
		xDelta.Sub(x1, x2)

		m := DivMod(yDelta, xDelta, sc.P)

		xm := new(big.Int)
		xm = xm.Mul(m, m)

		x = new(big.Int)
		x = x.Sub(xm, x1)
		x = x.Sub(x, x2)
		x = x.Mod(x, sc.P)

		y = new(big.Int)
		y = y.Sub(x1, x)
		y = y.Mul(m, y)
		y = y.Sub(y, y1)
		y = y.Mod(y, sc.P)

		return x, y
	}

	_scalarMult = func(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
		if len(k) == 1 {
			return x1, y1
		}
		m := new(big.Int)
		m = m.Mod(big.NewInt(int64(k[0])), big.NewInt(2))
		if m.Cmp(big.NewInt(0)) == 0 {
			h := new(big.Int)
			h = h.Div(big.NewInt(int64(k[0])), big.NewInt(2))
			c, d := sc.Double(x1, y1)
			return _scalarMult(c, d, k[1:])
		}
		e, f := _scalarMult(x1, y1, k[1:])
		return _add(e, f, x1, y1)
	}

	x, y = _scalarMult(x1, y1, k)
	return x, y
}

func (sc StarkCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return sc.ScalarMult(sc.Gx, sc.Gy, k)
}

func (sc StarkCurve) IsOnCurve(x, y *big.Int) bool {
	left := new(big.Int)
	left = left.Mul(y, y)
	left = left.Mod(left, sc.P)

	right := new(big.Int)
	right = right.Mul(x, x)
	right = right.Mul(right, x)
	right = right.Mod(right, sc.P)

	ri := new(big.Int)
	// ALPHA = big.NewInt(1)
	ri = ri.Mul(big.NewInt(1), x)

	right = right.Add(right, ri)
	right = right.Add(right, sc.B)
	right = right.Mod(right, sc.P)

	if left.Cmp(right) == 0 {
		return true
	} else {
		return false
	}
}

func (sc StarkCurve) InvModCurveSize(x *big.Int) *big.Int {
	return DivMod(big.NewInt(1), x, sc.N)
}

func (sc StarkCurve) GetYCoordinate(starkX *big.Int) *big.Int {
	y := new(big.Int)
	y = y.Mul(starkX, starkX)
	y = y.Mul(y, starkX)
	yin := new(big.Int)
	yin = yin.Mul(sc.Alpha, starkX)

	y = y.Add(y, yin)
	y = y.Add(y, sc.B)
	y = y.Mod(y, sc.P)

	// stark library checks for quad residue which is not in this implementation
	y = y.ModSqrt(y, sc.P)
	return y
}

func (sc StarkCurve) MimicEcMultAir(m, x1, y1, x2, y2 *big.Int) (x *big.Int, y *big.Int, err error) {
	// N_ELEMENT_BITS_ECDSA = 251
	if m.Cmp(big.NewInt(0)) != 1 || m.BitLen() > 502 {
		return x, y, fmt.Errorf("too many bits %v", m.BitLen())
	}

	psx := x2
	psy := y2
	for i := 0; i < 251; i++ {
		if psx == x1 {
			return x, y, fmt.Errorf("xs are the same")
		}
		// fmt.Println("INNER CHECK: ", psx, psy)
		// fmt.Println("INNER HASH: ", m)
		// fmt.Println("")
		if m.Bit(0) == 1 {
			psx, psy = sc.Add(psx, psy, x1, y1)
		}
		x1, y1 = sc.Double(x1, y1)
		m = m.Rsh(m, 1)
	}
	if m.Cmp(big.NewInt(0)) != 0 {
		return psx, psy, fmt.Errorf("m doesn't equal zero")
	}
	return psx, psy, nil
}

func DivMod(n, m, p *big.Int) *big.Int {
	q := new(big.Int)
	gx := new(big.Int)
	gy := new(big.Int)
	q = q.GCD(gx, gy, m, p)

	r := new(big.Int)
	r = r.Mul(n, gx)
	r = r.Mod(r, p)
	return r
}
