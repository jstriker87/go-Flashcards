package main

import (
    "fmt"
    "net/http"
    "log"
    "html/template"
    "strconv"
    "strings"
    "embed"
    "io/fs"
    "net"
)
type Flashcards struct {
    Question string
    Answer   string
}

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

var flashcardCount = 0
var flashcards = []Flashcards{
}
var version = 1.0   
var runCount = 0
var available = false

func parseTemplate(filename string) *template.Template {
    return template.Must(template.ParseFS(templatesFS, "templates/"+filename))
}

func showAnswer(w http.ResponseWriter, r *http.Request) {

        flashTemplate := parseTemplate("answer.html")
        data := map[string]Flashcards{
            "Flashcard": flashcards[flashcardCount],

            }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)

        }
}

func showQuestion(w http.ResponseWriter, r *http.Request) {
        runCount++
        if flashcardCount < len(flashcards){
        flashTemplate := parseTemplate("questions.html")
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
        flashTemplate := parseTemplate("index.html")
        data := map[string]int{
            "flashcardsnum": len(flashcards),
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func questionNeedsRevision (w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/question", http.StatusSeeOther)
    flashcardCount++

}


func questionOK (w http.ResponseWriter, r *http.Request) {
    flashcards = append(flashcards[:flashcardCount], flashcards[flashcardCount+1:]...) 
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
        

        flashTemplate:= parseTemplate("end.html")
            data := map[string]int{
            "Flashcard": len(flashcards),
        }

        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
}

func preSubmitQuestions(w http.ResponseWriter, r *http.Request) {
        flashTemplate:= parseTemplate("addquestions.html")
        data := map[string]int{
            "Flashcard": 0,
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func submitQuestions(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles("addquestions.html")
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
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func uploadQuestions(w http.ResponseWriter, r *http.Request) {
        flashTemplate:= parseTemplate("uploadquestions.html")
        data := map[string]int{
            "Flashcard": 0,
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }


func checkPort() int {
	port := 8000
	portstr := strconv.Itoa(port)
	var l net.Listener
	var err error
	for available != true {
		l, err = net.Listen("tcp", ":" + portstr)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		if strings.Contains(err.Error(), "in use") {
				port += 1
				portstr = strconv.Itoa(port)
			}
		} else {
			available = true
		}
	}
	defer l.Close()
    return port
}
func main() {
    port:= checkPort()
    if runCount < 1 {

        fmt.Println("Starting flashcards. Go to http://localhost:",port)

    }
    runCount++
    staticSubFS, err := fs.Sub(staticFS, "static")
    if err != nil {
        log.Fatal(err)
    }
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))
    http.HandleFunc("/", startFlashcards)
    http.HandleFunc("/question", showQuestion)
    http.HandleFunc("/needsRevision", questionNeedsRevision)
    http.HandleFunc("/ok", questionOK)
    http.HandleFunc("/answer", showAnswer)
    http.HandleFunc("/restart", restart)
    http.HandleFunc("/submitaddquestions", submitQuestions)
    http.HandleFunc("/addquestions", preSubmitQuestions);
    http.HandleFunc("/uploadquestions", uploadQuestions);
    http.HandleFunc("/end", endFlashcards)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
