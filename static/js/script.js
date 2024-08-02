let theme = localStorage.getItem("theme")
if (theme == null) {
    localStorage.setItem("theme","light");
}
function goPlay(musicurl,redirecturl){
    console.log("The dark mode status is: " + theme); 
    if (localStorage.getItem("toggleVolume") == "true") {
        currentAudio = new Audio(musicurl);
        currentAudio.play();
        sleep(700).then(() => { window.location.href = redirecturl; }); 
    }
    else {
        sleep(700).then(() => { window.location.href = redirecturl; }); 

    }

}

function noPlay(){
    window.speechSynthesis.cancel();
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function volume() {
    if (localStorage.getItem("toggleVolume") == "false") {
        localStorage.setItem("toggleVolume", "true");
        document.getElementById("volume-icon").src = `/static/images/volume-on-${theme}.png`;
        document.getElementById("volume-icon").title= "Volume On";
        tts();
        return;
    } else {
        localStorage.setItem("toggleVolume", "false");
        document.getElementById("volume-icon").src = `/static/images/volume-off-${theme}.png`;
        document.getElementById("volume-icon").title= "Volume Off";
        noPlay();
        return;

    }
}

function volumeOnLoad(){
    console.log("The dark mode status is: " + localStorage.getItem("toggleDarkMode"));
    if (localStorage.getItem("toggleVolume") == null) {
        localStorage.setItem("toggleVolume", "false");
        document.getElementById("volume-icon").src`/static/images/volume-off-${theme}.png`; 

    }
    if (localStorage.getItem("toggleVolume") == "true") {
        document.getElementById("volume-icon").src = `/static/images/volume-on-${theme}.png`;
    }else{
        document.getElementById("volume-icon").src = `/static/images/volume-off-${theme}.png`;
    }
    if (document.getElementById('card-count')) {

        var cardCount=document.getElementById("card-count").textContent;

    }
    if (document.getElementById('total-card-count')) {

        var totalCardCount=document.getElementById("total-card-count").textContent;
        var progressbarval=cardCount/totalCardCount*100; 
        document.getElementById("prog-value").setAttribute("value", progressbarval); 

    }

    tts();
}

function tts(){
    if (localStorage.getItem("toggleVolume") == "true") {
        var elements = document.getElementsByClassName("flashcard");
        for (var i = 0; i < elements.length; i++) {
            var text = elements[i].textContent; 
            var msg = new SpeechSynthesisUtterance();
            localStorage.setItem("currentAudio",msg);
            msg.text = text;
            window.speechSynthesis.speak(msg);
            
        }
    }
}

function toggleDarkMode() {
        var element = document.body;
        if (theme=='light'){
            theme = 'dark';
            element.classList.toggle("dark-mode");
            document.getElementById('volume-icon').src = `/static/images/volume-off-${theme}.png`; 
            document.getElementById('exit-icon').src = `/static/images/exit-${theme}.png`; 
            document.getElementById('dark-light-icon').src = `/static/images/light-dark-${theme}.png`; 
            const startflashcard = document.getElementsByClassName("startflashcard");
            const flashcard = document.getElementsByClassName("startflashcard");
            startflashcard[0].style.backgroundColor = '#000000';
            flashcard[0].style.backgroundColor = '#000000';
            return;
        }
        if (theme=='dark'){
            theme = 'light';
            element.classList.toggle("dark-mode");
            document.getElementById('volume-icon').src = `/static/images/volume-off-${theme}.png`; 
            document.getElementById('exit-icon').src = `/static/images/exit-${theme}.png`; 
            document.getElementById('dark-light-icon').src = `/static/images/light-dark-${theme}.png`; 
            const startflashcard = document.getElementsByClassName("startflashcard");
            const flashcard = document.getElementsByClassName("startflashcard");
            startflashcard[0].style.backgroundColor = '#ebebfa';
            flashcard[0].style.backgroundColor = '#ebebfa';
            return;
        }
    }
