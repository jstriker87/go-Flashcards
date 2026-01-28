package mobileapp

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



// Flashcards struct. This is used as the key structure for each flashcard
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


// Set the index of of which flashcard is being used to zero
var flashcardCountIndex = 0
// Set 'gameStarted' to false to indicate that a game has not yet begun
var gameStarted = false

// Stores all initial flashcards
var flashcards = []Flashcards{}

// Stores the flashcards that have been marked as 'ok' once all questions have been completed
var doneFlashcards = []Flashcards{}

// Sets counter for the number of flashcards that are marked as 'needs revision'
var needRevisionCount = 0


// This function is used when the user uploads questions using the upload file page
func SubmitUploadedQuestions(w http.ResponseWriter, r *http.Request) {
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

	// Once complete re-direct the user to the 'question' page to start the user working on the flashcards
	http.Redirect(w, r, "/question", http.StatusSeeOther)
}


// Helper to help process each template. Accepts the filename as an input and returns a pointer to a 'Template' with the desired filename in the 'templates' folder
func ParseTemplate(filename string) *template.Template {
	return template.Must(template.ParseFS(templatesFS, "templates/"+filename))
}


// This function shows the answer to the question 
func ShowAnswer(w http.ResponseWriter, r *http.Request) {

	// Load and parse the template file end.html, convert it to a template object, and store it in flashTemplate so I can execute it later”
	flashTemplate := ParseTemplate("answer.html")


	// Create a map using the current 'flashcardCountIndex' to be used to show the answer
	data := map[string]Flashcards{
		"Flashcard": flashcards[flashcardCountIndex],
	}
	if err := flashTemplate.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)

	}
}


func ShowQuestion(w http.ResponseWriter, r *http.Request) {

	// Set StartCount to the length of flashcards
	var startCardCount = len(flashcards)
	// If the length of flashcards is greater than zero then set the 'gameStarted' variable to true
	if len(flashcards) > 0 {
		gameStarted = true
	}
	// If the current index of flashcards is less than the length of flashcards and the flashcard has been completed then increment the flashcardCountIndex to skip that flasshcard
	for flashcardCountIndex < len(flashcards) && flashcards[flashcardCountIndex].Completed {
		flashcardCountIndex++
	}


	// Load and parse the template file questions.html, convert it to a template object, and store it in flashTemplate so I can execute it later”
	if flashcardCountIndex < len(flashcards) {
		flashTemplate := ParseTemplate("questions.html")
		// Create a struct with the flashcard, the card count (the flashcardCountIndex +1) and the starting CardCount (length of flashcards)
		type gameData struct {
			Flashcard      Flashcards
			CardCount      int
			StartCardCount int
		}

		// Create an instance of the struct above and name it 'theGameData'
		theGameData := gameData{

			Flashcard:      flashcards[flashcardCountIndex],
			CardCount:      flashcardCountIndex + 1,
			StartCardCount: startCardCount,
		}

		// Execute the templatand check that no errors occurs
		if err := flashTemplate.Execute(w, theGameData); err != nil {
			log.Println("Error executing template:", err)
		}

	// If there are no flashcards (i.e. the length of flashcards is zero), redirect the user to the end page
	} else {
		http.Redirect(w, r, "/end", http.StatusSeeOther)
	}
}

// This function is for the homepage index.html 

func StartFlashcards(w http.ResponseWriter, r *http.Request) {

	// Load and parse the template file questions.html, convert it to a template object, and store it in flashTemplate so I can execute it later”
	flashTemplate := ParseTemplate("index.html")

	// Execute the template with no data (as the variables needed are already global)
	if err := flashTemplate.Execute(w, nil); err != nil {
		log.Println("Error executing template:", err)
	}
}


// This function is used for when the user selects that a question 'needs revision'
func QuestionNeedsRevision(w http.ResponseWriter, r *http.Request) {
	// Add one to the 'Attempts' value of the current flashcard 
	flashcards[flashcardCountIndex].Attempts += 1
	// If the flashcardCountIndex is less than the length of flashcards array then increment the flashcardCountIndex by one (i.e. Go to the next flashcard)
	if flashcardCountIndex < len(flashcards) {
		flashcardCountIndex++
	}
	// Re-direct the user back to the 'question' page to load the next question
	http.Redirect(w, r, "/question", http.StatusSeeOther)
}


