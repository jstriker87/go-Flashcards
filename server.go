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
var needsRevision = false
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
        if needsRevision {

            flashcards = append(flashcards[:flashcardCount], flashcards[flashcardCount+1:]...)
            needsRevision = false

        }else{

            flashcardCount++
        }

        }else{
            http.Redirect(w, r, "/end", http.StatusSeeOther)           
            }
        }
func checkFlashcards (w http.ResponseWriter, r *http.Request) {
        if flashcardCount < len(flashcards){
        flashTemplate := template.Must(template.ParseFiles("questions.html"))
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

func startFlashcards (w http.ResponseWriter, r *http.Request) {

        flashTemplate := template.Must(template.ParseFiles("index.html"))
        data := map[string]int{
            "Flashcard": 0,
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func updateFlashcards (w http.ResponseWriter, r *http.Request) {
    needsRevision = true
    http.Redirect(w, r, "/results", http.StatusSeeOther)

}
func restartFlashcards (w http.ResponseWriter, r *http.Request) {
    needsRevision = false 
    flashcardCount = 0
    http.Redirect(w, r, "/start", http.StatusSeeOther)

}

func endFlashcards (w http.ResponseWriter, r *http.Request) {
    if flashcardCount == 0{
        http.Redirect(w, r, "/", http.StatusSeeOther)


    }
        flashTemplate := template.Must(template.ParseFiles("end.html"))
            data := map[string]int{
            "Flashcard": len(flashcards),
        }

        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func main() {

    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fileServer))
    http.HandleFunc("/", startFlashcards)
    http.HandleFunc("/start", checkFlashcards)
    http.HandleFunc("/update", updateFlashcards)
    http.HandleFunc("/results", flashcardResult)
    http.HandleFunc("/end", endFlashcards)
    http.HandleFunc("/restart", restartFlashcards)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
