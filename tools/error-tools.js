window.onerror = function(msg, url, line, col, error) {  
    console.log("msg " + msg); // エラーの内容
    console.log("url " + url); // エラーの内容
    console.log("line:col " + line + ":" + col ); // エラーの内容

};