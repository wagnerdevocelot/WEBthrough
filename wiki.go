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

// Há uma ineficiência neste código: renderizaTemplate chama ParseFiles sempre que uma página é
// renderizada. Uma abordagem melhor seria chamar ParseFiles uma vez na inicialização do programa,
// analisando todos os templates em um único *Template. Então, podemos usar o método ExecuteTemplate para
// renderizar um template específico.

// Primeiro, criamos uma variável global chamada templates e a inicializamos com ParseFiles.

// A função template.Must é um invólucro que entra em pânico quando é passado um valor error
// não nulo (nil) e, caso contrário, retorna o *Template inalterado. O panic é apropriado aqui;
// se os templates não puderem ser carregados, a única coisa sensata a fazer é sair do programa.

// A função ParseFiles recebe qualquer número de argumentos de string que identificam nossos arquivos de
// template e os parseam (não sei se existe essa palavra haha) em templates que são nomeados
// após o nome base do arquivo. Se adicionássemos mais modelos ao nosso programa, adicionaríamos
// seus nomes aos argumentos da chamada ParseFiles.

// templates faz o cacheamento dos templates
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// Em seguida, modificamos a função renderizaTemplate para chamar o método templates.ExecuteTemplate
// com o nome do template apropriado:

// Observe que o nome do template é o nome do arquivo do template, portanto, devemos anexar
// ".html" ao argumento tmpl.

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

// A função saveHandler tratará do envio de formulários localizados nas páginas de edição
func saveHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/save/"):]
	corpo := ler.FormValue("body")
	pagina := &Pagina{Titulo: titulo, Corpo: []byte(corpo)}
	pagina.salvar()
	err := pagina.salvar()
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
