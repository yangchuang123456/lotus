package types

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xorcare/golden"
)

func TestPoissonFunction(t *testing.T) {
	tests := []struct {
		lambdaBase  uint64
		lambdaShift uint
	}{
		{10, 10},      // 0.0097
		{209714, 20},  // 0.19999885
		{1036915, 20}, // 0.9888792038
		{1706, 10},    // 1.6660
		{2, 0},        // 2
		{5242879, 20}, //4.9999990
		{5, 0},        // 5
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("lam-%d-%d", test.lambdaBase, test.lambdaShift), func(t *testing.T) {
			b := &bytes.Buffer{}
			b.WriteString("icdf\n")

			lam := new(big.Int).SetUint64(test.lambdaBase)
			lam = lam.Lsh(lam, precision-test.lambdaShift)
			p, icdf := newPoiss(lam)

			fmt.Println("the icdf is:\r\n", icdf)

			b.WriteString(icdf.String())
			b.WriteRune('\n')

			for i := 0; i < 15; i++ {
				b.WriteString(p.next().String())
				b.WriteRune('\n')
			}
			golden.Assert(t, []byte(b.String()))
		})
	}
}

func TestLambdaFunction(t *testing.T) {
	tests := []struct {
		power      string
		totalPower string
		target     float64
	}{
		{"10", "100", .1 * 5.},
		{"1024", "2048", 0.5 * 5.},
		{"2000000000000000", "100000000000000000", 0.02 * 5.},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%s-%s", test.power, test.totalPower), func(t *testing.T) {
			pow, ok := new(big.Int).SetString(test.power, 10)
			assert.True(t, ok)
			total, ok := new(big.Int).SetString(test.totalPower, 10)
			assert.True(t, ok)
			lam := lambda(pow, total)
			assert.Equal(t, test.target, q256ToF(lam))
			golden.Assert(t, []byte(lam.String()))
		})
	}
}

func TestExpFunction(t *testing.T) {
	const N = 256

	step := big.NewInt(5)
	step = step.Lsh(step, 256) // Q.256
	step = step.Div(step, big.NewInt(N-1))

	x := big.NewInt(0)
	b := &bytes.Buffer{}

	b.WriteString("x, y\n")
	for i := 0; i < N; i++ {
		y := expneg(x)
		fmt.Fprintf(b, "%s,%s\n", x, y)
		x = x.Add(x, step)
	}

	golden.Assert(t, b.Bytes())
}

func q256ToF(x *big.Int) float64 {
	deno := big.NewInt(1)
	deno = deno.Lsh(deno, 256)
	rat := new(big.Rat).SetFrac(x, deno)
	f, _ := rat.Float64()
	return f
}

func TestElectionLam(t *testing.T) {
	p := big.NewInt(64)
	tot := big.NewInt(128)
	lam := lambda(p, tot)
	target := 64. * 5. / 128.
	if q256ToF(lam) != target {
		t.Fatalf("wrong lambda: %f, should be: %f", q256ToF(lam), target)
	}
}

var Res int64

func BenchmarkWinCounts(b *testing.B) {
	totalPower := NewInt(100)
	power := NewInt(100)
	ep := &ElectionProof{VRFProof: nil}
	var res int64

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ep.VRFProof = []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24), byte(i >> 32)}
		j := ep.ComputeWinCount(power, totalPower)
		res += j
	}
	Res += res
}

func TestWinCounts(t *testing.T) {
	//t.SkipNow()
	totalPower := NewInt(100)
	power := NewInt(30)

	f, _ := os.Create("output.wins")
	fmt.Fprintf(f, "wins\n")
	ep := &ElectionProof{VRFProof: nil}
	for i := uint64(0); i < 1000000; i++ {
		i := i + 1000000
		ep.VRFProof = []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24), byte(i >> 32)}
		j := ep.ComputeWinCount(power, totalPower)
		fmt.Fprintf(f, "%d\n", j)
	}
}

type proportion struct {
	PowerProportion *big.Rat
	WinProportion   float64
	WinCount        uint64
}

func calcElection(t *testing.T, power, totalPower BigInt) proportion {
	fmt.Println("the power and totalPower is:", power, totalPower)
	ep := &ElectionProof{VRFProof: nil}
	count := 0
	winCount := uint64(0)
	//calcNumber := uint64(1000000)
	calcNumber := uint64(10000)
	for i := uint64(0); i < calcNumber; i++ {
		i := i + calcNumber
		ep.VRFProof = []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24), byte(i >> 32)}
		j := ep.ComputeWinCount(power, totalPower)
		if j >= 1 {
			//fmt.Println("be elected, the winCount is:",j)
			count++
		}
		winCount += uint64(j)
	}
	winCount /= calcNumber
	fmt.Println("the count is:", count)
	tmp := big.NewRat(1, 1)
	powerProportion := tmp.SetFrac(power.Int, totalPower.Int)
	//fmt.Println("the win proportion is", powerProportion, float64(count)/float64(calcNumber))
	fmt.Println("")
	return proportion{powerProportion, float64(count) / float64(calcNumber), winCount}
}

func Test_WinCounts(t *testing.T) {
	totalPower := NewInt(1000_000_00)
	power := NewInt(1)

	/*	for ; ; power = BigMul(power, NewInt(10)) {
		calcElection(t,power,totalPower)
		if power.Equals(totalPower) {
			break
		}
	}*/

	power = NewInt(1)
	totalPower = NewInt(500)

	proportions := make([]proportion, 0)
	for ; ; power = BigAdd(power, NewInt(5)) {
		proportion := calcElection(t, power, totalPower)
		if BigCmp(power, NewInt(500)) == 1 {
			break
		}
		proportions = append(proportions, proportion)
	}

	fmt.Println("the test result is:")
	for _, v := range proportions {
		f, _ := v.PowerProportion.Float64()
		fmt.Printf("%v\r\n", f)
		//fmt.Printf("%v,%v\r\n", f, v.WinProportion)
		//	fmt.Printf("%v,%v,%v\r\n", f, v.WinProportion,v.WinCount)
	}

	fmt.Printf("@@@@@@@@@@@@@@@@@\r\n")

	for _, v := range proportions {
		fmt.Printf("%v\r\n", v.WinProportion)
		//fmt.Printf("%v,%v\r\n", f, v.WinProportion)
		//	fmt.Printf("%v,%v,%v\r\n", f, v.WinProportion,v.WinCount)
	}

/*	for _, v := range proportions {
		fmt.Printf("%v\r\n", v.WinProportion)
		//fmt.Printf("%v,%v\r\n", f, v.WinProportion)
		//	fmt.Printf("%v,%v,%v\r\n", f, v.WinProportion,v.WinCount)
	}*/
}
