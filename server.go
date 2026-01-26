package main

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Flashcards struct {
	Question  string
	Answer    string
	Attempts  int
	Completed bool
}

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

var flashcardCountIndex = 0
var StartingFlashcardCount = 0
var gameStarted = false

// Stores all initial flashcards
var flashcards = []Flashcards{}

// Stores the flashcards that have been marked as 'ok' once all questions have been completed
var doneFlashcards = []Flashcards{}

// Sets counter for the number of flashcards that are marked as 'needs revision'
var needRevisionCount = 0


// This function is used when the user uploads questions using the upload file pagee
func submitUploadedQuestions(w http.ResponseWriter, r *http.Request) {
	// Check that the request type is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the 'file' from the POST rquest
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create a new scanner to iterate through each line of the provided text file
	scanner := bufio.NewScanner(file)
	// Store in a string array and append each line of the array
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		lines = append(lines, text)
	}

	// Go back to the start of the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Iterate over each item question and answer the 'lines' array
	// Note that the for loop iterates over two items each time as each question and answer is two lines
	for i := 0; i < len(lines)-1; i += 2 {
		question := lines[i]
		answer := lines[i+1]
		// Create an instance of the flashcards struct and add the question, and answer from the two lines and set the number of attempts to zero, and 'Completed' to false
		flashcard := Flashcards{Question: question, Answer: answer, Attempts: 0, Completed: false}
		// Append the flashcard to the flashcards array 
		flashcards = append(flashcards, flashcard)
	}
	// Set 'StartingFlashcardCount' to the length of the flashcards
	StartingFlashcardCount = len(flashcards)
	// Once complete re-direct the user to the 'question' page to start the user working on the flashcards
	http.Redirect(w, r, "/question", http.StatusSeeOther)
}


// Helper to help process each template. Accepts the filename as an input and returns a pointer to a 'Template' with the desired filename in the 'templates' folder
func parseTemplate(filename string) *template.Template {
	return template.Must(template.ParseFS(templatesFS, "templates/"+filename))
}


// This function shows the answer to the question 
func showAnswer(w http.ResponseWriter, r *http.Request) {

	// Parse template for the answer
	flashTemplate := parseTemplate("answer.html")


	// Create a map using the current 'flashcardCountIndex' to be used to show the answer
	data := map[string]Flashcards{
		"Flashcard": flashcards[flashcardCountIndex],
	}
	if err := flashTemplate.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)

	}
}


func showQuestion(w http.ResponseWriter, r *http.Request) {
	var startCardCount = 0
	startCardCount = len(flashcards)
	if len(flashcards) > 0 {
		gameStarted = true
	}
	for flashcardCountIndex < len(flashcards) && flashcards[flashcardCountIndex].Completed {
		flashcardCountIndex++
	}
	if flashcardCountIndex < len(flashcards) {
		flashTemplate := parseTemplate("questions.html")
		type gameData struct {
			Flashcard      Flashcards
			CardCount      int
			StartCardCount int
		}
		theGameData := gameData{

			Flashcard:      flashcards[flashcardCountIndex],
			CardCount:      flashcardCountIndex + 1,
			StartCardCount: startCardCount,
		}

		if err := flashTemplate.Execute(w, theGameData); err != nil {
			log.Println("Error executing template:", err)
		}
	} else {
		http.Redirect(w, r, "/end", http.StatusSeeOther)
	}
}

func startFlashcards(w http.ResponseWriter, r *http.Request) {
	flashTemplate := parseTemplate("index.html")
	type gameData struct {
		FcLength       int
		GameHasStarted bool
	}
	theGameData := gameData{
		FcLength:       len(flashcards),
		GameHasStarted: gameStarted,
	}
	if err := flashTemplate.Execute(w, theGameData); err != nil {
		log.Println("Error executing template:", err)
	}
}

func questionNeedsRevision(w http.ResponseWriter, r *http.Request) {
	flashcards[flashcardCountIndex].Attempts += 1
	if flashcardCountIndex < len(flashcards) {
		flashcardCountIndex++
	}
	http.Redirect(w, r, "/question", http.StatusSeeOther)
}

func questionOK(w http.ResponseWriter, r *http.Request) {

	flashcards[flashcardCountIndex].Attempts += 1
	flashcards[flashcardCountIndex].Completed = true
	if flashcardCountIndex < len(flashcards) {
		flashcardCountIndex++
	}
	http.Redirect(w, r, "/question", http.StatusSeeOther)
}

func replay(w http.ResponseWriter, r *http.Request) {


	for i := len(flashcards) - 1; i >= 0; i-- {
		if flashcards[i].Completed {
			doneFlashcards = append(doneFlashcards, flashcards[i])
			flashcards = append(flashcards[:i], flashcards[i+1:]...)
		}
	}

	flashcardCountIndex = 0
	StartingFlashcardCount = 1
	http.Redirect(w, r, "/question", http.StatusSeeOther)

}

func restart(w http.ResponseWriter, r *http.Request) {
	flashcardCountIndex = 0
	flashcards = nil
	gameStarted = false
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func clearAndGoToMainMenu(w http.ResponseWriter, r *http.Request) {
	flashcards = []Flashcards{}
	flashcardCountIndex = 0
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func endFlashcards(w http.ResponseWriter, r *http.Request) {
	needRevisionCount = 0

	type gameData struct {
		AllFlashcards []Flashcards
		RevisionCount int
	}
	for _, item := range flashcards {
		if item.Completed == false {
			needRevisionCount++
		}
	}

	if needRevisionCount == 0 {

		flashcards = append(flashcards, doneFlashcards...)
		gameStarted = false

	}
	theGameData := gameData{
		RevisionCount: needRevisionCount,
		AllFlashcards: flashcards,
	}

	flashTemplate := parseTemplate("end.html")

	if err := flashTemplate.Execute(w, theGameData); err != nil {
		log.Println("Error executing template:", err)
	}
}

func addQuestions(w http.ResponseWriter, r *http.Request) {
	flashTemplate := parseTemplate("addquestions.html")
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
			question := r.FormValue("question" + strconv.Itoa(i))
			answer := r.FormValue("answer" + strconv.Itoa(i))
			if question == "" || answer == "" {
				continue
			}
			flashcard := Flashcards{Question: question, Answer: answer}
			flashcards = append(flashcards, flashcard)
		}
		StartingFlashcardCount = len(flashcards)
		http.Redirect(w, r, "/question", http.StatusSeeOther)
	}
}

func uploadQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	flashTemplate := parseTemplate("uploadquestions.html")
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
		l, err = net.Listen("tcp", ":"+portstr)
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
	port := checkPort()
	fmt.Printf("Starting flashcards at http://localhost:%d \n", port)
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
	http.HandleFunc("/restart", restart)
	http.HandleFunc("/submitaddquestions", submitQuestions)
	http.HandleFunc("/addquestions", addQuestions)
	http.HandleFunc("/uploadquestions", uploadQuestions)
	http.HandleFunc("/submituploadquestions", submitUploadedQuestions)
	http.HandleFunc("/mainmenu", clearAndGoToMainMenu)
	http.HandleFunc("/end", endFlashcards)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
