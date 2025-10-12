package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type solveFunc func(a, b, c float64, v [3]float64) float64

type stateFlag byte

const (
	groceriesInHands stateFlag = 1 << iota
	parcelInHands
	groceriesDelivered
	parcelDelivered
	carryingMask  = groceriesInHands | parcelInHands
	deliveredMask = groceriesDelivered | parcelDelivered
)

type state struct {
	Place byte
	F     stateFlag
}

type hop struct {
	Place byte
	Dist  float64
}

type visit struct {
	Hop int
	T   float64
}

func solve(a, b, c float64, v [3]float64) float64 {
	var (
		visited = make(map[state]visit, 256)
		queue   = make([]state, 0, 256)
	)

	rememberState := func(next state, vis visit) {
		if visNext, ok := visited[next]; !ok || vis.T < visNext.T {
			visited[next] = vis
			queue = append(queue, next)
		}
	}

	goToNextPlace := func(st state, vis visit, hops [2]hop) {
		var speed float64
		switch st.F & carryingMask {
		case 0:
			speed = v[0]
		case carryingMask:
			speed = v[2]
		default: // GroceriesInHands XOR ParcelInHands
			speed = v[1]
		}

		h := vis.Hop + 1
		for _, hop := range hops {
			t := vis.T + hop.Dist/speed
			next := st
			next.Place = hop.Place
			rememberState(next, visit{h, t})
		}
	}

	var ans float64 = 100500 // max hop = 100; min speed = 1; 3*100*1 < 100500
	visited[state{}] = visit{0, 0}
	queue = append(queue, state{})

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		vis := visited[st]

		if debugEnable {
			log.Printf("vis: %+v s: %v ", vis, st)
		}

		switch st.Place {
		case 0: // house
			if st.F&deliveredMask == deliveredMask {
				if debugEnable {
					log.Println("bingo!", vis)
				}
				ans = min(ans, vis.T)
				break
			}
			if st.F&groceriesInHands != 0 {
				next := st
				next.F &^= groceriesInHands
				next.F |= groceriesDelivered
				rememberState(next, visit{vis.Hop + 1, vis.T})
			}
			if st.F&parcelInHands != 0 {
				next := st
				next.F &^= parcelInHands
				next.F |= parcelDelivered
				rememberState(next, visit{vis.Hop + 1, vis.T})
			}
			goToNextPlace(st, vis, [2]hop{{1, a}, {2, b}})

		case 1: // supermarket
			if st.F&(groceriesDelivered|groceriesInHands) == 0 {
				next := st
				next.F |= groceriesInHands
				rememberState(next, visit{vis.Hop + 1, vis.T})
			}
			goToNextPlace(st, vis, [2]hop{{0, a}, {2, c}})

		case 2: // pick-up point
			if st.F&(parcelDelivered|parcelInHands) == 0 {
				next := st
				next.F |= parcelInHands
				rememberState(next, visit{vis.Hop + 1, vis.T})
			}
			goToNextPlace(st, vis, [2]hop{{0, b}, {1, c}})

		default:
			panic(fmt.Errorf("unknown place %v", st.Place))
		}
	}

	return ans
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	log.SetFlags(0)
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	a, b, c, err := scanThreeFloat(sc)
	if err != nil {
		panic(err)
	}
	v0, v1, v2, err := scanThreeFloat(sc)
	if err != nil {
		panic(err)
	}

	ans := solve(a, b, c, [3]float64{v0, v1, v2})
	writeFloat(bw, ans)
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout, solve)
}

// ----------------------------------------------------------------------------

type Sign interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Unsign interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

type Int interface {
	Sign | Unsign
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Int | Float
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func scanWord(sc *bufio.Scanner) (string, error) {
	if sc.Scan() {
		return sc.Text(), nil
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func _parseInt[X Int](b []byte) (X, error) {
	if ^X(0) < 0 {
		v, err := strconv.ParseInt(unsafeString(b), 0, int(unsafe.Sizeof(X(1)))<<3)
		return X(v), err
	} else {
		v, err := strconv.ParseUint(unsafeString(b), 0, int(unsafe.Sizeof(X(1)))<<3)
		return X(v), err
	}
}

func scanIntX[X Int](sc *bufio.Scanner) (X, error) {
	if !sc.Scan() {
		return 0, cmp.Or(sc.Err(), io.EOF)
	}
	return _parseInt[X](sc.Bytes())
}

func scanInts[X Int](sc *bufio.Scanner, buf []X) (_ []X, err error) {
	for n := 0; n < len(buf); n++ {
		buf[n], err = scanIntX[X](sc)
		if err != nil {
			return buf[:n], err
		}
	}
	return buf, nil
}

func scanTwoIntX[X Int](sc *bufio.Scanner) (X, X, error) {
	var buf [2]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], err
}

func scanThreeIntX[X Int](sc *bufio.Scanner) (X, X, X, error) {
	var buf [3]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], err
}

func scanFourIntX[X Int](sc *bufio.Scanner) (X, X, X, X, error) {
	var buf [4]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], err
}

