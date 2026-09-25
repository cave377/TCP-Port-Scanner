package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func worker_varredura(portas chan int, alvo string, wg *sync.WaitGroup) {
	for porta := range portas {
		end := fmt.Sprintf("%s:%d", alvo, porta)
		conexao, err := net.Dial("tcp", end)
		if err == nil {
			fmt.Println("Porta aberta: ", porta)
			conexao.Close()
		}

		wg.Done()
	}
}

func main() {
	alvo := flag.String("u", "", "url do alvo")
	porta := flag.String("p", "default", "porta para escaneamento (default escaneia da 1-1000)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: main.go [opções]\n\n")
		fmt.Fprintln(os.Stderr, "\nOpções:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExemplos:")
		fmt.Fprintf(os.Stderr, "-u example.com -p 22\n")
		fmt.Fprintf(os.Stderr, "-u example.com -p 22-100\n")
		fmt.Fprintf(os.Stderr, "-u example.com -p 22,53,135\n")
	}
	flag.Parse()
	if *alvo == "" {
		fmt.Fprintf(os.Stderr, "erro: a flag -u é obrigatória \n")
		flag.Usage()
	} else {
		workers_portas := make(chan int, 100)
		var wg sync.WaitGroup
		fmt.Println("Alvo -> ", *alvo)
		for i := 0; i < cap(workers_portas); i++ {
			go worker_varredura(workers_portas, *alvo, &wg)
		}
		switch {
		case strings.Contains(*porta, "-"):
			fmt.Println("Portas a verificar -> ", *porta)
			tratado := strings.Split(*porta, "-")
			i, _ := strconv.Atoi(tratado[0])
			f, _ := strconv.Atoi(tratado[1])
			for cont := i; cont <= f; cont++ {
				wg.Add(1)
				workers_portas <- cont
			}
			wg.Wait()
			close(workers_portas)

		case strings.Contains(*porta, ","):
			fmt.Println("Portas a verificar -> ", *porta)
			tratado := strings.Split(*porta, ",")
			for cont := 0; cont < len(tratado); cont++ {
				wg.Add(1)
				num, _ := strconv.Atoi(tratado[cont])
				workers_portas <- num
			}
			wg.Wait()
			close(workers_portas)

		case *porta == "default":
			fmt.Println("Portas a verificar -> ", *porta)
			for i := 1; i <= 1000; i++ {
				wg.Add(1)
				workers_portas <- i
			}
			wg.Wait()
			close(workers_portas)

		case !strings.ContainsAny(*porta, "-,"):
			fmt.Println("Portas a verificar -> ", *porta)
			tratado, _ := strconv.Atoi(*porta)
			wg.Add(1)
			workers_portas <- tratado
			wg.Wait()
			close(workers_portas)

		default:
			flag.Usage()
		}
	}
	fmt.Println("\n\nSaindo...")
}
