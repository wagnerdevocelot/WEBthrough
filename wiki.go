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

// O título da página (fornecido no URL) e campo do formulário, body são armazenados em uma nova Pagina.
// O método salvar() é então chamado para gravar os dados em um arquivo e o cliente é
// redirecionado para a /view/. O valor retornado por FormValue é do tipo string.
// Devemos converter esse valor para []byte antes que ele caiba na struct Pagina. Usamos []byte(corpo)
// para realizar a conversão.

// A função saveHandler tratará do envio de formulários localizados nas páginas de edição
func saveHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/save/"):]
	corpo := ler.FormValue("body")
	pagina := &Pagina{Titulo: titulo, Corpo: []byte(corpo)}
	pagina.salvar()
	http.Redirect(escrever, ler, "/view/"+titulo, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// para testar basta executar buildar o código novamente, executar ./wiki e abrir o browser na porta 8080

// O path em view e em seguida adicione um caminho que não esxiste, assim: http://localhost:8080/view/vapordev
// Seremos redirecionados para a pagina de Editando vapordev e quando salvar o conteudo da pagina estará
// em um arquivo .txt com o titulo tendo o nome do path e o conteudo com o corpo da edição.
