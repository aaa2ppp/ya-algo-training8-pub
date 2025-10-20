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
			args{strings.NewReader(`3 10 7
4 3
3 2
4 2
`)},
			`1.400000000000199973
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2 10 1
2 1
3 2
`)},
			`0.333333333333666715
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
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)

			want, err := strconv.ParseFloat(strings.TrimSpace(tt.wantOut), 64)
			if err != nil {
				t.Fatal(err)
			}

			got, err := strconv.ParseFloat(strings.TrimSpace(out.String()), 64)
			if err != nil {
				t.Fatal(err)
			}

			if math.Abs(want-got)/math.Max(want, 1) > 1e-6 {
				t.Errorf("run() = %v, want %v", got, want)
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
