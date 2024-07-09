function goPlay(musicurl,url){
  var music = new Audio(musicurl);
  music.play();
  sleep(1500).then(() => { window.location.href = url; }); 
  }

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
function volume() {

    console.log("The Start volume when vol function run is " + localStorage.getItem("toggleVolume"));
    if (localStorage.getItem("toggleVolume") == "false") {
        localStorage.setItem("toggleVolume", "true");
        document.getElementById("volume-icon").src = "/static/images/volume-on.svg";
        updateVolume(true);
        return;
    } else {
        localStorage.setItem("toggleVolume", "false");
        updateVolume(false);
        return;

    }
}

function updateVolume(isMuted){
    
    var okButton = document.getElementsByClassName("ok-button");
    var needsRevisonButton = document.getElementsByClassName("revision-button");
      
    if (isMuted==true){
        okButton.muted = true;
        needsRevisonButton.muted = true;
        document.getElementById("volume-icon").src = "/static/images/volume-on.svg";

    } 

    if (isMuted==false){
        okButton.muted = false;
        needsRevisonButton.muted = false;
        document.getElementById("volume-icon").src = "/static/images/volume-off.svg";
    }

}

function volumeOnLoad(){
    console.log("The volume setting on load is " + localStorage.getItem("toggleVolume"));
    if (localStorage.getItem("toggleVolume") == "true") {
        document.getElementById("volume-icon").src = "/static/images/volume-on.svg";
    }else{
        
        document.getElementById("volume-icon").src = "/static/images/volume-off.svg";
    }
}
