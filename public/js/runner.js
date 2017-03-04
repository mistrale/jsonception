function init(uuid) {
  script_editor[uuid] = CodeMirror.fromTextArea(document.getElementById("script_area_" + uuid), {
    lineNumbers: true,
    theme:"mdn-like",
    mode: "text/x-sh",
    styleActiveLine:true,
    matchBrackets: true
    });

  output_script_editor[uuid] = CodeMirror.fromTextArea(document.getElementById("output_script_area_" + uuid), {
    matchBrackets: true,
    theme:"mdn-like",
    lineNumbers: true,
    styleActiveLine:true
    });
}

function init_test(uuid) {
  config_editor[uuid] = CodeMirror.fromTextArea(document.getElementById("config_area_" + uuid), {
  matchBrackets: true,
  theme:"mdn-like",
  mode : "javascript",
  lineNumbers: true,
  styleActiveLine:true
  });
}

function init_history(uuid) {

  output_script_history_test[uuid] = CodeMirror.fromTextArea(document.getElementById("output_script_area_history_" + uuid), {
    matchBrackets: true,
    theme:"mdn-like",
    lineNumbers: true,
    styleActiveLine:true
    });

  event_ref_history_editor[uuid] = CodeMirror.fromTextArea(document.getElementById("event_ref_area_history_" + uuid), {
  matchBrackets: true,
  theme:"mdn-like",
  lineNumbers: true,
  styleActiveLine:true,
  readOnly : true,
  mode : "javascript",
  });

  event_log_history_editor[uuid] = CodeMirror.fromTextArea(document.getElementById("event_log_area_history_" + uuid), {
  matchBrackets: true,
  theme:"mdn-like",
  lineNumbers: true,
  styleActiveLine:true,
  readOnly : true,
  mode : "javascript",
  });
  console.log(output_script_history_test)

}

function init_run_test(uuid) {
  output_script_history_test[uuid].refresh()
  event_ref_history_editor[uuid].refresh()
  event_log_history_editor[uuid].refresh()

  output_script_history_test[uuid].setValue("")
  event_ref_history_editor[uuid].setValue("")
  event_log_history_editor[uuid].setValue("")



  queues_test_test[uuid] = []
  queues_test_ref[uuid] = []

  timers_test[uuid] = window.setInterval(function() {
    if (queues_test_test[uuid].length > 0) {
      var i = queues_test_test[uuid].shift();
      cursor = event_log_history_editor[uuid].getCursor()
      event_log_history_editor[uuid].replaceRange(i, cursor)
      event_log_history_editor[uuid].scrollTo(0, 10000000)
    }

    if (queues_test_ref[uuid].length > 0) {
      var i = queues_test_ref[uuid].shift();
      cursor = event_ref_history_editor[uuid].getCursor()
      event_ref_history_editor[uuid].replaceRange(i, cursor)
      event_ref_history_editor[uuid].scrollTo(0, 10000000)
    }
    if (isFinished_test[uuid] == true && queues_test_test[uuid].length == 0 && queues_test_ref[uuid].length == 0) {
      window.clearInterval(timers_test[uuid]);
    }
  }, 1 );
}

function init_run(uuid) {
  isFinished[uuid] = false

  queues[uuid] = []

  timers[uuid] = window.setInterval(function() {
    if (queues[uuid].length > 0) {
      var i = queues[uuid].shift();

      if (output_script_editor[uuid]) {

        console.log(output_script_editor[uuid].getCursor())
        cursor = output_script_editor[uuid].getCursor()
        if (cursor.ch >= 120) {
          i += "\n"
        }
        output_script_editor[uuid].replaceRange(i, cursor)
        output_script_editor[uuid].scrollTo(0, 10000000)
      }
      if (output_script_history_test[uuid]) {
        output_script_history_test[uuid].replaceRange(i, output_script_history_test[uuid].getCursor())
        output_script_history_test[uuid].scrollTo(0, 10000000)
      }
    }
    if (isFinished[uuid] == true && queues[uuid].length == 0) {

      window.clearInterval(timers[uuid]);
    }
  }, 0.1 );
}
