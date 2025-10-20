package main

import (
	"bufio"
	"flag"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

var (
	n    = flag.Int("n", 0, "length of a (by default rand.IntN(1e5)+1)")
	m    = flag.Int("m", 0, "length of b (by default rand.IntN(1e5)+1)")
	maxV = flag.Int("max", 1e4, "max value of ai or bj")
)

func main() {
	flag.Parse()
	rand := rand.New(rand.NewPCG(1, uint64(time.Now().Nanosecond())))

	if *n == 0 {
		*n = rand.IntN(1e5) + 1
	}
	if *m == 0 {
		*m = rand.IntN(1e5) + 1
	}

	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	writeInt(bw, *n)
	writeInts(bw, generator(rand, *n, *maxV)[1:])
	writeInt(bw, *m)
	writeInts(bw, generator(rand, *m, *maxV)[1:])
}

func generator(rand *rand.Rand, n int, maxV int) []int {
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rand.IntN(maxV) + 1
	}
	return a
}

type Sign interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Unsign interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

type Int interface {
	Sign | Unsign
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
