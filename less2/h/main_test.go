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
			args{strings.NewReader(`1`)},
			`1`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2`)},
			`1`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`3`)},
			`1`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`4`)},
			`2`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`5`)},
			`1`,
			true,
		},
		{
			"6",
			args{strings.NewReader(`6`)},
			`1`,
			true,
		},
		{
			"7",
			args{strings.NewReader(`7`)},
			`1`,
			true,
		},
		{
			"8",
			args{strings.NewReader(`8`)},
			`2`,
			true,
		},
		{
			"9",
			args{strings.NewReader(`9`)},
			`1`,
			true,
		},
		{
			"10",
			args{strings.NewReader(`10`)},
			`1`,
			true,
		},
		{
			"11",
			args{strings.NewReader(`11`)},
			`1`,
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
