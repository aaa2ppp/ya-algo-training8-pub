package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
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
			args{strings.NewReader(`9 1
0 0
1 0
1 1
0 1
-1 1
-1 0
-1 -1
0 -1
1 -1
`)},
			`12
`,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
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

func TestCalcSquares(t *testing.T) {
	tests := []struct {
		d    int32
		want []int32
	}{
		{0, []int32{0}},
		{1, []int32{0, 1}},
		{2, []int32{0, 1}},
		{3, []int32{0, 1}},
		{4, []int32{0, 1, 4}},
		{8, []int32{0, 1, 4}},
		{9, []int32{0, 1, 4, 9}},
		{15, []int32{0, 1, 4, 9}},
		{16, []int32{0, 1, 4, 9, 16}},
		{25, []int32{0, 1, 4, 9, 16, 25}},
		{100, []int32{0, 1, 4, 9, 16, 25, 36, 49, 64, 81, 100}},
		{1e8, nil}, // проверим отдельно длину и последний элемент
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("d=%d", tt.d), func(t *testing.T) {
			got := calcSquares(tt.d)

			if tt.d == 1e8 {
				// Для 1e8 не проверяем весь слайс, а только свойства
				if len(got) == 0 {
					t.Fatal("empty result for d=1e8")
				}
				if got[0] != 0 {
					t.Errorf("first element != 0: %d", got[0])
				}
				last := got[len(got)-1]
				if last != 1e8 {
					t.Errorf("last square != 1e8, got %d", last)
				}
				root := int32(math.Sqrt(float64(last)))
				if root*root != last {
					t.Errorf("last element %d is not a perfect square", last)
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("calcSquares(%d) = %v, want %v", tt.d, got, tt.want)
				}
			}

			// Общие инварианты
			for i, sq := range got {
				if sq < 0 || sq > tt.d {
					t.Errorf("element %d out of bounds [0, d]", sq)
				}
				r := int32(math.Sqrt(float64(sq)))
				if r*r != sq {
					t.Errorf("element %d is not a perfect square", sq)
				}
				if i > 0 && got[i-1] >= got[i] {
					t.Errorf("not strictly increasing at %d: %v", i, got)
				}
			}
		})
	}
}

func TestCalcSquareSumsSlow(t *testing.T) {
	tests := []struct {
		d    int32
		want [][2]int32
	}{
		{0, [][2]int32{{0, 0}}},
		{1, [][2]int32{{0, 1}}},
		{2, [][2]int32{{1, 1}}},
		{3, [][2]int32{}},
		{4, [][2]int32{{0, 2}}},
		{5, [][2]int32{{1, 2}}},
		{8, [][2]int32{{2, 2}}},
		{9, [][2]int32{{0, 3}}},
		{10, [][2]int32{{1, 3}}},
		{13, [][2]int32{{2, 3}}},
		{25, [][2]int32{{0, 5}, {3, 4}}},
		{50, [][2]int32{{1, 7}, {5, 5}}},
		{65, [][2]int32{{1, 8}, {4, 7}}},
		{325, [][2]int32{{1, 18}, {6, 17}, {10, 15}}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("d=%d", tt.d), func(t *testing.T) {
			got := calcSquareSumsSlow(tt.d)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcSquareSumsSlow(%d) = %v, want %v", tt.d, got, tt.want)
			}

			for _, p := range got {
				a, b := p[0], p[1]
				if a > b {
					t.Errorf("violation: a > b in %v", p)
				}
				if a*a+b*b != tt.d {
					t.Errorf("wrong sum in %v", p)
				}
			}
		})
	}
}

func TestCalcSquareSums(t *testing.T) {
	// Сравнение с эталоном на множестве значений
	testValues := []int32{
		0, 1, 2, 3, 4, 5, 8, 9, 10, 13, 16, 17, 25, 26, 50, 65, 85, 100,
		325, 425, 1000, 10000, 99999999, 100000000,
	}

	for _, d := range testValues {
		t.Run(fmt.Sprintf("d=%d", d), func(t *testing.T) {
			want := calcSquareSumsSlow(d)
			got := calcSquareSums(d)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("calcSquareSums(%d) = %v, want %v", d, got, want)
			}
		})
	}
}

