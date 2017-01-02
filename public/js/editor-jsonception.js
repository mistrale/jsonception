$("#clear-run-btn").click(function() {
  run_editor.setValue("#!/bin/bash")
});

$("#clear_button_exec").click(function() {
  out_editor.setValue("")
});

  var run_editor = CodeMirror.fromTextArea(document.getElementById("run_area"), {
    lineNumbers: true,
    theme:"mdn-like",
    mode: "text/x-sh",
    styleActiveLine:true,
    matchBrackets: true
  });

  var out_editor = CodeMirror.fromTextArea(document.getElementById("output_area"), {
    matchBrackets: true,
    theme:"mdn-like",
    lineNumbers: true,
    readOnly : true,
    styleActiveLine:true
    });

  var logevent_editor = CodeMirror.fromTextArea(document.getElementById("logevent_area"), {
      matchBrackets: true,
      theme:"mdn-like",
      mode:"javascript",
      lineNumbers: true,
      styleActiveLine:true
      });

  var config_editor = CodeMirror.fromTextArea(document.getElementById("config_area"), {
      matchBrackets: true,
      theme:"mdn-like",
      mode:"javascript",
      lineNumbers: true,
      styleActiveLine:true
      });
