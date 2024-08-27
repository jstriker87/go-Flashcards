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
* This program includes executable programs for Linux, Mac & Windows. To run the program it is as simple as double clicking on the program
* If you are compiling the program from source, Golang version 1.23 is required 
### Installing
* Download the executable for the platform of your choice

### Installing from source
1. Clone the repository
2. Navigate to the repository's  directory
3. Depending on your operating system type the command below
    - Windows:
        - GOOS=windows GOARCH=amd64 go build -o build/Flashcards.exe
    - Mac:
        - GOOS=darwin GOARCH=amd64 go build -o build/Flashcards-Mac
    - Linux:
        - GOOS=linux GOARCH=amd64 go build -o build/Flashcards-Linux
4. The Executable will then be created in the current folder
### Running the program
* Note that the program uses port 8000, however if this port is unavailable the program will check the port try the next port (8001), and continue until it finds a free port.
* When starting the program, the terminal window will display the port that has been used (e.g http://localhost:8000)
* For Windows 
    * Double click on the executable program
    * You may receive a warning to advise that the firewall needs permission to run the program. If you receive this message then press accept 
    * A web page will be loaded with the flashcards program
* For Linux / Mac 
    * Double click on the executable program
    * A web page will be loaded with the flashcards program

## Authors
@Jstriker87
## Version History
* 1.2
    * Dark / light mode now follows the users operating systems current light or dark mode
    * Re-factored questions and answers for when questions 'needs revision' so there is a 'completed' value that is set when the user sets the answer to 'OK' and also an 'Attempts' value that counts the number of times that the user has to retry a question before they set  it as 'OK' 
    * Added a summary at the end of the flashcards once the user has no 'needs revision' questions
* 1.2
    * Dark mode added
    * Icons updated and configured for light and dark mode
    * The logic for the 'add questions' and 'restart' buttons have been updated so that the 'add questions' button disappears when the user has added questions, which also then makes the 'restart' button visible
* 1.1    
    * Added progress bar when questions are displayed. Fonts have also been updated to be unified across the program
* 1.0.3.1
    * Fixed mistake causing audio to be played when starting flashcards, and also added titles to a few pages
* 1.0.3
    * Added sounds when selecting if the answer to a question is ok or if it needs revision
* 1.0.2 
    * Added functionality for users to upload a text file (.txt) with questions and answers for the Flashcards
* 1.0.1
    * Added port checker and updated Readme
* 1.0.0
    * Initial Release

## License
This project is licensed under the MIT License.

## Acknowledgements
This project uses code licensed from the sources below
[Google Font Roboto](https://fonts.google.com/specimen/Roboto/about)
[File Upload code by freshman-tech](https://github.com/Freshman-tech/file-upload)
