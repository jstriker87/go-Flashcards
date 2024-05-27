package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

type Flashcards struct {
    Question string
    Answer   string
}

var flashcardCount = 0
var flashcards = []Flashcards{
}

func showAnswer(w http.ResponseWriter, r *http.Request) {
        flashTemplate := template.Must(template.ParseFiles("answer.html"))
        data := map[string]Flashcards{
            "Flashcard": flashcards[flashcardCount],

            }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)

        }
}

func showQuestion(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("The flashcard count is %d \n",len(flashcards))
        fmt.Printf("The flashcardCount  count is %d \n",flashcardCount)
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

func questionNeedsRevision (w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}


func questionOK (w http.ResponseWriter, r *http.Request) {
    flashcards = append(flashcards[:flashcardCount], flashcards[flashcardCount+1:]...) 
    flashcardCount++
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}


func restart (w http.ResponseWriter, r *http.Request) {
    flashcardCount=0
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}

func endFlashcards (w http.ResponseWriter, r *http.Request) {
    if len(flashcards) == 0{
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

func openServerWebpage(url string) error {
    var cmd string
    var args []string
    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}

func preSubmitQuestions(w http.ResponseWriter, r *http.Request) {

        flashTemplate := template.Must(template.ParseFiles("add.html"))
        data := map[string]int{
            "Flashcard": 0,
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }


func submitQuestions(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles("add.html")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        for i := 1; i < 10; i++ {
            question:=r.FormValue("question"+strconv.Itoa(i))
            answer:=r.FormValue("answer"+strconv.Itoa(i))
            if question == "" || answer == "" {
			    continue
		    }
            flashcard := Flashcards{Question: question, Answer: answer}
		    flashcards = append(flashcards, flashcard)
        }
        fmt.Println(flashcards) 
    }
}

func main() {

    openServerWebpage("http://localhost:8000")
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fileServer))
    http.HandleFunc("/", startFlashcards)
    http.HandleFunc("/question", showQuestion)
    http.HandleFunc("/needsRevision", questionNeedsRevision)
    http.HandleFunc("/ok", questionOK)
    http.HandleFunc("/answer", showAnswer)
    http.HandleFunc("/restart", restart)
        http.HandleFunc("/submitaddquestions", submitQuestions)
    http.HandleFunc("/addquestions", preSubmitQuestions);
    http.HandleFunc("/end", endFlashcards)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
