package main

import (
	"bytes"
	"io"
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
			args{strings.NewReader(`4
06:45-10:20
07:36-11:26
19:00-22:35
20:08-23:58
7
06:35-10:10
07:15-11:10
11:00-14:48
14:00-17:48
15:40-19:28
18:35-22:23
20:20-23:55
`)},
			`7`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2
10:00-12:00
15:00-17:00
2
12:30-14:30
17:30-19:30
`)},
			`1`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`2
10:10-10:11
10:10-10:11
2
10:11-10:12
10:11-10:12
`)},
			`2`,
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
