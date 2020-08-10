// Aqui está um exemplo funcional completo de um servidor web simples:

package main

import (
	"fmt"
	"net/http"
)

// A função handler é do tipo http.HandlerFunc. Leva um http.ResponseWriter e um http.Request
// como seus argumentos.

// Um valor http.ResponseWriter monta a resposta do servidor HTTP; ao escrever nele enviamos dados
// para o cliente HTTP. Nesse caso enviei "Olá, eu adoro %s!"

// Um http.Request é uma estrutura de dados que representa a solicitação HTTP do client.
// O r.URL.Path é o componente do caminho do URL da solicitação. O [1:] significa
// "criar um sub-slice do Path do primeiro caractere até o fim."
// Isso remove o "/" inicial do nome do caminho.

// Então qualquer coisa escrita no browser depois da raiz / é adicionado a mensagem de "Olá, eu adoro..."

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá, eu adoro %s!", r.URL.Path[1:])
}

func main() {

}
