var myCodeMirror = CodeMirror(document.getElementById("editor"), {
  value: "function myScript(){return 100;}\n",
  mode:  "javascript"
});

myCodeMirror.setSize(null /*retain width*/, window.innerHeight-50);
document.getElementById("previewFrame").style.height = (window.innerHeight-50) + "px";

async function update() {
    const response = await fetch("/update/", {
        method: "POST",
        mode: "same-origin",
        cache: "no-cache",
        headers: {"Content-Type": "text/plain"},
        body: myCodeMirror.getValue(),
    });
    return response.json();
}

function save() {
    console.log("Saving");
    update().then((resp) => {
        console.log("Saved", resp);
    });
}

CodeMirror.keyMap.default['Ctrl-S'] = save;
