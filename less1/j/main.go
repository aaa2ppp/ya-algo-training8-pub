package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type solveFunc func(n int, t []string) []int

type List interface {
	SubList(idx, size int) *slice
	Set(idx int, val int)
	Add(val int)
	Get(idx int) int
}

type list struct {
	items []int
}

func NewList(items []int) *list {
	return &list{items}
}

func (l *list) SubList(idx, size int) *slice {
	return &slice{
		list: l,
		idx:  idx,
		size: size,
	}
}

func (l *list) Set(idx int, val int) {
	l.items[idx] = val
}

func (l *list) Add(val int) {
	l.items = append(l.items, val)
}

func (l *list) Get(idx int) int {
	return l.items[idx]
}

type slice struct {
	list *list
	idx  int
	size int
}

func (s *slice) SubList(idx, size int) *slice {
	return s.list.SubList(s.idx+idx, size)
}

func (s *slice) Set(idx int, val int) {
	s.list.Set(s.idx+idx, val)
}

func (l *slice) Add(_ int) {} // stub

func (s *slice) Get(idx int) int {
	return s.list.Get(s.idx + idx)
}

type Figon struct {
	lists map[string]List
}

func NewFigon() *Figon {
	return &Figon{
		lists: map[string]List{},
	}
}

func (f *Figon) parseNewList(s string) {
	dstName := f.parseListName(s)
	p := strings.IndexByte(s, '=')
	s = s[p+1:]
	args := f.parseArgs(s)
	f.lists[dstName] = NewList(args)
}

func (f *Figon) parseSubList(s string) {
	dstName := f.parseListName(s)
	p := strings.IndexByte(s, '=')
	s = s[p+1:]
	args := f.parseArgs(s)
	srcName := f.parseDotName(s)
	from, to := args[0], args[1]
	f.lists[dstName] = f.lists[srcName].SubList(from-1, to-from+1)
}

func (f *Figon) parseListName(s string) string {
	s = strings.TrimSpace(s)
	p := strings.IndexByte(s, ' ')
	s = s[p+1:]
	p = strings.IndexByte(s, ' ')
	return s[:p]
}

func (f *Figon) parseDotName(s string) string {
	s = strings.TrimSpace(s)
	p := strings.IndexByte(s, '.')
	return s[:p]
}

func (f *Figon) parseArgs(s string) []int {
	p := strings.IndexByte(s, '(')
	s = s[p+1:]
	p = strings.IndexByte(s, ')')

	parts := strings.Split(s[:p], ",")
	args := make([]int, 0, len(parts))
	for _, s := range parts {
		v, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		args = append(args, v)
	}
	return args
}

func (f *Figon) parseSet(s string) {
	name := f.parseDotName(s)
	args := f.parseArgs(s)
	idx, val := args[0], args[1]
	f.lists[name].Set(idx-1, val) // to 0-indexing
}

func (f *Figon) parseAdd(s string) {
	name := f.parseDotName(s)
	args := f.parseArgs(s)
	val := args[0]
	f.lists[name].Add(val)
}

func (f *Figon) parseGet(s string) int {
	name := f.parseDotName(s)
	args := f.parseArgs(s)
	idx := args[0]
	return f.lists[name].Get(idx - 1) // to 0-indexing
}

func solve(n int, text []string) []int {
	f := NewFigon()

	var ans []int
	for _, s := range text {
		s = strings.TrimSpace(strings.TrimSpace(s))
		switch {
		case strings.HasPrefix(s, "List "):
			if strings.Contains(s, " new ") {
				f.parseNewList(s)
			} else {
				f.parseSubList(s)
			}
		case strings.Contains(s, ".set("):
			f.parseSet(s)
		case strings.Contains(s, ".add("):
			f.parseAdd(s)
		case strings.Contains(s, ".get("):
			ans = append(ans, f.parseGet(s))
		default:
			panic(fmt.Errorf("unknown syntax: %q", s))
		}
	}
	return ans
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	log.SetFlags(0)
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n int
	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}

	text := make([]string, 0, n)
	for i := 0; i < n; i++ {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		text = append(text, line)
	}

	ans := solve(n, text)

	writeInts(bw, ans, writeOpts{sep: "\n", end: "\n"})
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

type writeOpts struct {
	sep   string
	begin string
	end   string
}

var defaultWriteOpts = writeOpts{
	sep: " ",
	end: "\n",
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
