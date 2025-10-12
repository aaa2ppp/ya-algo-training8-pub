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
			args{strings.NewReader(`3
List a = new List(2,3,5)
List b = a.subList(2,3)
b.get(1)
`)},
			`3
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5
List p = new List(2,4,8,16)
p.get(4)
List q = new List(3,9,27)
q.add(5)
q.get(4)
`)},
			`16
5
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`13
List x = new List(1,2,5,14,42)
List y = x.subList(1,4)
List z = y.subList(2,4)
y.set(1,7)
x.get(1)
z.get(1)
z.set(2,100)
x.get(3)
y.get(3)
x.add(132)
x.set(5,43)
x.get(5)
y.get(4)
`)},
			`7
2
100
100
43
14
`,
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