// This function is used for when the user selects that a question is 'ok'
func QuestionOK(w http.ResponseWriter, r *http.Request) {

	// Add one to the 'Attempts' value of the current flashcard 
	flashcards[flashcardCountIndex].Attempts += 1
	// Set the 'Completed' value of the current flashcard to true
	flashcards[flashcardCountIndex].Completed = true

	// If the flashcardCountIndex is less than the length of flashcards array then increment the flashcardCountIndex by one (i.e. Go to the next flashcard)
	if flashcardCountIndex < len(flashcards) {
		flashcardCountIndex++
	}

	// Re-direct the user back to the 'question' page to load the next question
	http.Redirect(w, r, "/question", http.StatusSeeOther)
}


// This function is used when the user has one or more questions that were marked as 'needs revision' and the user presses the 'replay' button to go through them again
func Replay(w http.ResponseWriter, r *http.Request) {

	// Iterate over each flashcard
	for i := len(flashcards) - 1; i >= 0; i-- {
		// If the current flashcards 'Completed' value is true
		// Add the flashcard to the 'doneFlashcards' array of Flashcards
		// Then remove that flashcard from the 'flashcards' array of flashcards
		if flashcards[i].Completed {
			doneFlashcards = append(doneFlashcards, flashcards[i])
			flashcards = append(flashcards[:i], flashcards[i+1:]...)
		}
	}

	// Set the current index of the flashcard being processed back to zero
	flashcardCountIndex = 0

	// Re-direct the user back to the 'question' page to load the next question
	http.Redirect(w, r, "/question", http.StatusSeeOther)

}

// This function is used when restarting the game to reset variables 
func Restart(w http.ResponseWriter, r *http.Request) {
	// Set flashcardCountIndex back to zero
	flashcardCountIndex = 0
	// Set flashcards back to nil
	flashcards = nil
	// Set 'gameStarted' back to false
	gameStarted = false
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// This function is used when a game has been fully completed and the user starts a new game
func ClearAndGoToMainMenu(w http.ResponseWriter, r *http.Request) {
	// Set the 'flashcards' and 'doneFlashcards' array back to empty
	flashcards = []Flashcards{}
	doneFlashcards =  []Flashcards{}
	// Set 'flashcardCountIndex' back to zero so any new round will start at the beginning
	flashcardCountIndex = 0
	// Re-redirect the user to the root page (which calls the startFlashcards function)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


// This function is for when a user completes a round
func EndFlashcards(w http.ResponseWriter, r *http.Request) {
	// Initially set 'needRevisionCount' (the number of flashcards that need to be reviewed) to zero
	needRevisionCount = 0

	// Create a struct to store the flashcards array and 'RevisionCount' (the number of flashcards that need revision)
	type gameData struct {
		AllFlashcards []Flashcards
		RevisionCount int
	}

	// Iterate over each flashcard. If the 'Completed' status of the flashcard is false then increment the 'needRevisionCount' value by one
	for _, item := range flashcards {
		if item.Completed == false {
			needRevisionCount++
		}
	}

	// If 'needRevisionCount' is zero (i.e. there are no flashcards that needs revision) then combine the 'flashcards' and 'doneFlashcards' into the 'flashcards' array
	// This is so that the final statistics page that displays each question and answer, and the number of attempts taken to answer the question, contains all the questions and answers, not just the ones that needed revision
	if needRevisionCount == 0 {

		flashcards = append(flashcards, doneFlashcards...)
		gameStarted = false

	}

	// Create an instance of the 'gameData' struct named 'theGameData' with the data
	theGameData := gameData{
		RevisionCount: needRevisionCount,
		AllFlashcards: flashcards,
	}


	// Load and parse the template file end.html, convert it to a template object, and store it in flashTemplate so I can execute it later”

	flashTemplate := ParseTemplate("end.html")

	// Execute the template using 'gameData'

	if err := flashTemplate.Execute(w, theGameData); err != nil {
		log.Println("Error executing template:", err)
	}
}

//  This function is used to set-up the process to add questions
func AddQuestions(w http.ResponseWriter, r *http.Request) {

	// Load and parse the template file addquestions.html, convert it to a template object, and store it in flashTemplate so I can execute it later”
	flashTemplate := ParseTemplate("addquestions.html")
	// Set main variables used in game logic to their initial values
	flashcards = nil
	gameStarted = false
	flashcardCountIndex = 0

	// Execute the template with no data
	if err := flashTemplate.Execute(w, nil); err != nil {
		log.Println("Error executing template:", err)
	}
}



// This function is used for processing questions and answers when adding them manually on the 'addquestions' page (not when uploading a text file)

func SubmitQuestions(w http.ResponseWriter, r *http.Request) {

	// Check that the method used is GET and if so load and parse the template file 'addquestions".html and execute it (In effect this reloads the page)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("addquestions.html")
		t.Execute(w, nil)
	// IF the method is not GET (i.e. the method is POST)
	} else {
		// Parse the form
		r.ParseForm()
		// Iterate over the questions (the maximum provided questions and answers is 20)

		for i := 1; i < 21; i++ {
			// Get each question and answer by checking the value along with the number of questions and answers
			question := r.FormValue("question" + strconv.Itoa(i))
			answer := r.FormValue("answer" + strconv.Itoa(i))

			// If either the question or answer is empty then skip it
			if question == "" || answer == "" {
				continue
			}
			// Create a new flashcard with the question and answer of that item
			flashcard := Flashcards{Question: question, Answer: answer}
			// Add the flashcard to the 'flashcards' array
			flashcards = append(flashcards, flashcard)
		}
		// Re-direct the user to the 'questions' page
		http.Redirect(w, r, "/question", http.StatusSeeOther)
	}
}


// This function is used for preparing to uploading questions
func UploadQuestions(w http.ResponseWriter, r *http.Request) {
	// Set the content type to header text/html as the user will be providing a text file for the questions and answers
	w.Header().Add("Content-Type", "text/html")
	// Load and parse the template file uploaduestions.html, convert it to a template object, and store it in flashTemplate so I can execute it later”
	flashTemplate := ParseTemplate("uploadquestions.html")
	if err := flashTemplate.Execute(w, nil); err != nil {
		log.Println("Error executing template:", err)
	}
}


// This function is used to check what port the program should use. The default port is 8000, but if it is not available it increments the port number by one until it finds one that is available
func CheckPort() int {
	var portAvailable = false
	port := 8000
	portstr := strconv.Itoa(port)
	// Create a new listener to check ports 
	var l net.Listener
	var err error
	// Check current number number to see if it is free
	for portAvailable != true {
		l, err = net.Listen("tcp", ":"+portstr)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			// If the error contains 'in use' the port is not available so increase the port number by one
			if strings.Contains(err.Error(), "in use") {
				port += 1
				portstr = strconv.Itoa(port)
			}
		// Once a port is found to be free set 'portAvailable' to true to stop the search for  a free port
		} else {
			portAvailable = true
		}
	}
	// Defer the closing of the listener so it doesn't stop prematurely
	defer l.Close()
	// Return the found port number
	return port
}

