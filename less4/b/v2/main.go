package main

import (
	"bufio"
	"io"
	"os"
	"slices"
	"strconv"
	"unsafe"
)

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n := scanInt32(sc)

	b := make([]int32, 0, n)
	t := make([]int32, 0, n)
	for range n {
		b = append(b, scanInt32(sc))
		t = append(t, scanInt32(sc))
	}

	m := scanInt32(sc)

	q := make([]int32, 0, m)
	for range m {
		q = append(q, scanInt32(sc))
	}

	for _, q := range q {
		j, _ := slices.BinarySearch(b, q)
		v := int64(t[j-1]) * int64(q)
		writeInt64(bw, v)
		bw.WriteByte('\n')
	}
}

func scanInt32(sc *bufio.Scanner) int32 {
	sc.Scan()
	v, err := strconv.ParseInt(unsafeString(sc.Bytes()), 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(v)
}

func writeInt64(bw *bufio.Writer, v int64) {
	b := strconv.AppendInt(bw.AvailableBuffer(), v, 10)
	bw.Write(b)
}

// ----------------------------------------------------------------------------

func main() {
	run(os.Stdin, os.Stdout)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
