package main

import (
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

// Observe que usamos quase exatamente o mesmo código de template em ambos os Handlers(view e edit).
// Vamos remover essa duplicação movendo o código de template para sua própria função

// Ela possui 3 parametros, ResponseWriter, tmpl que é uma string e um ponteiro para o struct Pagina
// O parseFiles vem pra função usando o segundo parametro como argumento e concatena com .html
// e então executa a função assim como faziamos nos Handler.

func renderizaTemplate(escrever http.ResponseWriter, tmpl string, pagina *Pagina) {
	template, _ := template.ParseFiles(tmpl + ".html")
	template.Execute(escrever, pagina)
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

// E modifique os Handlers para usar essa função

// viewHandler r o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/view/"):]
	pagina, _ := carregaPagina(titulo)
	// Agora só precisamos chamar a função com os parametros exigidos diminuindo a redundância
	renderizaTemplate(escrever, "view", pagina)
}

// E modifique os Handlers para usar essa função

// editHandler carrega um formulário de edição
func editHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/edit/"):]
	pagina, err := carregaPagina(titulo)
	if err != nil {
		pagina = &Pagina{Titulo: titulo}
	}
	// Agora só precisamos chamar a função com os parametros exigidos diminuindo a redundância
	renderizaTemplate(escrever, "edit", pagina)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Se comentarmos o registro do nosso saveHandler não implementado, podemos buildar e
// testar nosso programa mais uma vez.

// // $ go build wiki.go
// // $ ./wiki

// e na porta 8080 de localhost podemos mudar o path no browser entre view e edit
