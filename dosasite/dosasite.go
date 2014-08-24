package main

import(
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path"
    )
    
func main(){
    fs := http.FileServer(http.Dir("public"))
    http.Handle("/public/", http.StripPrefix("/public/",fs))
    
    http.HandleFunc("/", ServeTemplate)
    
    fmt.Println(fs)

    
    fmt.Println("listening ..")
    err := http.ListenAndServe(GetPort(), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
        return
    }
}


// Get the Port from the environment so we can run on Heroku
func GetPort() string {
    var port = os.Getenv("PORT")
    //set a default port if there is nothing in the environment
    if port == "" {
        port = "4747"
        fmt.Println("INFO: No PORT environemnt variable detected, default to " + port)
    }
    return ":" + port
}

func ServeTemplate(w http.ResponseWriter, r *http.Request){
    lp := path.Join("templates", "layout.html")
    fp := path.Join("templates", r.URL.Path)
    
    fmt.Println("fp is " + r.URL.Path)
    
    //Return a 404 if the template doesnt exist
    info, err := os.Stat(fp)
    if err != nil {
        if os.IsNotExist(err){
            fmt.Println("did not find template ", fp)
            http.NotFound(w,r)
            return
        }
    }
    
    if info.IsDir(){
        fmt.Println("did not find directory ")
        http.NotFound(w,r)
        return
    }
    
    templates, err :=template.ParseFiles(lp, fp)
    if err != nil{
        fmt.Println(err)
        http.Error(w, "500 Internal Server Error",500)
        return
    }
    
    fmt.Println(w)
    templates.ExecuteTemplate(w, "layout", nil)
    
}