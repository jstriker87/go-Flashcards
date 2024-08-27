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
    Attempts int
    Completed bool
}

//go:embed templates/*.html
var templatesFS embed.FS
//go:embed static
var staticFS embed.FS

var flashcardCountIndex = 0
var StartingFlashcardCount = 0
var gameStarted = false
var flashcards = []Flashcards{
}
var resultsSlice[]Flashcards
var needRevisionCount = 0

func createResultsSlice(){
    resultsSlice = make([]Flashcards, len(flashcards))
    }

func submitUploadedQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

    file, _, err := r.FormFile("file")
    defer file.Close()
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
        attempts:= 0
        completed:= false
        flashcard := Flashcards{Question: question, Answer: answer, Attempts: attempts, Completed: completed}
		flashcards = append(flashcards, flashcard)
    } 
    StartingFlashcardCount = len(flashcards)
    createResultsSlice()
    http.Redirect(w, r, "/question", http.StatusSeeOther)
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
        var startCardCount = 0
            if needRevisionCount > 0 {
                startCardCount = needRevisionCount
            } else {
                startCardCount = StartingFlashcardCount
            }
        if len(flashcards) > 0 {
            gameStarted = true
        }
        for flashcardCountIndex < len(flashcards) && flashcards[flashcardCountIndex].Completed {
                    flashcardCountIndex ++
                }
        if flashcardCountIndex < len(flashcards) { 
        flashTemplate := parseTemplate("questions.html")
        type gameData struct{
            Flashcard Flashcards
            CardCount int
            StartCardCount int

        }
        theGameData := gameData{ 
    
            Flashcard:  flashcards[flashcardCountIndex],
            CardCount: flashcardCountIndex + 1,
            StartCardCount: startCardCount,
        }
        if err := flashTemplate.Execute(w, theGameData); err != nil {
            log.Println("Error executing template:", err)
        }
    }else {
        http.Redirect(w, r, "/end", http.StatusSeeOther)
    }
}

func startFlashcards (w http.ResponseWriter, r *http.Request) {
        flashTemplate := parseTemplate("index.html")
         type gameData struct{
            FcLength int
            GameHasStarted bool
        }
        theGameData := gameData {
            FcLength: len(flashcards),
            GameHasStarted: gameStarted,
        }
        if err := flashTemplate.Execute(w, theGameData); err != nil {
            log.Println("Error executing template:", err)
        }
    }

func questionNeedsRevision (w http.ResponseWriter, r *http.Request) {
    flashcards[flashcardCountIndex].Attempts+=1
    if flashcardCountIndex <len(flashcards){
        flashcardCountIndex++
    }
    http.Redirect(w, r, "/question", http.StatusSeeOther)
}

func questionOK (w http.ResponseWriter, r *http.Request) {
    if (flashcards[flashcardCountIndex].Attempts == 0){
        flashcards[flashcardCountIndex].Attempts+=1
    }
    flashcards[flashcardCountIndex].Completed = true
    if flashcardCountIndex <len(flashcards){
        flashcardCountIndex++
    }
    http.Redirect(w, r, "/question", http.StatusSeeOther)
    }


func replay (w http.ResponseWriter, r *http.Request) {
    flashcardCountIndex=0
    StartingFlashcardCount = 1
    http.Redirect(w, r, "/question", http.StatusSeeOther)

}

func restart(w http.ResponseWriter, r *http.Request) {
    flashcardCountIndex= 0
    flashcards = nil
    gameStarted = false
    http.Redirect(w, r, "/", http.StatusSeeOther)

}

func clearAndGoToMainMenu (w http.ResponseWriter, r *http.Request) {
    flashcards = []Flashcards{
    }
    flashcardCountIndex=0
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func endFlashcards (w http.ResponseWriter, r *http.Request) {
    needRevisionCount = 0 
    for _, item := range flashcards {
        if item.Completed == false{
            needRevisionCount++
        }
    }
    if needRevisionCount == 0{
        gameStarted = false
        //http.Redirect(w, r, "/", http.StatusSeeOther)
    }
        type gameData struct{
            AllFlashcards []Flashcards
            RevisionCount int

        }
        theGameData := gameData{ 
            RevisionCount: needRevisionCount, 
            AllFlashcards:  flashcards,
        }

        flashTemplate:= parseTemplate("end.html")

        if err := flashTemplate.Execute(w,theGameData); err != nil {
            log.Println("Error executing template:", err)
        }
}

func addQuestions(w http.ResponseWriter, r *http.Request) {
        flashTemplate:= parseTemplate("addquestions.html")
        flashcards = nil
        gameStarted = false
        flashcardCountIndex = 0 
        
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
        StartingFlashcardCount = len(flashcards)
        http.Redirect(w, r, "/question", http.StatusSeeOther)
        createResultsSlice() 
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
    var portAvailable = false
	port := 8000
	portstr := strconv.Itoa(port)
	var l net.Listener
	var err error
	for portAvailable != true {
		l, err = net.Listen("tcp", ":" + portstr)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		if strings.Contains(err.Error(), "in use") {
				port += 1
				portstr = strconv.Itoa(port)
			}
		} else {
			portAvailable = true
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
    fmt.Printf("Starting flashcards at http://localhost:%d \n",port)
    openServerWebpage("http://localhost:" + strconv.Itoa(port))
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
    http.HandleFunc("/replay", replay)
    http.HandleFunc("/restart",restart)
    http.HandleFunc("/submitaddquestions", submitQuestions)
    http.HandleFunc("/addquestions", addQuestions);
    http.HandleFunc("/uploadquestions", uploadQuestions);
    http.HandleFunc("/submituploadquestions", submitUploadedQuestions);
    http.HandleFunc("/mainmenu",clearAndGoToMainMenu)
    http.HandleFunc("/end", endFlashcards)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
