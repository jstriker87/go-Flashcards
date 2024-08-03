let theme = localStorage.getItem("theme")
if (theme == null){
            localStorage.setItem("theme", "light");
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
    console.log(localStorage.getItem("toggleVolume"));
    console.log(localStorage.getItem("theme"));
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
        console.log("The theme is:" + theme);
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
    var element = document.body;
    element.classList.toggle(`${theme}-mode`);
    document.getElementById('volume-icon').src = `/static/images/volume-off-${theme}.png`; 
    document.getElementById('exit-icon').src = `/static/images/exit-${theme}.png`; 
    document.getElementById('dark-light-icon').src = `/static/images/light-dark-${theme}.png`; 
    const startflashcarddiv = document.getElementsByClassName("startflashcard");
    const flashcarddiv = document.getElementsByClassName("flashcard");
    if (startflashcarddiv.length > 0) {
    if (theme =='light'){
        console.log("Light theme");
        startflashcarddiv[0].style.backgroundColor = '#ebebfa';
    }
    if (theme =='dark'){
        console.log("Dark theme");
        startflashcarddiv[0].style.backgroundColor = '#000000';
        }
    if (flashcarddiv.length > 0) {
        if (theme =='light'){
            flashcarddiv[0].style.backgroundColor = '#ebebfa';
            console.log("Light theme");
        }
        if (theme =='dark'){
            flashcarddiv[0].style.backgroundColor = '#000000';
           }
        }
    }
}
