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
			args{strings.NewReader(`(a+b+c)*(d-a)
`)},
			strings.TrimLeft(`
         .----[*]----.   
         |           |   
   .----[+]-.     .-[-]-.
   |        |     |     |
.-[+]-.     c     d     a
|     |                  
a     b                  
`, "\r\n"),
			true,
		},
		{
			"2",
			args{strings.NewReader(`a
`)},
			`a
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`a+b
`)},
			strings.TrimLeft(`
.-[+]-.
|     |
a     b
`, "\r\n"),
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
				t.Errorf("run() = \n%v, want \n%v", gotOut, tt.wantOut)
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

func Test_InfixToPostfix(t *testing.T) {
	tests := []struct {
		expr string
		want []byte
	}{
		// Базовые бинарные операции
		{"a+b", []byte("ab+")},
		{"a-b", []byte("ab-")},
		{"a*b", []byte("ab*")},
		{"a/b", []byte("ab/")},
		{"a^b", []byte("ab^")},

		// Цепочки одинаковых операций (ассоциативность)
		{"a+b+c", []byte("ab+c+")},
		{"a-b-c", []byte("ab-c-")},
		{"a*b*c", []byte("ab*c*")},
		{"a/b/c", []byte("ab/c/")},
		{"a^b^c", []byte("abc^^")}, // правоассоциативность
		{"a^b^c^d", []byte("abcd^^^")},

		// Смешанные приоритеты без скобок
		{"a+b*c", []byte("abc*+")},
		{"a*b+c", []byte("ab*c+")},
		{"a+b*c-d", []byte("abc*+d-")},
		{"a*b+c/d", []byte("ab*cd/+")},
		{"a+b*c^d", []byte("abcd^*+")},
		{"a^b*c+d", []byte("ab^c*d+")},
		{"a+b^c*d", []byte("abc^d*+")},

		// Скобки меняют порядок
		{"(a+b)*c", []byte("ab+c*")},
		{"a*(b+c)", []byte("abc+*")},
		{"(a+b)*(c+d)", []byte("ab+cd+*")},
		{"(a+b+c)*(d-a)", []byte("ab+c+da-*")},
		{"a^(b+c)", []byte("abc+^")},
		{"(a^b)+c", []byte("ab^c+")},
		{"a+(b^c)", []byte("abc^+")},
		{"(a+b)^(c+d)", []byte("ab+cd+^")},
		{"a*(b+c^d)", []byte("abcd^+*")},
		{"(a*b)+(c/d)", []byte("ab*cd/+")},
		{"a/((b+c)*d)", []byte("abc+d*/")},

		// Вложенные скобки
		{"((a+b)*c)+d", []byte("ab+c*d+")},
		{"a+((b+c)*d)", []byte("abc+d*+")},
		{"(a+(b*c))^d", []byte("abc*+d^")},
		{"a^(b+(c*d))", []byte("abcd*+^")},

		// Унарный минус не поддерживается (по условию — только бинарные операции и переменные),
		// поэтому такие случаи не включаем.

		// Граничные случаи
		{"a", []byte("a")},
		{"z", []byte("z")},
		{"a+b-c*d/e^f", []byte("ab+cd*ef^/-")},

		// Длинное выражение с разными приоритетами и скобками
		{"(a+b*c)-(d/e^f)+g", []byte("abc*+def^/-g+")},
		{"a*(b+c^d-e)/f", []byte("abcd^+e-*f/")},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			if got := infixToPostfix(tt.expr); !bytes.Equal(got, tt.want) {
				t.Errorf("InfixToPostfix(%q) = %q, want %q", tt.expr, got, tt.want)
			}
		})
	}
}

func Test_PostfixToAST(t *testing.T) {
	tests := []struct {
		polish   string
		wantTree string
		wantRest []byte
	}{
		// Простой операнд
		{"a", "{a <nil> <nil>}", []byte{}},

		// Бинарная операция
		{"ab+", "{+ {a <nil> <nil>} {b <nil> <nil>}}", []byte{}},

		// Цепочка: a + b * c → ОПН: abc*+
		{"abc*+", "{+ {a <nil> <nil>} {* {b <nil> <nil>} {c <nil> <nil>}}}", []byte{}},

		// (a + b) * c → ОПН: ab+c*
		{"ab+c*", "{* {+ {a <nil> <nil>} {b <nil> <nil>}} {c <nil> <nil>}}", []byte{}},

		// a ^ b ^ c → ОПН: abc^^ (правоассоциативно)
		{"abc^^", "{^ {a <nil> <nil>} {^ {b <nil> <nil>} {c <nil> <nil>}}}", []byte{}},

		// Лишние символы в начале (остаток)
		{"xyabc*+", "{+ {a <nil> <nil>} {* {b <nil> <nil>} {c <nil> <nil>}}}", []byte("xy")},

		// // Один оператор, не хватает операндов → поведение undefined, но не паникует
		// {"+", "{+ <nil> <nil>}", []byte{}}, // или можно ожидать панику — зависит от требований

		// // Пустой ввод
		// {"", "<nil>", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.polish, func(t *testing.T) {
			root, rest := postfixToAST([]byte(tt.polish))
			gotTree := root.String()
			if gotTree != tt.wantTree {
				t.Errorf("buildTree(%q) tree = %q, want %q", tt.polish, gotTree, tt.wantTree)
			}
			if !bytes.Equal(rest, tt.wantRest) {
				t.Errorf("buildTree(%q) rest = %q, want %q", tt.polish, rest, tt.wantRest)
			}
		})
	}
}

func TestParseExpr(t *testing.T) {
	tests := []struct {
		input    string
		wantTree string
	}{
		{"a", "{a <nil> <nil>}"},
		{"a+b", "{+ {a <nil> <nil>} {b <nil> <nil>}}"},
		{"a*b+c", "{+ {* {a <nil> <nil>} {b <nil> <nil>}} {c <nil> <nil>}}"},
		{"(a+b)*c", "{* {+ {a <nil> <nil>} {b <nil> <nil>}} {c <nil> <nil>}}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			root := parseExprToAST(tt.input)
			if got := root.String(); got != tt.wantTree {
				t.Errorf("ParseExpr(%q) = %q, want %q", tt.input, got, tt.wantTree)
			}
		})
	}
}
