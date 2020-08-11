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

// renderizaTemplate faz o parse dos arquivos tratados pelos handlers
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

// E se você visitar /view/umaPaginaQueNaoExiste ?
// Você verá uma página contendo HTML. Isso ocorre porque ele ignora o valor de retorno do
// erro e continua tentando preencher o template sem dados. Em vez disso, se a página solicitada não existir,
// ele deve redirecionar o cliente para a página de edição para que o conteúdo possa ser criado:

// A função http.Redirect adiciona um código de status HTTP de http.StatusFound (302) e um Location
// header à resposta HTTP.

// viewHandler r o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/view/"):]
	pagina, err := carregaPagina(titulo)
	if err != nil {
		http.Redirect(escrever, ler, "/edit/"+titulo, http.StatusFound)
		return
	}
	renderizaTemplate(escrever, "view", pagina)
}

// editHandler carrega um formulário de edição
func editHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/edit/"):]
	pagina, err := carregaPagina(titulo)
	if err != nil {
		pagina = &Pagina{Titulo: titulo}
	}
	renderizaTemplate(escrever, "edit", pagina)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
