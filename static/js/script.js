function goPlay(musicurl,url){
  var music = new Audio(musicurl);
  music.play();
  sleep(1500).then(() => { window.location.href = url; }); 
  }

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