func TestCalcSquareSums2(t *testing.T) {
	tests := []struct {
		d    int32
		want [][2]int32
	}{
		{0, [][2]int32{{0, 0}}},
		{1, [][2]int32{{0, 1}}},
		{2, [][2]int32{{1, 1}}},
		{3, [][2]int32{}},
		{4, [][2]int32{{0, 2}}},
		{5, [][2]int32{{1, 2}}},
		{8, [][2]int32{{2, 2}}},
		{9, [][2]int32{{0, 3}}},
		{10, [][2]int32{{1, 3}}},
		{13, [][2]int32{{2, 3}}},
		{25, [][2]int32{{0, 5}, {3, 4}}},
		{50, [][2]int32{{1, 7}, {5, 5}}},
		{65, [][2]int32{{1, 8}, {4, 7}}},
		{100, [][2]int32{{0, 10}, {6, 8}}},
		{325, [][2]int32{{1, 18}, {6, 7}, {10, 15}}},
		{100000000, [][2]int32{{0, 10000}}}, // 10^8 = (10^4)^2
		{99999999, [][2]int32{}},            // не представимо
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("d=%d", tt.d), func(t *testing.T) {
			got := calcSquareSums(tt.d)
			want := calcSquareSumsSlow(tt.d) // эталон

			if !reflect.DeepEqual(got, want) {
				t.Errorf("calcSquareSums(%d) = %v, want %v", tt.d, got, want)
			}

			// Дополнительные проверки гарантий:
			for _, pair := range got {
				if pair[0] > pair[1] {
					t.Errorf("pair %v violates ordering (a <= b)", pair)
				}
				a, b := pair[0], pair[1]
				if a*a+b*b != tt.d {
					t.Errorf("pair %v sum != %d", pair, tt.d)
				}
			}
		})
	}
}

// Фаззинг: быстрая проверка на корректность до 100k
func TestCalcSquareSums_Fuzz(t *testing.T) {
	const maxD = 100_000 // достаточно для уверенности; 1e6 тоже пройдёт за ~1-2 сек
	for d := int32(0); d <= maxD; d++ {
		got := calcSquareSums(d)
		want := calcSquareSumsSlow(d)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch at d=%d: got %v, want %v", d, got, want)
		}
	}
}

func TestCalcSquareSums_Random(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000; i++ {
		d := int32(rand.Intn(100_000_000-100_000) + 100_001) // [100_000, 1e8]
		got := calcSquareSums(d)
		want := calcSquareSumsSlow(d)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch at d=%d: got %v, want %v", d, got, want)
		}
	}
}

