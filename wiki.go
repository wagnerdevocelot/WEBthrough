package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

// salvar persiste os dados da página
func (p *Pagina) salvar() error {
	nomeDoArquivo := p.Titulo + ".txt"
	return ioutil.WriteFile(nomeDoArquivo, p.Corpo, 0600)
}

// Além de salvar páginas, também queremos carregar páginas:
func carregaPagina(titulo string) (*Pagina, error) {
	nomeDoArquivo := titulo + ".txt"
	corpo, err := ioutil.ReadFile(nomeDoArquivo)
	if err != nil {
		return nil, err
	}
	return &Pagina{Titulo: titulo, Corpo: corpo}, nil
}

// viewHandler escreve o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/view/"):]
	pagina, _ := carregaPagina(titulo)
	fmt.Fprintf(escrever, "<h1>%s</h1><div>%s</div>", pagina.Titulo, pagina.Corpo)
}

// Para usar esse Handler, reescrevemos nossa main func para inicializar http usando o viewHandler para
// manipular todas as solicitações no caminho /view/.

func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Vamos criar alguns dados para Pagina, crie um novo arquivo teste.txt no diretório de wiki.go e escreva
// Olá Mundo dentro dele, sem aspas. Vamos compilar o nosso código e tentar servir uma página wiki.

// $ go build wiki.go
// $ ./wiki

// Com este servidor da web em execução, uma visita a http://localhost:8080/view/teste
// deve mostrar uma página intitulada "teste" contendo as palavras "Olá, mundo!".
// Se quiser um caminho diferente basta usar um arquivo com um nome diferente no lugar de teste
