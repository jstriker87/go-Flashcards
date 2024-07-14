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
}



