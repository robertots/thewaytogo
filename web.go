package main

import (
    "fmt"
    "net/http"
    "time"
	"io/ioutil"
	"html/template"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
    filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	
	if err != nil {
		return nil, err
	}
	
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    	
    t, err := template.ParseFiles("./tmpl/" + tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    
	fmt.Println("\n New instance called at ", time.Now())    
    fmt.Println("\n write ---------------------------- \n")    
    fmt.Println(w)
    fmt.Println("\n reqst ---------------------------- \n")    
    fmt.Println(r)        
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}	
    }
    t, _ := template.ParseFiles("./tmpl/edit.html")
    t.Execute(w, p)
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
    title := "home"
    p, err := loadPage(title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "home", p)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    
    return func(w http.ResponseWriter, r *http.Request) {
    
        fmt.Println("\n reqst ---------------------------- \n")    
        //fmt.Println(r) 
        fmt.Println(validPath.FindStringSubmatch(r.URL.Path));
        
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func main() {
    // http.HandleFunc("/", handler)
    http.Handle("/misc/", http.StripPrefix("/misc/", http.FileServer(http.Dir("misc"))))
    
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/", homeHandler)
    http.ListenAndServe(":8080", nil)
}
