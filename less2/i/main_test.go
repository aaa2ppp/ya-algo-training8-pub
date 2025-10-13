package main

import (
	"bytes"
	"io"
	"math/rand/v2"
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
		// Базовые тесты из условия
		{
			"1",
			args{strings.NewReader(`3 3
1 2 3
6 5 4
7 8 9
`)},
			`9
`,
			false,
		},
		{
			"2",
			args{strings.NewReader(`3 3
2 2 2
2 3 2
2 2 2
`)},
			`2
`,
			false,
		},
		{
			"3",
			args{strings.NewReader(`1 1
1`)},
			`1`,
			false,
		},
		{
			"4",
			args{strings.NewReader(`3 3
2 4 5
2 3 4
2 2 2
`)},
			`4
`,
			false,
		},

		// Крайние случаи (deepseek такой deepseek)

		// 5. Все элементы одинаковые - цепочка длины 1
		{
			"5_all_identical",
			args{strings.NewReader(`3 4
5 5 5 5
5 5 5 5
5 5 5 5
`)},
			`1
`,
			false,
		},

		// 6. Полная последовательность 1-16
		{
			"6_full_sequence_1_16",
			args{strings.NewReader(`4 4
1 2 3 4
8 7 6 5
9 10 11 12
16 15 14 13
`)},
			`16
`,
			false,
		},

		// 7. Цепочка длины 2
		{
			"7_chain_length_2",
			args{strings.NewReader(`2 2
1 3
2 4
`)},
			`2
`,
			false,
		},

		// 8. Длинная вертикальная цепочка
		{
			"8_long_vertical_chain",
			args{strings.NewReader(`5 2
1 2
3 4
5 6
7 8
9 10
`)},
			`2
`,
			false,
		},

		// 9. Длинная горизонтальная цепочка
		{
			"9_long_horizontal_chain",
			args{strings.NewReader(`2 5
1 2 3 4 5
10 9 8 7 6
`)},
			`10
`,
			false,
		},

		// 10. Цепочка с нулем
		{
			"10_chain_with_zero",
			args{strings.NewReader(`3 3
0 1 2
5 4 3
6 7 8
`)},
			`9
`,
			false,
		},

		// 11. Несколько цепочек разной длины
		{
			"11_multiple_chains",
			args{strings.NewReader(`3 3
1 2 1
4 3 4
5 6 5
`)},
			`6
`,
			false,
		},

		// 12. Одна строка с последовательностью
		{
			"12_single_row_sequence",
			args{strings.NewReader(`1 5
1 2 3 4 5
`)},
			`5
`,
			false,
		},

		// 13. Один столбец с последовательностью
		{
			"13_single_column_sequence",
			args{strings.NewReader(`5 1
1
2
3
4
5
`)},
			`5
`,
			false,
		},

		// 14. Зигзагообразная цепочка
		{
			"14_zigzag_chain",
			args{strings.NewReader(`3 3
1 3 2
4 6 5
7 9 8
`)},
			`2
`,
			false,
		},

		// 15. Максимальная цепочка не из самого большого числа
		{
			"15_max_chain_not_from_max_value",
			args{strings.NewReader(`3 3
9 1 2
8 7 3
5 4 6
`)},
			`3
`,
			false,
		},

		// 16. Большие числа с длинной цепочкой
		{
			"16_large_numbers_long_chain",
			args{strings.NewReader(`2 3
1000000 1000001 1000002
1000005 1000004 1000003
`)},
			`6
`,
			false,
		},

		// 17. Изолированные цепочки
		{
			"17_isolated_chains",
			args{strings.NewReader(`3 3
1 5 2
4 3 6
7 8 9
`)},
			`3
`,
			false,
		},

		// 18. Препятствие в середине
		{
			"18_obstacle_in_middle",
			args{strings.NewReader(`3 4
1 2 3 4
10 9 8 5
11 12 7 6
`)},
			`12
`,
			false,
		},

		// 19. Все элементы убывают - цепочки длины 1
		{
			"19_all_decreasing",
			args{strings.NewReader(`3 3
9 8 7
6 5 4
3 2 1
`)},
			`3
`,
			false,
		},

		// 20. Только одна возможная цепочка длины 2
		{
			"20_single_chain_length_2",
			args{strings.NewReader(`3 3
1 1 1
1 2 1
1 1 1
`)},
			`2
`,
			false,
		},

		// 21. Цепочка в углу
		{
			"21_corner_chain",
			args{strings.NewReader(`3 3
1 2 1
1 1 1
1 1 1
`)},
			`2
`,
			false,
		},

		// 22. Две отдельные цепочки одинаковой длины
		{
			"22_two_equal_chains",
			args{strings.NewReader(`3 3
1 2 4
3 5 5
4 3 2
`)},
			`3
`,
			false,
		},

		// 23. Спиральная цепочка
		{
			"23_spiral_chain",
			args{strings.NewReader(`3 3
1 2 3
8 9 4
7 6 5
`)},
			`9
`,
			false,
		},

		// 24. Минимальная матрица 2x2 с цепочкой
		{
			"24_minimal_2x2_with_chain",
			args{strings.NewReader(`2 2
1 2
4 3
`)},
			`4
`,
			false,
		},

		// 25. Только один элемент может продолжить цепочку
		{
			"25_single_chain_continuation",
			args{strings.NewReader(`3 3
1 1 1
1 2 1
1 3 1
`)},
			`3
`,
			false,
		},
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
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

func BenchmarkSolve(b *testing.B) {
	n, m := 1000, 1000
	mx := makeMatrix[int32](n, m)
	for i := range n {
		for j := range m {
			mx[i][j] = rand.Int32N(int32(10*n*m)) + 1
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		solve(mx)
	}
}
