package main

import (
	"bytes"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func Test_run_solve(t *testing.T) {
	test_run(t, solve)
}

func Test_run_slowSolve(t *testing.T) {
	test_run(t, slowSolve)
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
			args{strings.NewReader(`4
4 2 2 4
`)},
			`14
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1
10`)},
			`0`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`4
0 0 0 4
`)},
			`0
`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`4
1 0 0 4
`)},
			`0
`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`4
1 0 2 4
`)},
			`4
`,
			true,
		},
		{
			"6",
			args{strings.NewReader(`4
1 1 2 4
`)},
			`4
`,
			true,
		},
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

func generateA(rand *rand.Rand, maxN int) []int {
	n := rand.Intn(maxN) + 1
	a := make([]int, 0, n)
	for range n {
		v := rand.Intn(maxN + 1)
		a = append(a, v)
	}
	return a
}

func Fuzz_solve(f *testing.F) {
	const maxN = 10
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))
		a := generateA(rand, maxN)
		want := slowSolve(a)
		got := solve(a)
		if got != want {
			t.Logf("a: %v", a)
			t.Errorf("solve() = %d, want %d", got, want)
		}
	})
}

var benchAns int

func Benchmark_solve(b *testing.B) {
	const maxN = 100_000
	rand := rand.New(rand.NewSource(1))
	a := generateA(rand, maxN)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAns = solve(a)
	}
}
