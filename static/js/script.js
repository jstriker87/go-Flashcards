
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
        document.getElementById("volume-icon").src = "/static/images/volume-on-light.png";
        document.getElementById("volume-icon").title= "Volume On";
        tts();
        return;
    } else {
        localStorage.setItem("toggleVolume", "false");
        document.getElementById("volume-icon").src = "/static/images/volume-off-light.png";
        document.getElementById("volume-icon").title= "Volume Off";
        noPlay();
        return;

    }
}

function volumeOnLoad(){
    console.log("The dark mode status is: " + localStorage.getItem("toggleDarkMode"));
    if (localStorage.getItem("toggleVolume") == null) {
        localStorage.setItem("toggleVolume", "false");
        document.getElementById("volume-icon").src = "/static/images/volume-off-light.png";

    }
    if (localStorage.getItem("toggleVolume") == "true") {
        document.getElementById("volume-icon").src = "/static/images/volume-on-light.png";
    }else{
        document.getElementById("volume-icon").src = "/static/images/volume-off-light.png";
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
    if (localStorage.getItem("toggleDarkMode") == "false") {
        localStorage.setItem("toggleDarkMode", "true");
        console.log("The toggled dark mode status is: " + localStorage.getItem("toggleDarkMode"));
        element.classList.toggle("dark-mode");
        return;
    }
    if (localStorage.getItem("toggleDarkMode") == "true") {
        localStorage.setItem("toggleDarkMode", "false");
        console.log("The toggled dark mode status is: " + localStorage.getItem("toggleDarkMode"));
        element.classList.toggle("light-mode");
        return;
    }
    console.log("The dark mode status is: " + localStorage.getItem("toggleDarkMode"));

    }
