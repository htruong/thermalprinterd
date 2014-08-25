package main 

import (
  "net/http"
  "fmt"
  "text/template"
  "flag"
  "os/exec"
  "strings"
)

var dev = flag.String("device", "/dev/ttyO1", "Output printer device")
var port = flag.Int("port", 8080, "Port to listen to")
var cmdline = flag.String("exec", "/usr/local/bin/serialprinter", "Command to execute")


func main() {
  http.HandleFunc("/dest", posthandler)
  http.HandleFunc("/", formhandler)
  http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func formhandler (w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, `<!DOCTYPE html>
<html>
  <body>
    <form method="POST" action="/dest">
      <div>
      <textarea name="content" placeholder="Something to print"></textarea>
      </div>
      <button>Submit</button>
    </form>
  </body>
</html>`)
}

func posthandler (w http.ResponseWriter, r *http.Request) {
  stringToPrint := r.FormValue("content")
  cmd := exec.Command(*cmdline, "-s", *dev)
  cmd.Stdin = strings.NewReader(stringToPrint)
  printerErr := cmd.Run()

  if printerErr != nil {
    if err := postTemplate.Execute(w, stringToPrint); err != nil {
      http.Error(w, "Template error: ", http.StatusInternalServerError)
    }
  }
}

var postTemplate = template.Must(template.New("PostTemplate").Parse(postTemplateHTML)) 
const postTemplateHTML = `
<!DOCTYPE html>
<html>
  <body>
    <h3>Printed successfully</h3>
    <pre>{{.|html}}</pre>
  </body>
</html>`
