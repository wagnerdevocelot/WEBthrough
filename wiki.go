package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

// caminhovalido oferece a obtemTitulo uma expressão regular que caso não combine causa um panic
var caminhovalido = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// templates faz o cacheamento dos templates
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// obtemTitulo usa a expressão regular em caminhovalido para validar um path
func obtemTitulo(escrever http.ResponseWriter, ler *http.Request) (string, error) {
	coincide := caminhovalido.FindStringSubmatch(ler.URL.Path)
	if coincide == nil {
		http.NotFound(escrever, ler)
		return "", errors.New("Titulo da página inválido")
	}
	return coincide[2], nil // O titulo é a segunda subexpressão
}

// renderizaTemplate faz o parse dos arquivos tratados pelos handlers
func renderizaTemplate(escrever http.ResponseWriter, tmpl string, pagina *Pagina) {
	err := templates.ExecuteTemplate(escrever, tmpl+".html", pagina)
	if err != nil {
		http.Error(escrever, err.Error(), http.StatusInternalServerError)
	}
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

// Por fim, removemos as chamadas para obtemTitulo das funções de handler, tornando-as muito mais simples

// viewHandler r o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request, titulo string) {
	pagina, err := carregaPagina(titulo)
	if err != nil {
		http.Redirect(escrever, ler, "/edit/"+titulo, http.StatusFound)
		return
	}
	renderizaTemplate(escrever, "view", pagina)
}

// Por fim, removemos as chamadas para obtemTitulo das funções de handler, tornando-as muito mais simples

// editHandler carrega um formulário de edição
func editHandler(escrever http.ResponseWriter, ler *http.Request, titulo string) {
	pagina, err := carregaPagina(titulo)
	if err != nil {
		pagina = &Pagina{Titulo: titulo}
	}
	renderizaTemplate(escrever, "edit", pagina)
}

// Por fim, removemos as chamadas para obtemTitulo das funções de handler, tornando-as muito mais simples

// A função saveHandler tratará do envio de formulários localizados nas páginas de edição
func saveHandler(escrever http.ResponseWriter, ler *http.Request, titulo string) {
	corpo := ler.FormValue("body")
	pagina := &Pagina{Titulo: titulo, Corpo: []byte(corpo)}
	err := pagina.salvar()
	if err != nil {
		http.Error(escrever, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(escrever, ler, "/view/"+titulo, http.StatusFound)
}

// criaHandler usa obtemTitulo como argumento e encapsula todas as condições de erro dos outros Handlers
func criaHandler(função func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(escrever http.ResponseWriter, ler *http.Request) {
		coincide := caminhovalido.FindStringSubmatch(ler.URL.Path)
		if coincide == nil {
			http.NotFound(escrever, ler)
			return
		}
		função(escrever, ler, coincide[2])
	}
}

func main() {
	http.HandleFunc("/view/", criaHandler(viewHandler))
	http.HandleFunc("/edit/", criaHandler(editHandler))
	http.HandleFunc("/save/", criaHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Recompile o código e execute o aplicativo:

// $ go build wiki.go
// $ ./wiki

// Visite http://localhost:8080/view/ANewPage deve apresentar o formulário de edição da página. Você
// poderá inserir algum texto, clicar em 'Salvar' e ser redirecionado para a página recém-criada.
