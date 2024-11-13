package main

import (
	"flag"
	"fmt"
	"os"
)

type Interpreter struct {
	tape []byte
	ptr  int
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		tape: make([]byte, 30000),
		ptr:  0,
	}
}

func (i *Interpreter) Run(code string) {
	buf := make([]byte, 1)
	for j := 0; j < len(code); j++ {
		switch code[j] {
		case '>':
			// ポインタをインクリメントする
			if i.ptr < len(i.tape)-1 {
				i.ptr++
			}
		case '<':
			// ポインタをデクリメントする
			if i.ptr > 0 {
				i.ptr--
			}
		case '+':
			// ポインタが指す値をインクリメントする
			i.tape[i.ptr]++
		case '-':
			// ポインタが指す値をデクリメントする
			i.tape[i.ptr]--
		case '.':
			// ポインタが指す値を出力に書き出す
			fmt.Print(string(i.tape[i.ptr]))
		case ',':
			// 入力から1バイト読み込んで、ポインタが指す先に代入する
			_, err := os.Stdin.Read(buf)
			if err != nil {
				os.Exit(1)
			}
			i.tape[i.ptr] = buf[0]
		case '[':
			// ポインタが指す値が0なら、対応する ] の直後にジャンプする
			if i.tape[i.ptr] == 0 {
				depth := 1
				for depth > 0 {
					j++
					if j >= len(code) {
						fmt.Fprintf(os.Stderr, "missing matching ]")
						os.Exit(1)
					}
					if code[j] == '[' {
						depth++
					} else if code[j] == ']' {
						depth--
					}
				}
			}
		case ']':
			// ポインタが指す値が0でないなら、対応する [ （の直後）にジャンプする
			if i.tape[i.ptr] != 0 {
				depth := 1
				for depth > 0 {
					j--
					if j < 0 {
						fmt.Fprintf(os.Stderr, "missing matching [")
						os.Exit(1)
					}
					if code[j] == ']' {
						depth++
					} else if code[j] == '[' {
						depth--
					}
				}
			}
		}
	}
}

func readCode(file *os.File) string {
	buf := make([]byte, 1024)
	code := ""
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file: %v", err)
			os.Exit(1)
		}
		code += string(buf[:n])
	}
	return code
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing argument")
		os.Exit(1)
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	code := readCode(file)

	intprt := NewInterpreter()

	intprt.Run(code)
}
