/*
 *  Wiki from http://golang.org/doc/articles/wiki/
 *  Created by Taylor Brazelton
 */

package main

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
)

//Wiki structure:
type Page struct {
	Title string
	Body []byte
}

//Global Variables:
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Handle URL's prefixed with /view/
func viewHandler(w http.ResponseWriter, r *http.Request) {
    //Get the title from the URL
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    //Load the page from the file.
    p, err := loadPage(title)
    //Handle error if someone tries to view a non-existent file. Redirect them to the edit page so they can create one.
    if err != nil {
	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	return
    }
    //Render the template from the parameters in the object.
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    //Check if there were any issues with loading the page.
    if err != nil {
	p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/*
 * Template Code:
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/*
 * Validation Functions:
 */
//Validate the URL of which to open.
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}


/*
 *Data functions:
 */
//Persistent storage of the pages...
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//Load pages from persistent storage.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

//Main function
func main() {
    
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    http.ListenAndServe(":8080", nil)
    
/*
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
*/
}
