// Implementação do algoritmo de travessia de Tarry
package main

import (
	"fmt"
	"time"
)

type Token struct {
	Sender string
}

type Neighbour struct {
	Id   string
	From chan Token
	To   chan Token
}

type Path struct {
	From  string
	Nodes [5]string
}

func redirect(in chan Token, neigh Neighbour) {
	token := <-neigh.From
	in <- token
}

func process(id string, token Token, neighs ...Neighbour) {
	var pai Neighbour

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Token, 1)
	nmap := make(map[string]Neighbour)
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(in, neigh)
	}

	if token.Sender == "init" {
		// Processo iniciador
		fmt.Printf("* %s é raiz.\n", id)
		token.Sender = id
		neighs[0].To <- token
		size := len(neighs)

		for i := 1; i < size; i++ {
			tk := <-in
			fmt.Printf("From %s to %s\n", tk.Sender, id)
			tk.Sender = id
			neighs[i].To <- tk
		}

		tk := <-in
		fmt.Printf("From %s to %s\n", tk.Sender, id)
		fmt.Println("Fim!")
	} else {
		// Processo não iniciador
		tk := <-in
		fmt.Printf("From %s to %s\n", tk.Sender, id)

		for _, neigh := range neighs {
			if pai.Id == "" {
				pai = nmap[tk.Sender]
				fmt.Printf("Explorando o processo %s \n", id)
			}

		// Entrega o token para o vizinho se ele não for o pai
		if pai.Id != neigh.Id {
			tk.Sender = id
			neigh.To <- tk
			tk = <-in
			fmt.Printf("From %s to %s\n", tk.Sender, id)
		}
		}

		// Token volta para o pai depois de ter passado enviado para todos os vizinhos
		tk.Sender = id
		pai.To <- tk
	}

}

func main() {

	pW := make(chan Token, 1)
	pS := make(chan Token, 1)
	pR := make(chan Token, 1)
	wP := make(chan Token, 1)
	wS := make(chan Token, 1)
	sP := make(chan Token, 1)
	sW := make(chan Token, 1)
	rQ := make(chan Token, 1)
	rP := make(chan Token, 1)
	qR := make(chan Token, 1)

	
	go process("W", Token{}, Neighbour{"S", sW, wS}, Neighbour{"P", pW, wP})
	go process("S", Token{}, Neighbour{"W", wS, sW},Neighbour{"P", pS, sP})
	go process("R", Token{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP})
	go process("Q", Token{}, Neighbour{"R", rQ, qR})
	go process("P", Token{"init"}, Neighbour{"W", wP, pW}, Neighbour{"S", sP, pS}, Neighbour{"R", rP, pR})
	time.Sleep(10 * time.Second);
	fmt.Printf("* Fim da execução do algoritmo.")
}
