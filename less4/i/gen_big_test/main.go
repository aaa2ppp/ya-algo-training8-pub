package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Tree struct {
	X, Y int
}

func main() {
	d := 85399925 // 288 offsets
	
	if len(os.Args) > 1 {
		v, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("argument must be int: %v", err)
		}
		if v <= 0 {
			log.Fatalf("argument must be positive")
		}
		d = v
	}

	// Генерируем 100k уникальных точек в диапазоне [-1e8, 1e8]
	trees := make([]Tree, 100_000)
	for i := range trees {
		trees[i] = Tree{
			X: i,
			Y: i * 2,
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	fmt.Fprintln(w, len(trees), d)
	for _, t := range trees {
		fmt.Fprintln(w, t.X, t.Y)
	}
}
