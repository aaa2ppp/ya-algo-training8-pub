package main

import (
	"bufio"
	"os"
	"slices"
	"strconv"
)

func main() {
	m := make([]int32, 1e8+1)

	for i := 0; i <= 10_000; i++ {
		a := i * i
		for j := i + 1; j <= 10_000; j++ {
			v := a + j*j
			if v > 1e8 {
				break
			}
			if i == 0 {
				m[v] += 4
			} else if i == j {
				m[v] += 4
			} else {
				m[v] += 8
			}
		}
	}

	m[0] = -1

	pairs := make([][2]int32, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, [2]int32{int32(k), v})
	}

	slices.SortFunc(pairs, func(a, b [2]int32) int {
		return int(b[1] - a[1])
	})

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for _, p := range pairs[:20] {
		b := strconv.AppendInt(w.AvailableBuffer(), int64(p[0]), 10)
		w.Write(b)
		w.WriteByte(' ')
		b = strconv.AppendInt(w.AvailableBuffer(), int64(p[1]), 10)
		w.Write(b)
		w.WriteByte('\n')
	}
}
