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

// Capturar a condição de erro em cada handler apresenta muitos códigos repetidos. E se pudéssemos
// envolver cada um dos handlers em uma função que faz essa validação e verificação de erro?
// As literal func de Go fornecem um meio poderoso de abstrair funcionalidades que pode nos ajudar aqui.

// Primeiro, redfinimos os argumentos da função de cada um dos handlers para aceitar uma string de titulo

// viewHandler r o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request, titulo string) {
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

// Primeiro, redfinimos os argumentos da função de cada um dos handlers para aceitar uma string de titulo

// editHandler carrega um formulário de edição
func editHandler(escrever http.ResponseWriter, ler *http.Request, titulo string) {
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

// Primeiro, redfinimos os argumentos da função de cada um dos handlers para aceitar uma string de titulo

// A função saveHandler tratará do envio de formulários localizados nas páginas de edição
func saveHandler(escrever http.ResponseWriter, ler *http.Request, titulo string) {
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

// Agora vamos definir uma função que encapsula uma função do tipo handler
// e retorna uma função do tipo http.HandlerFunc

// podemos pegar a função obtemTitulo e usá-la aqui como argumento (com algumas pequenas modificações)

// A função retornada é chamada de clojure porque contém valores definidos fora dela. Nesse caso,
// a variável função (o único argumento para criaHandler ) é anexada pelo clojure.
// A variável função será um de nossos handlers de salvar, editar ou visualizar. (save/edit/view)

// O clojure retornado por criaHandler é uma função que recebe um http.ResponseWriter e
// http.Request (em outras palavras, um http.HandlerFunc ).

// O clojure extrai o titulo do caminho da solicitação e o valida com o caminhovalido regexp.
// Se titulo for inválido, um erro será gravado no ResponseWriter usando a função http.NotFound.
// Se o titulo for válido, a função de handler encapsulada será o argumento função.
// chamado com o ResponseWriter, Request e titulo como argumentos.

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

// Agora podemos agrupar as funções de handler com criaHandler na main func,
// antes de serem registradas no pacote http

func main() {
	http.HandleFunc("/view/", criaHandler(viewHandler))
	http.HandleFunc("/edit/", criaHandler(editHandler))
	http.HandleFunc("/save/", criaHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
