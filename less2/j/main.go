package main

import (
	"bufio"
	"cmp"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type solveFunc func(l int, offers []offer) (cost int, chunks []int)

const (
	MaxPrice  = 1000
	MaxLength = 100
	OverCost  = MaxLength*MaxPrice + 500 // 100500 :)
)

type offer struct {
	P int // price if N < R
	R int
	Q int // price if N >= R
	F int // max amount
}

func (o offer) Cost(n int) int {
	if n > o.F {
		return OverCost
	}
	if n < o.R {
		return o.P * n
	}
	return o.Q * n
}

type cell struct {
	C int // total cost
	N int // offer N
}

func logMatrix(mx [][]cell) {
	for _, r := range mx {
		log.Println(r)
	}
}

func solve(l int, offers []offer) (cost int, chunks []int) {
	n := len(offers)
	m := 0
	t := 0 // total on the market
	for i := range n {
		m = max(offers[i].R)
		t += offers[i].F
	}
	m += l

	if t < l {
		if debugEnable {
			log.Printf("too few on the market %d < %d", t, l)
		}
		return -1, nil
	}

	mx := makeMatrix[cell](n, m+1)

	minimum := OverCost
	iMin, jMin := -1, -1

	for j := 0; j <= m; j++ {
		c := offers[0].Cost(j)
		mx[0][j] = cell{N: j, C: c}
		if j >= l && c < minimum {
			minimum = c
			iMin, jMin = 0, j
		}
	}

	for i := 1; i < n; i++ {
		for j := 0; j <= m; j++ {
			mx[i][j] = cell{N: -1, C: OverCost}
			for k := 0; k <= j; k++ {
				c := mx[i-1][j-k].C + offers[i].Cost(k)
				if c < mx[i][j].C {
					mx[i][j] = cell{N: k, C: c}
				}
			}
			if j >= l && mx[i][j].C < minimum {
				minimum = mx[i][j].C
				iMin, jMin = i, j
			}
		}
	}

	if debugEnable {
		log.Println("mx")
		logMatrix(mx)
	}

	if minimum == OverCost {
		panic("it's impossible!")
	}

	ans := make([]int, n)
	for i, j := iMin, jMin; i >= 0 && j > 0; {
		offerN := mx[i][j].N
		ans[i] = offerN
		i--
		j -= offerN
	}

	return minimum, ans
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	log.SetFlags(0)
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, l, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	offers := make([]offer, 0, n)
	for i := 0; i < n; i++ {
		p, r, q, f, err := scanFourInt(sc)
		if err != nil && err != io.EOF {
			panic(err)
		}
		offers = append(offers, offer{p, r, q, f})
	}

	cost, chunks := solve(l, offers)

	writeInt(bw, cost)
	writeInts(bw, chunks)
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
