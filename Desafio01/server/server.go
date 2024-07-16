package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	//_ "github.com/mattn/go-sqlite3"
)

/*
Requisitos desse modulo:
	Deverá consumir a API contendo o câmbio de Dólar e Real no
	endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e
	em seguida deverá retornar no formato JSON o resultado para o cliente

	Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida,
	sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e
	o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.

	Os 2 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.

	O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a
	ser utilizada pelo servidor HTTP será a 8080.

Extras:
	https://mholt.github.io/json-to-go/
*/

type BuscaMoeda struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	println("Executanto Server")

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", BuscaCotacao)
	http.ListenAndServe(":8080", mux)

	//http.HandleFunc("/cotacao", BuscaCotacao)
	//http.ListenAndServe(":8080", nil)
}

func BuscaCotacao(w http.ResponseWriter, r *http.Request) {

	// consumir a API contendo o câmbio de Dólar e Real

	c := http.Client{Timeout: time.Duration(900) * time.Millisecond}

	var url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
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

	var data BuscaMoeda
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao tratar resposta: %v\n", err)
	}

	fmt.Println(data.Usdbrl.Bid)

	// Salvando dados no Banco de Dados

	// Abrindo conexão com o banco de dados
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Defina um tempo limite no driver SQL, se necessário
	db.SetConnMaxLifetime(time.Millisecond * 10)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Crie uma tabela
	createTableSQL := `CREATE TABLE IF NOT EXISTS valores (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
		"valor" TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Insira um registro
	insertUserSQL := `INSERT INTO valores (valor) VALUES (?)`
	_, err = db.Exec(insertUserSQL, data.Usdbrl.Bid)
	if err != nil {
		log.Fatal(err)
	}

	// Retorno o valor do Dolar

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data.Usdbrl.Bid))
}
