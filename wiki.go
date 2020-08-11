package main

import (
	// Pacote para auxilio no tratamento e criação de erros
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	// pacote de expressões regulares
	"regexp"
)

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

// Como você deve ter observado, este programa tem uma falha de segurança séria: um usuário pode fornecer um
// caminho arbitrário para ser lido/escrito no servidor. Para atenuar isso, podemos escrever uma função para
// validar o título com uma expressão regular.

// Primeiro, adicione "regexp" à lista import. Então, podemos criar uma variável global para armazenar nossa
// expressão de validação:

// A função regexp.MustCompile irá analisar e compilar a expressão regular e retornar um regexp.Regexp.
// MustCompile é diferente de Compile porque entrará em panic se a compilação da expressão falhar, enquanto
// Compile retorna um error como um segundo parâmetro.

var caminhovalido = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// templates faz o cacheamento dos templates
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// Agora, vamos escrever uma função que usa a expressão caminhovalido para validar o caminho e extrair
// o título da página:

// Se o título coincide, ele será retornado junto com um valor nil de erro. Se o título for inválido,
// a função gravará um erro "404 Not Found" na conexão HTTP e retornará um erro ao handler.
// Para criar um novo erro customizado, temos que importar o pacote errors.

func obtemTitulo(escrever http.ResponseWriter, ler *http.Request) (string, error) {
	coincide := caminhovalido.FindStringSubmatch(ler.URL.Path)
	if coincide == nil {
		http.NotFound(escrever, ler)
		return "", errors.New("Titulo da página inválido")
	}
	return coincide[2], nil // The title is the second subexpression.
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

// Vamos colocar uma chamada para obtemTitulo em cada um dos handlers:

// viewHandler r o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo, err := obtemTitulo(escrever, ler)
	if err != nil {
		return
	}
	pagina, err := carregaPagina(titulo)
	if err != nil {
		http.Redirect(escrever, ler, "/edit/"+titulo, http.StatusFound)
		return
	}
	renderizaTemplate(escrever, "view", pagina)
}

// Vamos colocar uma chamada para obtemTitulo em cada um dos handlers:

// editHandler carrega um formulário de edição
func editHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo, err := obtemTitulo(escrever, ler)
	if err != nil {
		return
	}
	pagina, err := carregaPagina(titulo)
	if err != nil {
		pagina = &Pagina{Titulo: titulo}
	}
	renderizaTemplate(escrever, "edit", pagina)
}

// Vamos colocar uma chamada para obtemTitulo em cada um dos handlers:

// A função saveHandler tratará do envio de formulários localizados nas páginas de edição
func saveHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo, err := obtemTitulo(escrever, ler)
	if err != nil {
		return
	}
	corpo := ler.FormValue("body")
	pagina := &Pagina{Titulo: titulo, Corpo: []byte(corpo)}
	pagina.salvar()
	err = pagina.salvar()
	if err != nil {
		http.Error(escrever, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(escrever, ler, "/view/"+titulo, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
