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

type solveFunc func(n int, edges []edge) int

const MaxN = int(1e5)

type edge struct {
	a int
	b int
}

func solve(n int, edges []edge) int {
	nodes := make([][]int, n+1)
	for _, e := range edges {
		nodes[e.a] = append(nodes[e.a], e.b)
		nodes[e.b] = append(nodes[e.b], e.a)
	}

	if debugEnable {
		log.Println("nodes:", nodes)
	}

	minimum := MaxN + 1

	var dfs func(node, prev int, l int) int

	dfs = func(node, prev int, l1 int) int {
		if debugEnable {
			log.Println("->", node, l1)
		}

		minL1 := l1       // first minimum
		minL2 := MaxN + 1 // second minimum
		minL3 := MaxN + 1 // to return

		for _, next := range nodes[node] {
			if next == prev {
				continue
			}
			l2 := dfs(next, node, l1+1)
			if debugEnable {
				log.Println("<-", node, l2)
			}

			minL3 = min(minL3, l2)

			if l2 < minL1 {
				minL2, minL1 = minL1, l2
			} else if l2 < minL2 {
				minL2 = l2
			}
		}

		if minL3 == MaxN+1 {
			minL2 = 0
			minL3 = 0
		}

		minimum = min(minimum, minL1+minL2)

		if debugEnable {
			log.Println("==", node, minL1, minL2, minL1+minL2)
		}

		return minL3 + 1
	}

	for node := 1; node <= n; node++ {
		if len(nodes[node]) == 1 {
			l2 := dfs(node, 0, 0)
			minimum = min(minimum, l2)
			break
		}
	}

	return minimum
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	log.SetFlags(0)
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	if debugEnable {
		log.Println("n:", n)
	}

	edges := make([]edge, 0, n)
	for i := 0; i < n-1; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil && err != io.EOF {
			panic(err)
		}
		edges = append(edges, edge{a, b})
	}

	if debugEnable {
		log.Println("edges:", edges)
	}

	ans := solve(n, edges)
	writeInt(bw, ans)
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