func scanFiveIntX[X Int](sc *bufio.Scanner) (X, X, X, X, X, error) {
	var buf [5]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], buf[4], err
}

var (
	scanInt      = scanIntX[int]
	scanTwoInt   = scanTwoIntX[int]
	scanThreeInt = scanThreeIntX[int]
	scanFourInt  = scanFourIntX[int]
	scanFiveInt  = scanFiveIntX[int]
)

func scanFloat(sc *bufio.Scanner) (float64, error) {
	if !sc.Scan() {
		return 0, cmp.Or(sc.Err(), io.EOF)
	}
	return strconv.ParseFloat(unsafeString(sc.Bytes()), 64)
}

func scanFloats(sc *bufio.Scanner, buf []float64) (_ []float64, err error) {
	for n := 0; n < len(buf); n++ {
		buf[n], err = scanFloat(sc)
		if err != nil {
			return buf[:n], err
		}
	}
	return buf, nil
}

func scanTwoFloat(sc *bufio.Scanner) (float64, float64, error) {
	var buf [2]float64
	_, err := scanFloats(sc, buf[:])
	return buf[0], buf[1], err
}

func scanThreeFloat(sc *bufio.Scanner) (float64, float64, float64, error) {
	var buf [3]float64
	_, err := scanFloats(sc, buf[:])
	return buf[0], buf[1], buf[2], err
}

func scanFourFloat(sc *bufio.Scanner) (float64, float64, float64, float64, error) {
	var buf [4]float64
	_, err := scanFloats(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], err
}

func scanFiveFloat(sc *bufio.Scanner) (float64, float64, float64, float64, float64, error) {
	var buf [5]float64
	_, err := scanFloats(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], buf[4], err
}

type writeOpts struct {
	sep   string
	begin string
	end   string
}

var defaultWriteOpts = writeOpts{
	sep: " ",
	end: "\n",
}

func _appendInt[T Int](b []byte, v T) []byte {
	if ^T(0) < 0 {
		b = strconv.AppendInt(b, int64(v), 10)
	} else {
		b = strconv.AppendUint(b, uint64(v), 10)
	}
	return b
}

func _writeInt[X Int](bw *bufio.Writer, v X) (int, error) {
	if bw.Available() < 24 {
		bw.Flush()
	}
	return bw.Write(_appendInt(bw.AvailableBuffer(), v))
}

func writeInt[X Int](bw *bufio.Writer, v X, opts ...writeOpts) error {
	var opt writeOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultWriteOpts
	}

	bw.WriteString(opt.begin)
	_writeInt(bw, v)
	_, err := bw.WriteString(opt.end)
	return err
}

func writeInts[X Int](bw *bufio.Writer, a []X, opts ...writeOpts) error {
	var opt writeOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultWriteOpts
	}

	bw.WriteString(opt.begin)

	if len(a) != 0 {
		_writeInt(bw, a[0])
	}

	for i := 1; i < len(a); i++ {
		bw.WriteString(opt.sep)
		_writeInt(bw, a[i])
	}

	_, err := bw.WriteString(opt.end)
	return err
}

func writeFloat(bw *bufio.Writer, v float64, opts ...writeOpts) error {
	var opt writeOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultWriteOpts
	}

	b := bw.AvailableBuffer()
	b = append(b, opt.begin...)
	b = strconv.AppendFloat(b, v, 'g', -1, 64)
	b = append(b, opt.end...)
	_, err := bw.Write(b)

	return err
}

func writeFloats(bw *bufio.Writer, a []float64, opts ...writeOpts) error {
	var opt writeOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultWriteOpts
	}

	bw.WriteString(opt.begin)

	if len(a) != 0 {
		b := bw.AvailableBuffer()
		b = strconv.AppendFloat(b, a[0], 'g', -1, 64)
		bw.Write(b)
	}

	for i := 1; i < len(a); i++ {
		b := bw.AvailableBuffer()
		b = append(b, opt.sep...)
		b = strconv.AppendFloat(b, a[i], 'g', -1, 64)
		bw.Write(b)
	}

	_, err := bw.WriteString(opt.end)
	return err
}

// ----------------------------------------------------------------------------

func gcd[I Int](a, b I) I {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func gcdx[I Int](a, b I, x, y *I) I {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1, y1 I
	d := gcdx(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return d
}

func abs[N Sign | Float](a N) N {
	if a < 0 {
		return -a
	}
	return a
}

func sign[N Sign | Float](a N) N {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	mx := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		mx[i] = buf[j : j+m]
	}
	return mx
}