func TestCalcOffsets(t *testing.T) {
	tests := []struct {
		name string
		d    int32
		want [][2]int32
	}{
		{
			name: "d=0 → [[0,0]]",
			d:    0,
			want: [][2]int32{{0, 0}}, // но d>=1, так что опционально
		},
		{
			name: "d=1 → [[0,1]]",
			d:    1,
			want: [][2]int32{{0, 1}, {0, -1}, {1, 0}, {-1, 0}},
		},
		{
			name: "d=2 → [[1,1]]",
			d:    2,
			want: [][2]int32{{1, 1}, {1, -1}, {-1, -1}, {-1, 1}},
		},
		{
			name: "d=5 → [[1,4]]",
			d:    5,
			want: [][2]int32{
				{1, 2}, {1, -2}, {-1, -2}, {-1, 2},
				{2, 1}, {2, -1}, {-2, -1}, {-2, 1},
			},
		},
		{
			name: "d=25 → [[0,25], [9,16]]",
			d:    25,
			want: [][2]int32{
				// from [0,25] → a=0, b=5
				{0, 5}, {0, -5}, {5, 0}, {-5, 0},
				// from [9,16] → a=3, b=4
				{3, 4}, {3, -4}, {-3, -4}, {-3, 4},
				{4, 3}, {4, -3}, {-4, -3}, {-4, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcOffsets(tt.d)

			// Приведём к множеству для сравнения без учёта порядка
			gotSet := make(map[[2]int32]bool)
			for _, v := range got {
				gotSet[v] = true
			}
			wantSet := make(map[[2]int32]bool)
			for _, v := range tt.want {
				wantSet[v] = true
			}

			if !reflect.DeepEqual(gotSet, wantSet) {
				t.Errorf("mismatch:\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestMaxOffsets(t *testing.T) {
	cases := [][2]int32{
		{85399925, 288},
		{77068225, 288},
		{72325565, 256},
		{58482125, 256},
		{48612265, 256},
		{88843105, 256},
		{48868625, 256},
		{93485125, 256},
		{82681625, 256},
		{41907125, 256},
		{92412125, 256},
		{60029125, 256},
		{80913625, 256},
		{71488625, 256},
		{69633785, 256},
		{65692250, 256},
		{29641625, 256},
		{62840245, 256},
		{84919250, 256},
		{83814250, 256},
		{97224530, 256},
		{77709125, 256},
		{86553545, 256},
		{54172625, 256},
		{80144545, 256},
		{32846125, 256},
		{95910685, 256},
		{87322625, 256},
		{42459625, 256},
		{97737250, 256},
		{89311625, 256},
		{62349625, 256},
		{59283250, 256},
		{74615125, 256},
		{90969125, 256},
		{69090125, 256},
		{71300125, 256},
		{99146125, 256},
		{90527125, 256},
		{86880625, 240},
		{66438125, 240},
		{73620625, 240},
		{68095625, 240},
		{96273125, 240},
		{95168125, 240},
		{52073125, 240},
	}
	for _, d := range cases {
		offsets := calcOffsets(d[0])
		if len(offsets) != int(d[1]) {
			t.Errorf("d=%d: calcOffsets()= %d, want %d", d, len(offsets), d[1])
		}
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		d     int32
		trees []Tree
		want  int32
	}{
		{
			name:  "empty",
			d:     1,
			trees: []Tree{},
			want:  0,
		},
		{
			name:  "single tree",
			d:     5,
			trees: []Tree{{0, 0}},
			want:  0,
		},
		{
			name:  "d=1, horizontal pair",
			d:     1,
			trees: []Tree{{0, 0}, {1, 0}},
			want:  1,
		},
		{
			name:  "d=1, vertical pair",
			d:     1,
			trees: []Tree{{0, 0}, {0, 1}},
			want:  1,
		},
		{
			name:  "d=1, no pair",
			d:     1,
			trees: []Tree{{0, 0}, {2, 0}},
			want:  0,
		},
		{
			name:  "d=2, diagonal (1,1)",
			d:     2,
			trees: []Tree{{0, 0}, {1, 1}},
			want:  1,
		},
		{
			name:  "d=2, two diagonal pairs",
			d:     2,
			trees: []Tree{{0, 0}, {1, 1}, {-1, -1}},
			want:  2, // (0,0)-(1,1), (0,0)-(-1,-1)
		},
		{
			name:  "d=5, mixed (1,2)",
			d:     5,
			trees: []Tree{{0, 0}, {1, 2}, {2, 1}, {-1, -2}},
			want:  3, // (0,0) с тремя остальными
		},
		{
			name: "d=25, multiple representations: 0²+5² and 3²+4²",
			d:    25,
			trees: []Tree{
				{0, 0},
				{5, 0}, {-5, 0}, {0, 5}, {0, -5}, // from 0²+5²
				{3, 4}, {3, -4}, {-3, 4}, {-3, -4},
				{4, 3}, {4, -3}, {-4, 3}, {-4, -3},
			},
			want: 12, // (0,0) с 4 + 8 = 12 соседями
		},
		{
			name: "d=25, full graph: every point connected to origin only",
			d:    25,
			trees: []Tree{
				{0, 0}, {5, 0}, {3, 4},
			},
			want: 2,
		},
		{
			name:  "d=100000000 (1e8 = 10000² + 0²)",
			d:     100000000,
			trees: []Tree{{0, 0}, {10000, 0}, {0, -10000}},
			want:  2,
		},
		{
			name:  "d=99999999 (not sum of two squares)",
			d:     99999999,
			trees: []Tree{{0, 0}, {1, 2}, {3, 4}},
			want:  0,
		},
		{
			name: "d=0 — but d>=1 per constraints, skip",
		},
	}

	for _, tt := range tests {
		if tt.name == "d=0 — but d>=1 per constraints, skip" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got := solve(tt.d, tt.trees)
			if got != tt.want {
				t.Errorf("solve(d=%d, trees=%v) = %d, want %d", tt.d, tt.trees, got, tt.want)
			}
		})
	}
}

func bruteForceSolve(d int32, trees []Tree) int32 {
	count := 0
	n := len(trees)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := trees[i].X - trees[j].X
			dy := trees[i].Y - trees[j].Y
			if dx*dx+dy*dy == d {
				count++
			}
		}
	}
	return int32(count)
}

func TestSolve_Fuzz_BruteForce(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	const trials = 10000

	for _, maxN := range []int{2, 5, 10, 20, 40, 80, 160, 320} {
		for trial := 0; trial < trials; trial++ {
			// Случайное d в разумных пределах
			d := int32(rand.Intn(200) + 1) // d <= 200 → макс. расстояние ~14

			// Случайное число деревьев
			n := rand.Intn(maxN + 1)
			trees := make([]Tree, n)
			for i := 0; i < n; i++ {
				x := rand.Intn(21) - 10 // [-10, 10]
				y := rand.Intn(21) - 10
				trees[i] = Tree{X: int32(x), Y: int32(y)}
			}

			want := bruteForceSolve(d, trees)
			got := solve(d, trees)

			if got != want {
				t.Fatalf("FAIL on trial %d: d=%d, trees=%v\nsolve=%d, brute=%d", trial, d, trees, got, want)
			}
		}
	}
}

func BenchmarkSolve_WorstCase(b *testing.B) {
	// Генерируем 100k уникальных точек в диапазоне [-1e8, 1e8]
	trees := make([]Tree, 100_000)
	for i := range trees {
		// Чтобы избежать коллизий в map, делаем уникальные точки
		trees[i] = Tree{
			X: int32(i),
			Y: int32(i * 2),
		}
	}

	d := int32(88830625) // 180 offsets

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(d, trees)
	}
}
