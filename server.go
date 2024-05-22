package main

import (
	"log"
	"net/http"
	"html/template"
)

type Flashcards struct {
    Question string
    Answer   string
}

var flashcardCount = 0
var flashcards = []Flashcards{
    {"What is the capital of England", "London"},
    {"What is the capital of Scotland", "Edinburgh"},
}

func flashcardResult (w http.ResponseWriter, r *http.Request) {
            if flashcardCount < len(flashcards){
            flashTemplate := template.Must(template.ParseFiles("results.html"))
            data := map[string]Flashcards{
                "Flashcard": flashcards[flashcardCount],
            }
            if err := flashTemplate.Execute(w, data); err != nil {
                log.Println("Error executing template:", err)
            }
        if r.URL.Path != "/favicon.ico" { 
            flashcardCount++
        }
        }else{
            http.Redirect(w, r, "/end", http.StatusSeeOther)           
        }
}
func checkFlashcards (w http.ResponseWriter, r *http.Request) {
        if flashcardCount < len(flashcards){
        flashTemplate := template.Must(template.ParseFiles("index.html"))
        data := map[string]Flashcards{
            "Flashcard": flashcards[flashcardCount],
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }else {
        http.Redirect(w, r, "/end", http.StatusSeeOther)
    }
}


func endFlashcards (w http.ResponseWriter, r *http.Request) {
        flashTemplate := template.Must(template.ParseFiles("end.html"))
        data:=-1
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func main() {

    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fileServer))
    http.HandleFunc("/", checkFlashcards)
    http.HandleFunc("/results", flashcardResult)
    http.HandleFunc("/end", endFlashcards)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
