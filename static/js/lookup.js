function beacon_btn_click(e) {
  var targ = e.target;
  var inner = targ.innerHTML;
  var tr = targ.parentElement.parentElement;
  var targ_class_1 = "record-row-" + targ.getAttribute('bkey');
  var targ_class_2 = "record-point-row-" + targ.getAttribute('bkey');

  var targ_rows_1 = document.getElementsByClassName(targ_class_1);
  var targ_rows_2 = document.getElementsByClassName(targ_class_2);

  //var targ_rows = [...targ_rows_1, ...targ_rows_2];
  var targ_rows = [...targ_rows_1, ...targ_rows_2];

  if (inner === "+") {
    for (var i=0; i <targ_rows_1.length; i++) {
      targ_rows[i].removeAttribute('hidden');
    }
    targ.innerHTML = "-";
  } else if (inner === "-") {
    for (var i=0; i <targ_rows.length; i++) {
      targ_rows[i].setAttribute('hidden', 'hidden');
    }
    targ.innerHTML = "+";
  }
}

function record_btn_click(e) {
  var targ = e.target;
  var inner = targ.innerHTML;
  var tr = targ.parentElement.parentElement;
  var targ_class = "record-point-row-" + targ.getAttribute('bkey');

  var targ_rows = document.getElementsByClassName(targ_class);

  if (inner === "+") {
    for (var i=0; i <targ_rows.length; i++) {
      targ_rows[i].removeAttribute('hidden');
    }
    targ.innerHTML = "-";
  } else if (inner === "-") {
    for (var i=0; i <targ_rows.length; i++) {
      targ_rows[i].setAttribute('hidden', 'hidden');
    }
    targ.innerHTML = "+";
  }
}

function ijs_setup() {
  var beacon_btns = document.getElementsByClassName('beacon-expand-btn');
  for (var i=0; i<beacon_btns.length; i++) {
    beacon_btns[i].addEventListener('click', beacon_btn_click);
  }
  var record_btns = document.getElementsByClassName('record-expand-btn');
  for (var i=0; i<record_btns.length; i++) {
    record_btns[i].addEventListener('click', record_btn_click);
  }
}

document.addEventListener("DOMContentLoaded", ijs_setup);
