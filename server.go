// server.go

// to run the server start it from the command line in a terminal shell
// command  cd to where goserver is installed
// command  ./goserver

package main

import (
	"errors"
	"fmt"
	"html"

	"github.com/user/goserver/datamanager"
	//"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
)

// the Page Object definition
type Page struct {
	Title string
	Body  []byte
}

// regexp used to validate the requested path for security purposes
var validPath = regexp.MustCompile("^/([a-z.-])+/([a-zA-Z0-9./-]+)$")

/****************************************************************************************/
/*  the getAbsPath function accepts a file path as a string constant and validates the  */
/*  path. The path is converted to an absolute path, and returned to calling function   */
/****************************************************************************************/
func getAbsPath(path string) (string, error) {
	// get the current working directory
	currentDir, _ := os.Getwd()
	//fmt.Print("matching path : " + path + "\n")
	// validate the path
	m := validPath.FindStringSubmatch(path)
	// if path is not valid return URL is invalid error
	if m == nil {
		fmt.Print("match = nil\n")
		err := errors.New("URL is invalid\n")
		return "", err
	}
	//fmt.Print("match = " + m[0] + "\n")
	//return the valid absolute path and no errors
	return currentDir + m[0], nil
}

// to be imp[lemented
/*func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}*/

/****************************************************************************************/
/*  the getPage function accepts a file pathname as a string constant and converts it   */
/*  to an absolute path. The file is retrieved from storage and new Page object is      */
/*  created and returned to calling function                                            */
/****************************************************************************************/
func getPage(urlpath string) (*Page, error) {
	// get the absolute filepath
	filepath, err := getAbsPath(urlpath)
	// if the getAbsPath routine fails return the error
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	// else fill in the placeholders to construct the Page Object
	title := path.Base(filepath)
	body, err := ioutil.ReadFile(filepath)
	// if the io routine fails return the error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// else return the Page Object with no errors
	return &Page{Title: title, Body: body}, nil
}

/****************************************************************************************/
/*  the getDataHandler function queries the incoming http request for the URL query     */
/*  parameters, and passes the parameters to the appropriate routine handler defined    */
/*  by the index parameter.                                                             */
/****************************************************************************************/
func getDataHandler(w http.ResponseWriter, r *http.Request) {
	index := r.URL.Query().Get("index")
	fmt.Printf("index = %q\n", index)
	switch index {
	case "clubindex":
		{
			fmt.Printf("index = clubindex\n")
			user := "Guest"
			country := r.URL.Query().Get("country")
			club := r.URL.Query().Get("club")
			language := r.URL.Query().Get("lang")

			html := datamanager.ClubIndex(user, country, club, language)
			fmt.Fprintf(w, "%s", html)
		}
	default:
		fmt.Printf("Unrecognised index value : " + index)
	}
}

/****************************************************************************************/
/*  the getHandler function queries the incoming http request for the URL path of the   */
/*  static file, retrieves the file from storage and writes it to the http response     */
/*  and terminates the response, which is then returned by the server. All index files  */
/*  are stored in the index folder created in the same directory as the go server       */
/****************************************************************************************/
func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.URL.Path = "/html/index.html"
	}
	//if(!client is logged in)
	r.URL.Path = "/index" + r.URL.Path
	p, err := getPage(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// it is possible to use templating to create the final html file
	if path.Base(r.URL.Path) == "/html/index.html" {
		//t, _ :=template.ParseFiles("view.html")
		//t.Execute(w, p)
	}
	//fmt.Print("Ext = " + path.Ext(r.URL.Path) + "\n")
	// correct response header content type for css files
	if path.Ext(r.URL.Path) == ".css" {
		w.Header().Add("Content-Type", "text/css")
	}
	fmt.Fprintf(w, "%s", p.Body)
}

/****************************************************************************************/
/*  the handler function queries the incoming http request for pre-allocated URL paths  */
/*  any other paths are served as static files using GET, POST and PUT methods          */
/****************************************************************************************/
func handler(w http.ResponseWriter, r *http.Request) {
	// first check for pre-allocated URL paths and pass to appropriate handler functions
	if r.URL.Path == "/getData" {
		fmt.Printf("getData, %q\n", r.URL.RawQuery)
		getDataHandler(w, r)
		return
	}
	// the http request method is selected and the response/request methods
	// are passed to the relevant handler function
	if r.Method == "GET" {
		fmt.Printf("GET, %q\n", html.EscapeString(r.URL.Path))
		getHandler(w, r)
	} else if r.Method == "POST" {
		fmt.Printf("POST, %q\n", html.EscapeString(r.URL.Path))
	} else if r.Method == "PUT" {
		fmt.Printf("PUT, %q\n", html.EscapeString(r.URL.Path))
	} else {
		http.Error(w, "Invalid request method.\n", 405)
	}
}

/****************************************************************************************/
/*  the main function just has one route to the route handler function                  */
/*  its main task is to setup the go server                                             */
/****************************************************************************************/
func main() {

	// at the moment all requests are passed to the route handler function
	http.HandleFunc("/", handler)
	fmt.Printf("Server started using port : 8000\n")
	// the server is started with a logger to log fatal errors
	// it is handled by the "net/http" go package and listens on port 8000
	// it is a blocking procedure
	log.Fatal(http.ListenAndServe(":8000", nil))
	// when running in a terminal it can be shutdown pressing the Ctrl + C keys together
}
