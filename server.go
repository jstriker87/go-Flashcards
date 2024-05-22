package main 


import ( 
    "fmt"
	"html/template"
	"log"
	"net/http"
)

type Flashcards struct {
    Question string
    Answer   string
}

var flashcardCount = 0
var flashcardOk = true
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
        flashcards = append(flashcards[:flashcardCount], flashcards[flashcardCount+1:]...)
        fmt.Println(flashcards)
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

func updateFlashcards (w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/results", http.StatusSeeOther)
    }
func resetFlashcards (w http.ResponseWriter, r *http.Request) {
    fmt.Println("beep")
    flashcardCount = 0
    http.Redirect(w, r, "/", http.StatusSeeOther)
    }


func endFlashcards (w http.ResponseWriter, r *http.Request) {
        fmt.Printf("The flashcard count is %d \n",flashcardCount)
        flashTemplate := template.Must(template.ParseFiles("end.html"))
            data := map[string]int{
            "Flashcard": len(flashcards),
        }
        
        flashcardCount = 0
        fmt.Printf("The updated flashcard count is %d \n",flashcardCount)
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func main() {

    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fileServer))
    http.HandleFunc("/", checkFlashcards)
    http.HandleFunc("/results", flashcardResult)
    http.HandleFunc("/update", updateFlashcards)
    http.HandleFunc("/restart", resetFlashcards)
    http.HandleFunc("/end", endFlashcards)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
