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

// Existem vários locais em nosso programa onde os erros são ignorados. Esta é uma má prática, até porque,
// quando ocorre um erro, o programa terá um comportamento indesejado. A melhor solução é tratar os erros e
// retornar uma mensagem de erro ao usuário. Dessa forma, se algo der errado, o servidor funcionará
// exatamente como queremos e o usuário poderá ser notificado.

// Primeiro, vamos lidar com os erros em renderizaTemplate:

// A função http.Error envia um código de resposta HTTP especificado (neste caso "Erro interno do servidor")
// e uma mensagem de erro. A decisão de colocar isso em uma função separada já está valendo a pena pois
// tanto viewHandler e editHandler que fazem uso de renderizaTemplate possuem tratamento de erro.

// renderizaTemplate faz o parse dos arquivos tratados pelos handlers
func renderizaTemplate(escrever http.ResponseWriter, tmpl string, pagina *Pagina) {
	template, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(escrever, err.Error(), http.StatusInternalServerError)
		return
	}
	template.Execute(escrever, pagina)
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

// Agora vamos consertar saveHandler:

// A função salvar é uma função customizada que criamos ela não vem direto de packages do go então
// então ela não foi implementada com retorno de erro, por isso nesse caso atribuimos a função
// na variável err pra poder aplicar a mesma implementação de tratamento de erros padrão.

// Quaisquer erros que ocorram durante pagina.salvar() serão relatados ao usuário.

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
