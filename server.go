package main

import (
    "fmt"
    "net/http"
    "log"
    "html/template"
    "strconv"
    "strings"
    "embed"
    "io"
    "io/fs"
    "net"
    "bufio"
    "os/exec"
    "runtime"
)
type Flashcards struct {
    Question string
    Answer   string
}

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

var flashcardCountIndex = 0
var flashcardCount = 1
var StartingFlashcardCount = 0


var flashcards = []Flashcards{
}
var version = 1.0   
var runCount = 0
var available = false
const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

    file, fileHeader, err := r.FormFile("file")
    defer file.Close()
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

    scanner := bufio.NewScanner(file)
    var lines []string
    for scanner.Scan() {
        text:= scanner.Text()
        lines = append(lines,text)
    }
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

    for i:=0;i<len(lines)-1;i+=2{
        question:=lines[i]
        answer:=lines[i+1]
        flashcard := Flashcards{Question: question, Answer: answer}
		flashcards = append(flashcards, flashcard)
    } 
    http.Redirect(w, r, "/", http.StatusSeeOther)
	}

func parseTemplate(filename string) *template.Template {
    return template.Must(template.ParseFS(templatesFS, "templates/"+filename))
}

func showAnswer(w http.ResponseWriter, r *http.Request) {

        flashTemplate := parseTemplate("answer.html")
        data := map[string]Flashcards{
            "Flashcard": flashcards[flashcardCountIndex],

            }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)

        }
}

func showQuestion(w http.ResponseWriter, r *http.Request) {
        runCount++
        if flashcardCountIndex < len(flashcards){
        flashTemplate := parseTemplate("questions.html")
        type gameData struct{
            Flashcard Flashcards
            CardCount int
            StartCardCount int

        }
        theGameData := gameData{ 
    
            Flashcard:  flashcards[flashcardCountIndex],
            CardCount: flashcardCount,
            StartCardCount: StartingFlashcardCount,
        }
        if err := flashTemplate.Execute(w, theGameData); err != nil {
            log.Println("Error executing template:", err)
        }
    }else {
        http.Redirect(w, r, "/end", http.StatusSeeOther)
    }
}

func startFlashcards (w http.ResponseWriter, r *http.Request) {
        StartingFlashcardCount = len(flashcards)
        flashTemplate := parseTemplate("index.html")
        data := map[string]int{
            "flashcardsnum": len(flashcards),
        }
        if err := flashTemplate.Execute(w, data); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func questionNeedsRevision (w http.ResponseWriter, r *http.Request) {
    flashcardCount++
    flashcardCountIndex++
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}


func questionOK (w http.ResponseWriter, r *http.Request) {
    flashcardCount++
    flashcards = append(flashcards[:flashcardCountIndex], flashcards[flashcardCountIndex+1:]...) 
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}


func restart (w http.ResponseWriter, r *http.Request) {
    flashcardCountIndex=0
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}


func clearAndGoToMainMenu (w http.ResponseWriter, r *http.Request) {
    flashcards = []Flashcards{
    }
    flashcardCountIndex=0
    http.Redirect(w, r, "/", http.StatusSeeOther)
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
        w.Header().Add("Content-Type", "text/html")
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

func main() {
    port:= checkPort()
    if runCount < 1 {

        fmt.Printf("Starting flashcards at http://localhost:%d",port)
        openServerWebpage("http://localhost:" + strconv.Itoa(port))

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
    http.HandleFunc("/submituploadquestions", uploadHandler);
    http.HandleFunc("/mainmenu",clearAndGoToMainMenu)
    http.HandleFunc("/end", endFlashcards)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
