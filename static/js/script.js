let theme = localStorage.getItem("theme")

let prefersColourScheme = window.matchMedia('(prefers-color-scheme: light)').matches ? "light" : "dark";
let overrideColourScheme = localStorage.getItem("overrideColourScheme") 

window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change",   e => e.matches && toggleDarkMode());
window.matchMedia("(prefers-color-scheme: light)").addEventListener("change",   e => e.matches && toggleDarkMode());

if (theme == null){
            localStorage.setItem("theme", prefersColourScheme);
    }
    
function goPlay(musicurl,redirecturl){
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
    setelementsTheme();


    if (localStorage.getItem("toggleVolume") == null) {
        localStorage.setItem("toggleVolume", "false");
    }

    if (localStorage.getItem("toggleVolume") == "true") {
        document.getElementById("volume-icon").src = `/static/images/volume-on-${theme}.png`;
        tts();
    } else {
        document.getElementById("volume-icon").src = `/static/images/volume-off-${theme}.png`;
    }
    if (document.getElementById('card-count')) {

        var cardCount=document.getElementById("card-count").textContent;

    }

    if (document.getElementById('total-card-count')) {
        var totalCardCount = document.getElementById("total-card-count").textContent;
        var progressbarval = cardCount / totalCardCount * 100;
        document.getElementById("prog-value").setAttribute("value", progressbarval);
    }


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
            localStorage.setItem("theme", "dark");
            element.classList.toggle("dark-mode");
            window.location.reload();
            return;
        }
        if (theme=='dark'){
            localStorage.setItem("theme", "light");
            element.classList.toggle("dark-mode");
            window.location.reload();
            return;
        }
    }

function setelementsTheme(){
    if (theme === prefersColourScheme) {
        theme = prefersColourScheme;
        localStorage.setItem("overrideColourScheme", "false");
    } else{ 
        localStorage.setItem("overrideColourScheme", "true");
    }

    var element = document.body;
    element.classList.toggle(`${theme}-mode`);
    document.getElementById('volume-icon').src = `/static/images/volume-off-${theme}.png`; 
    document.getElementById('exit-icon').src = `/static/images/exit-${theme}.png`; 
    document.getElementById('dark-light-icon').src = `/static/images/light-dark-${theme}.png`; 
    const startflashcarddiv = document.getElementsByClassName("startflashcard");
    const flashcarddiv = document.getElementsByClassName("flashcard");
    if (startflashcarddiv.length > 0) {
    if (theme =='light'){
        startflashcarddiv[0].style.backgroundColor = '#ebebfa';
    }
    if (theme =='dark'){
        startflashcarddiv[0].style.backgroundColor = '#000000';
        }
    }
    if (flashcarddiv.length > 0) {
        if (theme =='light'){
            flashcarddiv[0].style.backgroundColor = '#ebebfa';
        }
        if (theme =='dark'){
            flashcarddiv[0].style.backgroundColor = '#000000';
           }
        }
    }