// The function is used to automatically open the web page for the flashcards 
func OpenServerWebpage(url string) error {
	// Set empty cmd string and create a string array calls 'args'
	var cmd string
	var args []string
	// runtime.GOOS is the architecture of the system
	switch runtime.GOOS {
	// If the user is using Windows 
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	// If the user is using Mac
	case "darwin":
		cmd = "open"
	// If the user is using linux
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	// Add the url to args
	args = append(args, url)
	// Return an execution command to open the browser with the flashcards url
	return exec.Command(cmd, args...).Start()
}


// The main function is run a program launch

func StartFlashcardsServer() {
	// Get the available port
	port := CheckPort()
	fmt.Printf("Starting flashcards at http://localhost:%d \n", port)
	// Open the flashcards web page
	OpenServerWebpage("http://localhost:" + strconv.Itoa(port))
	// Create a sub-filesystem in the 'static directory' inside staticFS (the file system that you start from
	//This means that when the program runs it makes the 'static' folder the root of the program
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	// For the http.Handle line below
		//Registers the handler for all URLs starting with /static/ but removes the 'static' from any url as we are already serving the file server from the 'static' directory
		// Create a file server that serves files from the 'staticSubFS' directory (the 'static' directory)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))


	// Below is a list of paths that the user can go to and which function it calls. Each function is one in this server, and relates to one or more page. These pages are called from the 
	// html files in this program in the 'templates' folder
	http.HandleFunc("/", StartFlashcards)
	http.HandleFunc("/question", ShowQuestion)
	http.HandleFunc("/needsRevision", QuestionNeedsRevision)
	http.HandleFunc("/ok", QuestionOK)
	http.HandleFunc("/answer", ShowAnswer)
	http.HandleFunc("/replay", Replay)
	http.HandleFunc("/restart", Restart)
	http.HandleFunc("/submitaddquestions", SubmitQuestions)
	http.HandleFunc("/addquestions", AddQuestions)
	http.HandleFunc("/uploadquestions", UploadQuestions)
	http.HandleFunc("/submituploadquestions", SubmitUploadedQuestions)
	http.HandleFunc("/mainmenu", ClearAndGoToMainMenu)
	http.HandleFunc("/end", EndFlashcards)
	// Start a web server using the available port
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

