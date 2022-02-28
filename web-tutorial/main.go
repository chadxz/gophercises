package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := fmt.Sprintf("%s.txt", p.Title)
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := fmt.Sprintf("%s.txt", title)
	body, err := os.ReadFile(filename)
	if err != nil {
		return &Page{}, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(w, "<h1>%s</h1><main>%s</main>", title, p.Body)
	if err != nil {
		panic(err)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	_, err = fmt.Fprintf(w, `<h1>Editing %s</h1>
<form action="/save/%s" method="post">
	<textarea name="body">%s</textarea>
	<input type="submit" value="Save" />
</form>`, title, title, p.Body)
	if err != nil {
		panic(err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/view/%s", title), http.StatusFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
