package server

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"github.com/russross/blackfriday"
)

type Page struct {
    Title string
    Body  []byte
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)

        if m == nil {
            http.NotFound(w, r)
            return
        }

        fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)


    if err != nil {
        http.Redirect(w,r, "/edit/"+title, http.StatusFound)
    }

    p.Body = []byte(blackfriday.MarkdownBasic(p.Body))

    renderTemplate(w, "view", p)
}

func (p *Page) save() error {
    filename := p.Title + ".md"
    return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".md"
    body, err := os.ReadFile(filename)

    if err != nil {
        return nil, err
    }

    return &Page{Title: title, Body: body}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)

    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}

    err := p.save()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func homePage(w http.ResponseWriter, r *http.Request, title string) {
    p := &Page{Title: title, Body: []byte("home page")}

    renderTemplate(w, "home", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func Run() {
    http.HandleFunc("/", makeHandler(homePage))
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
