# Flashcard Revision Program

## Description
* Flashcard program written in Go that uses a web server
* The program offers a very simple web Interface that allows users to add questions and then test knowledge
* After each questions is displayed, and are shown the answer followed by an option to either confirm that your knowledge is ok, or that you need further revision
* If you select 'ok' then the question is removed from the list of questions
* If you select 'needs revision' then the question is not removed from the list of questions
* At the end of the flashcards if you have any questions that you selected as 'needs revision' then you will be given the option to restart with the remaining questions 
## Getting Started

### Dependencies
* This program includes executable programs for Linux, Mac & Windows
* If you are compiling the program from source, Golang version 1.23 is required 
* Note that the program uses port 8000, so it must be free to be used. Future versions will check the port try the next port (8001), and continue until it finds a free port 

### Installing
* Download the executable for the platform of your choice

### Installing from source
1. Clone the repository
2. Navigate to the repository's  directory
3. In your command prompt or terminal type `go build Flashcards`
4. The Executable will then be created in the current folder

### Executing program
* For Windows 
    * Run the executable 
    * Open your web browser and go to http://localhost:8000
    * You may receive a warning to advise that the firewall needs permission to run the program. If you receive this message then press accept 
* For Linux / Mac 
    * Navigate to the folder where the executable is located
    * Execute the program by typing `./Flashcards`
    * The webpage http://localhost:8000 should automatically be displayed, but if it doesn't then go to the webpage 
## Authors
@Jstriker87

## Version History
* 0.1
    * Initial Release

## License
This project is licensed under the MIT License.
