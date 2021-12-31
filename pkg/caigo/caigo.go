package caigo

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// N_ELEMENT_BITS_ECDSA = math.floor(math.log(FIELD_PRIME, 2))
// assert N_ELEMENT_BITS_ECDSA == 251

// N_ELEMENT_BITS_HASH = FIELD_PRIME.bit_length()
// assert N_ELEMENT_BITS_HASH == 252

func Verify(msgHash, r, s *big.Int, pub ecdsa.PublicKey, sc StarkCurve) bool {
	fmt.Println("M R S PUB SC: ", msgHash, r, s, pub)

	w := sc.InvModCurveSize(s)

	if s.Cmp(big.NewInt(0)) != 1 || s.Cmp(sc.N) != -1 {
		return false
	}
	if r.Cmp(big.NewInt(0)) != 1 || r.BitLen() > 502 {
		return false
	}
	if w.Cmp(big.NewInt(0)) != 1 || w.BitLen() > 502 {
		return false
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.BitLen() > 502 {
		return false
	}
	if !sc.IsOnCurve(pub.X, pub.Y) {
		return false
	}

	rSig := new(big.Int)
	rSig = rSig.Set(r)
	
	zGx, zGy, err := sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if err != nil {
		return false
	}

	rQx, rQy, err := sc.MimicEcMultAir(r, pub.X, pub.Y, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}

	inX, inY := sc.Add(zGx, zGy, rQx, rQy)
	wBx, wBy, err := sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}

	outX, _ := sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if rSig.Cmp(outX) == 0 {
		return true
	}
	return false
}

func XToPubKey(x string) (ecdsa.PublicKey) {
	crv := SC()
	xin := HexToBN(x)

	yout := crv.GetYCoordinate(xin)

	return ecdsa.PublicKey{
		Curve: crv,
		X: xin,
		Y: yout,
	}
}

func StrToBig(str string) (*big.Int) {
	b, _ := new(big.Int).SetString(str, 10)
	
	return b
}

func HexToBN(hexString string) (n *big.Int) {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n = new(big.Int)
	n.SetString(numStr, 16)
	return n
}

func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

func BigToHex(in *big.Int) (string) {
	return fmt.Sprintf("0x%x", in)
}