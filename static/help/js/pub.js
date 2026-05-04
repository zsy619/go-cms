$(document).ready(function () {
    $(".navbar-burger").click(function () {
        $(".navbar-burger").toggleClass("is-active");
        $(".navbar-menu").toggleClass("is-active");
    });

    $(".wcms").click(function () {
        window.open("/link/wcms")
    });
    $(".lcms").click(function () {
        window.open("/link/lcms")
    });
    $(".mcms").click(function () {
        window.open("/link/mcms")
    });
    $(".lcms-arm").click(function () {
        window.open("/link/lcms-arm")
    });

    var url = location.href;
    $("ul[class='menu-list'] a").each(function () {
        var aurl = $(this).prop("href");
        if (url.indexOf(aurl) > -1) {
            $(this).addClass("is-active");
        }
    });
});

function createCodeMirror(id, typex) {
    var editor = CodeMirror.fromTextArea(document.getElementById(id), {
        theme: "dracula",
        lineNumbers: true,
        //括号匹配
        matchBrackets: true,
        lineNumberFormatter: function (number) {
            // 将行号格式化为三位数
            var str = String(number);
            while (str.length < 4) {
                str = "0" + str;
            }
            return str;
        },
        smartIndent: true,
        viewportMargin: Infinity,
        styleActiveLine: true,
        matchBrackets: true,
        //代码折叠
        lineWrapping: true,
        foldGutter: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
        mode: { name: typex, globalVars: true },
        extraKeys: {
            "F11": function (cm) {
                cm.setOption("fullScreen", !cm.getOption("fullScreen"));
            },
            "Esc": function (cm) {
                if (cm.getOption("fullScreen")) cm.setOption("fullScreen", false);
            }
        }
    });
    editor.setSize("auto", 'auto');
    editor.setOption("readOnly", true);
    // editor.setOption("lineWrapping", true);
    return editor;
}