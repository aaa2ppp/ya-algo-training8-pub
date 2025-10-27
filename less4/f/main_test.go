package main

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"
	"testing"
)

func Test_run_solve(t *testing.T) {
	test_run(t, solve)
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
			args{strings.NewReader(`3 2 0
-4 -1 1
13 6 3
-7 -6 1
1 5
`)},
			`4.333333333
5.000000000
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2 2 0
4 2 1
-11 -8 2
2 6
`)},
			`5.500000000
6.000000000
`,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
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
			want, err := toFloats(tt.wantOut)
			if err != nil {
				t.Fatalf("can't parse wantOut: %v", err)
			}

			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)

			got, err := toFloats(out.String())
			if err != nil {
				t.Fatalf("can't parse output: %v", err)
			}

			if !check(got, want) {
				t.Errorf("run() = %v, want %v", got, want)
			}
		})
	}
}

func check(got, want []float64) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if d := math.Abs(got[i] - want[i]); d > 1e-6 || d/want[i] > 1e-6 {
			return false
		}
	}
	return true
}

func toFloats(s string) ([]float64, error) {
	lines := strings.Split(s, "\n")
	if n := len(lines); lines[n-1] == "" {
		lines = lines[:n-1]
	}
	res := make([]float64, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		v, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
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
