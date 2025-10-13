package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
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
		// Базовые тесты
		{
			"1 - базовый случай",
			args{strings.NewReader(`3
0 1 1
0.5 1.5 1.5
1 2 1
`)},
			`2`, // интервалы 0 и 2 (1 + 1)
			true,
		},
		{
			"2 - отрицательные координаты",
			args{strings.NewReader(`3
-2 -1 1
-1.5 -0.5 1.5
-1 0 1
`)},
			`2`, // интервалы 0 и 2 (1 + 1)
			true,
		},
		{
			"3 - zero intervals",
			args{strings.NewReader(`0`)},
			`0`,
			true,
		},

		// Краевые случаи
		{
			"4 - single interval",
			args{strings.NewReader(`1
0 1 2.5
`)},
			`2.5`,
			false,
		},
		{
			"5 - intervals with same start/end times",
			args{strings.NewReader(`3
0 1 1
0 2 2
1 2 1
`)},
			`2`, // берем самый тяжелый интервал (вес 2)
			false,
		},
		{
			"6 - nested intervals",
			args{strings.NewReader(`3
0 5 1
1 4 2
2 3 3
`)},
			`3`, // берем самый тяжелый вложенный интервал
			false,
		},
		{
			"7 - very small intervals (непересекающиеся)",
			args{strings.NewReader(`3
0 0.0001 1
0.0001 0.0002 1
0.0002 0.0003 1
`)},
			`3`, // все три интервала непересекающиеся
			false,
		},
		{
			"8 - large weights",
			args{strings.NewReader(`3
0 1 1000
0.5 1.5 2000
1 2 3000
`)},
			`4000`, // интервалы 0 и 2 (1000 + 3000)
			false,
		},
		{
			"9 - intervals in reverse order (непересекающиеся)",
			args{strings.NewReader(`3
3 4 1
2 3 1
1 2 1
`)},
			`3`, // все три интервала непересекающиеся
			false,
		},
		{
			"10 - many intervals with same endpoints",
			args{strings.NewReader(`5
0 1 1
0 1 1
0 1 1
0 1 1
0 1 1
`)},
			`1`, // можно взять только один интервал
			false,
		},
		{
			"11 - intervals with negative coordinates",
			args{strings.NewReader(`3
-5 -3 1
-4 -2 2
-3 -1 3
`)},
			`4`, // интервалы 0 и 2 (1 + 3)
			false,
		},
		{
			"12 - floating point precision (непересекающиеся)",
			args{strings.NewReader(`3
0.1 0.2 0.1
0.2 0.3 0.1
0.3 0.4 0.1
`)},
			`0.3`, // все три интервала непересекающиеся
			false,
		},
		{
			"13 - disjoint intervals",
			args{strings.NewReader(`3
0 1 1
2 3 1
4 5 1
`)},
			`3`, // все три интервала непересекающиеся
			false,
		},
		{
			"14 - all intervals start together",
			args{strings.NewReader(`4
0 1 1
0 2 1
0 3 1
0 4 1
`)},
			`1`, // можно взять только один интервал
			false,
		},
		{
			"15 - all intervals end together",
			args{strings.NewReader(`4
0 10 1
1 10 1
2 10 1
3 10 1
`)},
			`1`, // можно взять только один интервал
			false,
		},
		{
			"16 - one interval contains all others",
			args{strings.NewReader(`4
0 10 1
1 2 2
3 4 3
5 6 4
`)},
			`9`, // интервалы 1, 2, 3 (2 + 3 + 4)
			false,
		},
		{
			"17 - intervals with very close boundaries",
			args{strings.NewReader(`3
0 1 1
1 2 1
1 1.0000001 1
`)},
			`2`, // интервалы 0 и 1
			false,
		},
		{
			"18 - все интервалы пересекаются",
			args{strings.NewReader(`100
0 100 0.01
` + strings.Repeat("0 100 0.01\n", 99))},
			`0.01`, // можно взять только один интервал
			false,
		},

		// Сложные случаи выбора
		{
			"19 - сложный выбор между группами",
			args{strings.NewReader(`5
0 2 1
1 3 2
2 4 3
3 5 2
4 6 1
`)},
			`5`, // интервалы 0, 2, 4 (1 + 3 + 1)
			false,
		},
		{
			"20 - интервалы с разными стратегиями выбора",
			args{strings.NewReader(`4
0 3 5
1 2 3
2 4 3
3 5 5
`)},
			`10`, // // интервалы 0 и 3 (5 + 5) - НЕ пересекаются!
			false,
		},
		{
			"21 - граничный случай с точками соприкосновения",
			args{strings.NewReader(`3
0 1 2
1 2 2
2 3 2
`)},
			`6`, // все три интервала непересекающиеся
			false,
		},
		{
			"22 - один очень тяжелый интервал vs несколько легких",
			args{strings.NewReader(`4
0 2 10
1 3 1
2 4 1
3 5 1
`)},
			`11`, // интервалы 0 и 3 (10 + 1) - интервалы 2 и 3 пересекаются!
			false,
		},
		{
			"23 - жадный выбор не оптимален",
			args{strings.NewReader(`3
0 3 3
2 4 3
3 5 3
`)},
			`6`, // интервалы 0 и 2 (3 + 3)
			false,
		},
		{
			"24 - цепочка интервалов",
			args{strings.NewReader(`4
0 1 2
1 2 3
2 3 4
3 4 5
`)},
			`14`, // все интервалы (2 + 3 + 4 + 5)
			false,
		},
		{
			"25 - альтернативные пути",
			args{strings.NewReader(`4
0 2 5
1 3 4
2 4 3
3 5 2
`)},
			`8`, // интервалы 0 и 2 (5 + 3) - НЕ пересекаются: [0,2) и [2,4)
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)

			var want, got float64
			fmt.Fscan(strings.NewReader(tt.wantOut), &want)
			fmt.Fscan(out, &got)
			if math.Abs(want-got) > 1e-4 || (want != 0 && math.Abs(want-got)/want > 1e-4) {
				t.Errorf("run() = %v, want %v", got, want)
			}
		})
	}
}
