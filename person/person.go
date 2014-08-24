package main

import(
    "log"
    "os"
    "text/template"
    
    )
    
type Person struct {
    Name string
    Emails []string
    }
    
const tmpl = 
    `{{$name := .Name}}
    {{range .Emails}}
        Name is {{$name}}, email id is {{.}}
    {{end}}`
    
func main(){
    person := Person{
        Name: "Satish",
        Emails: []string{"satish@rubylearning.org", "satishtalim@gmail.com"},
    }
    
    t := template.New("Person template")
    
    t,err := t.Parse(tmpl)
    if err != nil {
        log.Fatal("Parse", err)
    }
        
    err = t.Execute(os.Stdout, person)
    if err != nil {
        log.Fatal("Execute", err)
        return
    }
}