package main

//

/*
Requisitos desse modulo:
	Deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

	O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON).
	Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

	O contexto devera retornar erro nos logs caso o tempo de execução seja insuficiente.

	O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
*/

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	println("Executanto Cliente")

	// consumir a API contendo o câmbio de Dólar e Real

	c := http.Client{Timeout: time.Duration(300) * time.Millisecond}

	var url = "http://localhost:8080/cotacao"
	//req, err := http.Get(url)
	req, err := c.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro na busca: %v\n", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro na resposta: %v\n", err)
	}

	//f, err := os.Create("coltacaodolar.txt")
	f, err := os.OpenFile("coltacaodolar.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	a, err := f.WriteString("Dolar: " + string(res) + "\n")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ok %d", a)
}
