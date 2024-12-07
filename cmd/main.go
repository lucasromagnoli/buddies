package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		println(err.Error())
	}
}

func run(ctx context.Context) error {
	go dota2()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	return nil
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Lê o corpo da requisição
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Gera o nome do arquivo com a data e hora atual
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("data_%s.json", currentTime)

	// Cria o arquivo
	file, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Erro ao criar o arquivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Escreve o corpo da requisição no arquivo
	_, err = file.Write(body)
	if err != nil {
		http.Error(w, "Erro ao escrever no arquivo", http.StatusInternalServerError)
		return
	}

	// Não envia resposta ao cliente, apenas escreve no arquivo
	fmt.Println("Corpo da requisição gravado em:", fileName)
}

func dota2() {
	// Rota para /post
	http.HandleFunc("/", postHandler)

	fmt.Println("Servidor rodando na porta 44444...")
	err := http.ListenAndServe(":44444", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
