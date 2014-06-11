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
)

//Wiki structure:
type Page struct {
	Title string
	Body []byte
}
l
//Handle URL's prefixed with /view/
func viewHandler(w http.ResponseWriter, r *http.Request) {
    //Get the title from the URL
    title := r.URL.Path[len("/view/"):]
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
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    //Check if there were any issues with loading the page.
    if err != nil {
	p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func savehandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/*
 * Template Code:
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
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
    //http.HandleFunc("/save/", saveHandler)
    http.ListenAndServe(":8080", nil)
    
/*
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
*/
}
