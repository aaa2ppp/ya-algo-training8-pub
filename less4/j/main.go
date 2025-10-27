package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"unsafe"
)

type solveFunc func(l, w int, cars []car) []int

// Fraction натуральная дробь, чтобы не выходить из целочисенных вычислений
type Fraction struct {
	N, M int
}

func NewFraction(n, m int) Fraction {
	if m == 0 {
		panic("NewFraction: denominator cannot be 0")
	}
	if m < 0 {
		n = -n
		m = -m
	}
	d := gcd(abs(n), m)
	return Fraction{n / d, m / d}
}

func (f Fraction) Less(f2 Fraction) bool  { return f.N*f2.M < f2.N*f.M }
func (f Fraction) Equal(f2 Fraction) bool { return f.N == f2.N && f.M == f2.M }
func (f Fraction) Positive() bool         { return f.N > 0 }

// Point в простанстве-времени
type Point struct {
	T, X, Y Fraction
}

func (f Fraction) String() string {
	return fmt.Sprintf("%d/%d", f.N, f.M)
}

type EventType int

const (
	_ EventType = iota
	// порядок важен! для сортировки
	Collision
	Finish
)

type Event struct {
	Point
	Type EventType
}

func solve(l, w int, cars []car) []int {
	n := len(cars)

	events := make([]Event, 0, n*(n+1)) // каждое событие это точка в простанстве-времени
	collision := make(map[Point][]int)  // столкновения в точке пространства-времени
	finish := make(map[Fraction][]int)  // финиш во времени
	eliminated := make([]bool, n)       // выбывшие из соревнования

	// ищем столкновения между участниками
	for i := 0; i < n; i++ {
		a := cars[i]
		for j := i + 1; j < n; j++ {
			b := cars[j]
			if a.vx == b.vx && a.vy == b.vy {
				// движутся с одной скоростью по параллельным или совпадающим траекториям
				continue
			}
			// if a.vx*b.vy == a.vy*b.vx {
			// 	// траектории паралельны или совпадают
			// 	continue
			// }
			if (a.x-b.x)*(a.vy-b.vy) != (a.y-b.y)*(a.vx-b.vx) {
				// нет общей точки в пространстве, где будут в одно и тоже время
				continue
			}

			var t Fraction
			if a.vx != b.vx {
				t = NewFraction(a.x-b.x, b.vx-a.vx)
			} else {
				t = NewFraction(a.y-b.y, b.vy-a.vy)
			}

			if !t.Positive() {
				// пересечение в прошлом
				continue
			}

			p := Point{
				T: t,
				X: NewFraction(a.x*t.M+t.N*a.vx, t.M),
				Y: NewFraction(a.y*t.M+t.N*a.vy, t.M),
			}

			events = append(events, Event{p, Collision})
			collision[p] = append(collision[p], i, j)
		}
	}

	// ищем столкновения с бортами и финиш
	for i, c := range cars {

		if c.vy != 0 {
			// траектория пересекается с бортом
			t := NewFraction(-c.y, c.vy)
			if !t.Positive() {
				t = NewFraction(w-c.y, c.vy)
			}
			p := Point{
				T: t,
				X: NewFraction(c.x*t.M+t.N*c.vx, t.M),
				Y: NewFraction(c.y*t.M+t.N*c.vy, t.M),
			}
			events = append(events, Event{p, Collision})
			collision[p] = append(collision[p], i)
		}

		if c.vx != 0 {
			t := NewFraction(l-c.x, c.vx)
			if t.Positive() {
				// движется в сторону финиша
				p := Point{
					T: t,
					X: NewFraction(c.x*t.M+t.N*c.vx, t.M),
					Y: NewFraction(c.y*t.M+t.N*c.vy, t.M),
				}
				events = append(events, Event{p, Finish})
				finish[p.T] = append(finish[p.T], i)
			}
		}
	}

	// сортируем по времени
	slices.SortFunc(events, func(a, b Event) int {
		if a.T.Less(b.T) {
			return -1
		} else if b.T.Less(a.T) {
			return 1
		} else {
			return int(a.Type - b.Type)
		}
	})

	if debugEnable {
		log.Println("events   :", events)
		log.Println("collision:", collision)
		log.Println("finish   :", finish)
	}

	// разбирием события
	var parts []int // учасники события
	for _, e := range events {

		switch e.Type {

		case Collision:
			parts = parts[:0]
			for _, c := range collision[e.Point] {
				if eliminated[c] {
					continue
				}
				parts = append(parts, c)
			}

			// если в столкновении более одного участника или это столкновение с бортом
			if len(parts) > 1 || e.Y.Equal(Fraction{0, 1}) || e.Y.Equal(Fraction{w, 1}) {
				// то все учасники выбывают из соревнования
				for _, c := range parts {
					eliminated[c] = true
				}
			}

		case Finish:
			parts = parts[:0]
			for _, c := range finish[e.T] {
				if eliminated[c] {
					continue
				}
				parts = append(parts, c+1) // to 1-indexing
			}
			if len(parts) > 0 {
				return parts
			}
		}
	}

	return nil
}

type car struct {
	x, y, vx, vy int
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	log.SetFlags(0)
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, l, w, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	cars := make([]car, 0, n)
	for range n {
		x, y, vx, vy, err := scanFourInt(sc)
		if err != nil {
			panic(err)
		}
		cars = append(cars, car{x, y, vx, vy})
	}

	ans := solve(l, w, cars)
	writeInt(bw, len(ans))
	writeInts(bw, ans)
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
