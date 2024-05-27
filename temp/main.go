package main

import (
    "fmt"
    "html/template"
    "net/http"
    "log"
    "strconv"
)

func submitQuestions(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Submitting")
    if r.Method == "GET" {
        t, _ := template.ParseFiles("add.html")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        for i := 1; i <= 10; i++ {
            question:=r.FormValue("question"+strconv.Itoa(i))
            answer:=r.FormValue("question"+strconv.Itoa(i))
            if question == "" || answer == "" {
			continue
		    }
            // Add logic to submit question 
        }
        
    }
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


func main() {
    http.HandleFunc("/", preSubmitQuestions);
    http.HandleFunc("/submitquestion", submitQuestions);
    http.ListenAndServe(":8000", nil)
}

