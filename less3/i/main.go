package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"unsafe"
)

type solveFunc func(string) string

func solve(expr string) string {
	tree := parseExprToAST(expr)
	h, w := tree.calcSizes()

	// make canvas
	n, m := h, w+1
	buf := make([]byte, n*m)
	for i := range buf {
		buf[i] = ' '
	}

	canvas := make([][]byte, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		canvas[i] = buf[j : j+m]
	}

	tree.printToCanvas(canvas, 0, 0)
	for i := range h {
		canvas[i][w] = '\n'
	}

	return unsafeString(buf)
}

// infixToPostfix converts an infix expression (e.g. "a+b*c") into postfix notation (e.g. "abc*+").
func infixToPostfix(expr string) []byte {
	var (
		polish []byte
		ops    stack[byte]
	)

	prior := func(c byte) int {
		switch c {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		case '^':
			return 3
		}
		return 0
	}

	for _, c := range []byte(expr) {
		switch c {
		case '(':
			ops.push(c)
		case ')':
			for !ops.empty() {
				c := ops.pop()
				if c == '(' {
					break
				}
				polish = append(polish, c)
			}
		case '+', '-', '*', '/':
			for !ops.empty() && prior(c) <= prior(ops.top()) {
				polish = append(polish, ops.pop())
			}
			ops.push(c)
		case '^':
			for !ops.empty() && prior(c) < prior(ops.top()) {
				polish = append(polish, ops.pop())
			}
			ops.push(c)
		default:
			polish = append(polish, c)
		}
	}

	for !ops.empty() {
		polish = append(polish, ops.pop())
	}

	return polish
}

// exprNode represents a node in the abstract syntax tree of an expression.
type exprNode struct {
	Value  byte
	Left   *exprNode
	Right  *exprNode
	Height int
	Width  int
}

func (n *exprNode) String() string {
	var b strings.Builder
	n.writeTo(&b)
	return b.String()
}

func (n *exprNode) writeTo(w *strings.Builder) {
	if n == nil {
		w.WriteString("<nil>")
		return
	}
	w.WriteByte('{')
	w.WriteByte(n.Value)
	w.WriteByte(' ')
	n.Left.writeTo(w)
	w.WriteByte(' ')
	n.Right.writeTo(w)
	w.WriteByte('}')
}

func (node *exprNode) calcSizes() (h, w int) {
	if node == nil {
		panic("calcNodeSize: node cannot be nil")
	}

	if node.Left == nil && node.Right == nil {
		node.Height = 1
		node.Width = 1
	} else {
		lh, lw := node.Left.calcSizes()
		rh, rw := node.Right.calcSizes()
		node.Height = max(lh, rh) + 2
		node.Width = lw + rw + 5
	}

	return node.Height, node.Width
}

func (node *exprNode) printToCanvas(canvas [][]byte, i, j int) int {
	if node == nil {
		panic("printToCanvas: node cannot be nil")
	}

	if node.Left == nil && node.Right == nil {
		canvas[i][j] = node.Value
		return j
	}

	center := j + node.Left.Width + 2

	lc := node.Left.printToCanvas(canvas, i+2, j)
	rc := node.Right.printToCanvas(canvas, i+2, center+3)

	jj := lc

	canvas[i][jj] = '.'
	jj++

	for ; jj < center-1; jj++ {
		canvas[i][jj] = '-'
	}

	copy(canvas[i][jj:], []byte{'[', node.Value, ']'})
	jj += 3

	for ; jj < rc; jj++ {
		canvas[i][jj] = '-'
	}

	canvas[i][jj] = '.'

	canvas[i+1][lc] = '|'
	canvas[i+1][rc] = '|'

	return center
}

// postfixToAST builds an abstract syntax tree from a postfix expression.
func postfixToAST(postfix []byte) (*exprNode, []byte) {
	n := len(postfix)
	if n == 0 {
		panic("postfixToAST: input postfix cannot be empty")
	}

	c := postfix[n-1]
	postfix = postfix[:n-1]

	switch c {
	case '+', '-', '*', '/', '^':
		right, postfix := postfixToAST(postfix)
		left, postfix := postfixToAST(postfix)
		return &exprNode{Value: c, Left: left, Right: right}, postfix
	default:
		return &exprNode{Value: c}, postfix
	}
}

// parseExprToAST parses an infix expression and returns its AST.
func parseExprToAST(expr string) *exprNode {
	postfix := infixToPostfix(expr)
	root, _ := postfixToAST(postfix)
	return root
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	log.SetFlags(0)
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	expr, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	expr = strings.TrimSpace(expr)

	ans := solve(expr)
	bw.WriteString(ans)
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout, solve)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// ----------------------------------------------------------------------------

type stack[T any] []T

func (s stack[T]) empty() bool {
	return len(s) == 0
}

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s stack[T]) top() T {
	n := len(s)
	return s[n-1]
}

func (s *stack[T]) pop() T {
	old := *s
	n := len(old)
	v := old[n-1]
	*s = old[:n-1]
	return v
}
