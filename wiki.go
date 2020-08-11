package main

import (
	"fmt"

	// O html/template package faz parte da biblioteca padrão do Go. Podemos usar html/template
	// para manter o HTML em um arquivo separado, o que nos permite alterar o layout
	// de nossa página de edição sem modificar o código Go.

	// =======================================> Vá para o arquivo edit.html

	"html/template"
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

// Modifique editHandler para usar o template, em vez do HTML hardcoded:

// A função template.ParseFiles lerá o conteúdo de edit.html e retornará um ponteiro *template.Template.

// O método template.Execute executa o template, escrevendo o HTML gerado no http.ResponseWriter.
// Os identificadores .Titulo e .Corpo lá no edit.html referem-se a Titulo e Corpo
// presentes aqui na struct Pagina.

// As diretivas de temlate são colocadas entre chaves duplas. A instrução printf "%s" .Corpo
// é uma chamada de função que tem .Corpo como uma string em vez de um fluxo de bytes como saída.
// O html/template package ajuda a garantir que apenas HTML seguro e com aparência correta seja
// gerado por ações de template.

// editHandler carrega um formulário de edição
func editHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/edit/"):]
	pagina, err := carregaPagina(titulo)
	if err != nil {
		pagina = &Pagina{Titulo: titulo}
	}
	template, _ := template.ParseFiles("edit.html")
	template.Execute(escrever, pagina)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
