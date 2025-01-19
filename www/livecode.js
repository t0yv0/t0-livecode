var codeMirrorEditor = undefined;

htmx.onLoad(function(content) {
    if (content.classList.contains("editor")) {
        initCodeMirror(content);
    } else {
        var editor = content.querySelector(".editor");
        if (editor) {
            initCodeMirror(editor);
        }
    }
});

function currentCode() {
    return codeMirrorEditor.getValue();
}

function initCodeMirror(editor) {
    var t = editor.innerText;
    editor.innerText = "";
    codeMirrorEditor = CodeMirror(editor, {value: t, mode: "javascript"});
    codeMirrorEditor.setSize(null /*retain width*/, window.innerHeight-50);
    CodeMirror.keyMap["default"]["Ctrl-S"] = function () {
        document.querySelector("#run").click();
    };
}
