package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

func main() {
	_ = bufio.NewReader(os.Stdin)
	// no flush method

	writer := bufio.NewWriter(os.Stdout)
	writer.Flush()

	_ = strings.NewReader("awe") // обертка над строкой, предоставляющая reader api

	_ = bytes.NewBuffer([]byte("awe")) // аналогичная обертка над слайсом байт

	_ = bufio.NewScanner(os.Stdin)

	//file, err := os.Open("")
	//file.Read()

	//st, err := file.Stat()
	//st.Sys()
}
