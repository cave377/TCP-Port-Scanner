package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func worker_varredura(portas chan int, alvo string, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	for porta := range portas {
		end := fmt.Sprintf("%s:%d", alvo, porta)
		conexao, err := net.DialTimeout("tcp", end, 5*time.Second)
		if err == nil {
			fmt.Println("\nPorta aberta: ", porta)
			conexao.Close()
		}
		bar.Add(1)
		wg.Done()
	}
}

func main() {
	alvo := flag.String("u", "", "url do alvo")
	porta := flag.String("p", "default", "porta para escaneamento (default escaneia da 1-1000)")
	threads := flag.Int("t", 50, "quantidade de workers (default: 50)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: main.go [opções]\n\n")
		fmt.Fprintln(os.Stderr, "\nOpções:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExemplos:")
		fmt.Fprintf(os.Stderr, "-u example.com -p 22\n")
		fmt.Fprintf(os.Stderr, "-u example.com -p 22-100 -t 20\n")
		fmt.Fprintf(os.Stderr, "-u example.com -p 22,53,135\n")
	}
	flag.Parse()

	if *alvo == "" {
		fmt.Fprintf(os.Stderr, "erro: a flag -u é obrigatória \n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("## Iniciando TCP-SCAN ##")
	fmt.Println("- Alvo -> ", *alvo)
	fmt.Println("- Quantidade de Threads -> ", *threads)

	workers_portas := make(chan int, *threads)
	var wg sync.WaitGroup

	for i := 0; i < cap(workers_portas); i++ {
		go func() {
		}()
	}

	switch {
	case strings.Contains(*porta, "-"):
		fmt.Println("- Portas a verificar -> ", *porta)
		tratado := strings.Split(*porta, "-")
		i, err1 := strconv.Atoi(strings.TrimSpace(tratado[0]))
		f, err2 := strconv.Atoi(strings.TrimSpace(tratado[1]))
		if err1 != nil || err2 != nil || i > f {
			fmt.Fprintln(os.Stderr, "erro: intervalo inválido, use 22-100")
			os.Exit(1)
		}

		total := f - i + 1
		bar := progressbar.Default(int64(total))
		for w := 0; w < cap(workers_portas); w++ {
			go worker_varredura(workers_portas, *alvo, &wg, bar)
		}

		for cont := i; cont <= f; cont++ {
			wg.Add(1)
			workers_portas <- cont
		}
		wg.Wait()
		close(workers_portas)

	case strings.Contains(*porta, ","):
		fmt.Println("- Portas a verificar -> ", *porta)
		tratado := strings.Split(*porta, ",")

		total := len(tratado)
		bar := progressbar.Default(int64(total))
		for w := 0; w < cap(workers_portas); w++ {
			go worker_varredura(workers_portas, *alvo, &wg, bar)
		}

		for cont := 0; cont < len(tratado); cont++ {
			num, err := strconv.Atoi(strings.TrimSpace(tratado[cont]))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ignorando porta inválida: %q\n", tratado[cont])
				bar.Add(1)
				continue
			}
			wg.Add(1)
			workers_portas <- num
		}
		wg.Wait()
		close(workers_portas)

	case *porta == "default":
		fmt.Println("- Portas a verificar -> ", *porta)

		total := 1000
		bar := progressbar.Default(int64(total))
		for w := 0; w < cap(workers_portas); w++ {
			go worker_varredura(workers_portas, *alvo, &wg, bar)
		}

		for i := 1; i <= total; i++ {
			wg.Add(1)
			workers_portas <- i
		}
		wg.Wait()
		close(workers_portas)

	case !strings.ContainsAny(*porta, "-,"):
		fmt.Println("- Portas a verificar -> ", *porta)
		tratado, err := strconv.Atoi(*porta)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro: porta inválida")
			os.Exit(1)
		}

		bar := progressbar.Default(int64(1))
		for w := 0; w < cap(workers_portas); w++ {
			go worker_varredura(workers_portas, *alvo, &wg, bar)
		}

		wg.Add(1)
		workers_portas <- tratado
		wg.Wait()
		close(workers_portas)

	default:
		flag.Usage()
	}

	fmt.Println("\n\nSaindo...")
}
