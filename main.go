package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func varredura(inicio int, fim int, alvo string) {
	for door := inicio; door <= fim; door++ {
		address := fmt.Sprintf("%s:%d", alvo, door)
		conexao, err := net.Dial("tcp", address)
		if err == nil {
			fmt.Println("Porta aberta: ", door)
		}
		conexao.Close()
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
		fmt.Fprintf(os.Stderr, "-u http://www.example.com -p 22\n")
		fmt.Fprintf(os.Stderr, "-u http://www.example.com -p 22-100\n")
		fmt.Fprintf(os.Stderr, "-u http://www.example.com -p 22,53,135\n")
	}
	flag.Parse()
	if *alvo == "" {
		fmt.Fprintf(os.Stderr, "erro: a flag -u é obrigatória \n")
		flag.Usage()
	} else {
		fmt.Println("Alvo -> ", *alvo)
		switch {
		case strings.Contains(*porta, "-"):
			fmt.Println("Portas a verificar -> ", *porta)
			tratado := strings.Split(*porta, "-")
			i, _ := strconv.Atoi(tratado[0])
			f, _ := strconv.Atoi(tratado[1])
			varredura(i, f, *alvo)

		case strings.Contains(*porta, ","):
			fmt.Println("Portas a verificar -> ", *porta)
			tratado := strings.Split(*porta, ",")
			for cont := 0; cont < len(tratado); cont++ {
				door, _ := strconv.Atoi(tratado[cont])
				address := fmt.Sprintf("%s:%d", *alvo, door)
				conexao, err := net.Dial("tcp", address)
				if err == nil {
					fmt.Println("Porta aberta: ", tratado[cont])
				}
				conexao.Close()
			}
		case *porta == "default":
			fmt.Println("Portas a verificar -> ", *porta)
			varredura(1, 1000, *alvo)
		case !strings.ContainsAny(*porta, "-,"):
			fmt.Println("Portas a verificar -> ", *porta)
			tratado, _ := strconv.Atoi(*porta)
			address := fmt.Sprintf("%s:%d", *alvo, tratado)
			conexao, err := net.Dial("tcp", address)
			if err == nil {
				fmt.Println("Porta aberta: ", tratado)
			}
			conexao.Close()
		default:
			flag.Usage()
		}
	}
	fmt.Println("\n\nSaindo...")
}
