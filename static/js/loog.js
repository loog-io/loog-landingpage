function loog(id){
    var r = new XMLHttpRequest();
    r.open("POST", "http://192.168.0.103:5050/loog", true);

    var d = new FormData();
    d.append("s", id);
    d.append("l", document.location);
    if(document.referrer){d.append("r", document.referrer);}
    d.append("w", window.innerWidth);
    d.append("h", window.innerHeight);
    if(sessionStorage.loogviews){
        sessionStorage.loogviews = Number(sessionStorage.loogviews) + 1;
        d.append("v", sessionStorage.loogviews);
    } else {
        sessionStorage.loogviews = 1;
        d.append("v", 1);
    }

    r.send(d);
}
loog(loog_site);