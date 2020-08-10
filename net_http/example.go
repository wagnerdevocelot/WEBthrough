// Aqui está um exemplo funcional completo de um servidor web simples:

package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá, eu adoro %s!", r.URL.Path[1:])
}

// A main func começa com uma chamada para http.HandleFunc, que informa ao pacote http para lidar com
// todas as solicitações para o web root ( "/" ) com handler.

// Em seguida chama http.ListenAndServe, especificando que deve escutar na porta 8080 em qualquer
// interface ( ":8080" ). (Não se preocupe com seu segundo parâmetro, nil por enquanto.) nil será
// bloqueado até que o programa seja encerrado.

// ListenAndServe sempre retorna um erro, e só retorna quando ocorre um erro inesperado.
// Para registrar esse erro, envolvemos com a chamada de função log.Fatal.

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Se você executar este programa e acessar o URL:

// http://localhost:8080/ubuntu

// o programa apresentaria uma página contendo:

// Olá, eu adoro ubuntu!
