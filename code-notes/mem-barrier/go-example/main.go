package main

import (
	"log"
	"runtime"
)

var (
	done bool
	str  string
)

func setup() {
	str = "hello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, world"

	done = true
	// потенциально возможен reordering операций(в данном случае присваивания):
	// 1. компилятором - оптимизация и перестановка независимых инструкций, если это не нарушит ход исполнения ОДНОПОТОЧНОГО кода
	// 2. процессором - конвейерная обработка, поход в ОЗУ в случае cache miss для адреса 1-ой инструкции и проч.
	// проблема в том, что компилятор/процессор не видят связь потоков между собой через переменные памяти
	// для предотвращения переупорядочивания используются барьеры памяти

	if done {
		log.Println("str:", str)
	}
}

func main() {
	go setup()

	for !done {
		runtime.Gosched()
	}

	if len(str) == 0 {
		log.Println("REORDERING! str:", str)
	}
	// log.Println("str main:", str)
}
