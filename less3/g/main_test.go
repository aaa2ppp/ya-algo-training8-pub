package main

import (
	"bytes"
	"io"
	"math/rand/v2"
	"strings"
	"testing"
	"time"
)

func Test_run_slowSolve(t *testing.T) {
	test_run(t, slowSolve)
}

func Test_run_solve1(t *testing.T) {
	test_run(t, solve1)
}

func test_run(t *testing.T, solve solveFunc) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`3
1 2 3
3
1 2 3
`)},
			`0
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`4
1 4 3 6
3
8 1 1
`)},
			`34
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1
1
1
2`)},
			`0`,
			true,
		},
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

func generator(rand *rand.Rand, n int, maxV int) []int {
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rand.IntN(maxV) + 1
	}
	return a
}

func Test_solve1(t *testing.T) {
	const (
		maxN = 1000
		maxV = 100
	)

	rand := rand.New(rand.NewPCG(1, 2))
	for range 10000 {
		a := generator(rand, rand.IntN(maxN)+1, maxV)
		b := generator(rand, rand.IntN(maxN)+1, maxV)
		want := slowSolve(a, b)
		got := solve1(a, b)
		if got != want {
			t.Log("a:", a[1:])
			t.Log("b:", b[1:])
			t.Fatalf("got = %v, want %v", got, want)
		}
	}
}

func Benchmark_solve1(b *testing.B) {
	const (
		maxN = int(1e6)
		maxV = int(1e4)
	)

	rand := rand.New(rand.NewPCG(1, uint64(time.Now().Nanosecond())))

	aa := generator(rand, maxN, maxV)
	bb := generator(rand, maxN, maxV)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		solve1(aa, bb)
	}
}
