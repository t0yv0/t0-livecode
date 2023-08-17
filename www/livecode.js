function livecode(prog) {
    var cm = CodeMirror(document.getElementById("editor"), {mode: "javascript"});
    cm.setSize(null /*retain width*/, window.innerHeight-50);
    document.getElementById("previewFrame").style.height = (window.innerHeight-50) + "px";
    document.getElementById("previewFrame").src = "/program/"+prog+"/";

    async function fetchSource() {
        const response = await fetch("/program/"+prog+"/script.js", {
            method: "GET",
            mode: "same-origin",
            cache: "no-cache",
        });
        return response.text();
    }

    fetchSource().then(fetched => {
        cm.setValue(fetched);
    });

    async function update() {
        const response = await fetch("/program/"+prog, {
            method: "POST",
            mode: "same-origin",
            cache: "no-cache",
            headers: {"Content-Type": "text/plain"},
            body: cm.getValue(),
        });
        return response.json();
    }

    function save() {
        console.log("Saving", prog);
        update().then((resp) => {
            console.log("Saved", prog);
            document.getElementById("previewFrame").contentWindow.location.reload();
        });
    }

    CodeMirror.keyMap["default"]["Ctrl-S"] = save;
}
