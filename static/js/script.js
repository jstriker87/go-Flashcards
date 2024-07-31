function goPlay(musicurl,url){
  if (localStorage.getItem("toggleVolume") == "true") {
    var music = new Audio(musicurl);
    music.play();
  }
  sleep(700).then(() => { window.location.href = url; }); 
  }

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
function volume() {
    if (localStorage.getItem("toggleVolume") == "false") {
        localStorage.setItem("toggleVolume", "true");
        document.getElementById("volume-icon").src = "/static/images/volume-on.svg";
        document.getElementById("volume-icon").title= "Volume On";
        return;
    } else {
        localStorage.setItem("toggleVolume", "false");
        document.getElementById("volume-icon").src = "/static/images/volume-off.svg";
        document.getElementById("volume-icon").title= "Volume Off";
        return;

    }
}

function volumeOnLoad(){
    if (localStorage.getItem("toggleVolume") == null) {
        localStorage.setItem("toggleVolume", "false");
    }
    if (localStorage.getItem("toggleVolume") == "true") {
        document.getElementById("volume-icon").src = "/static/images/volume-on.svg";
    }else{
        document.getElementById("volume-icon").src = "/static/images/volume-off.svg";
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
            msg.text = text;
            window.speechSynthesis.speak(msg);
        }
    }
}
