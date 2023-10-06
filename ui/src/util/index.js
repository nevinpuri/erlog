export function timeConverter(UNIX_timestamp) {
  var a = new Date(UNIX_timestamp * 1000);
  var months = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];
  var year = a.getFullYear();
  var month = months[a.getMonth()];
  var date = a.getDate();
  var hour = a.getHours();
  //   var min = a.getMinutes();
  //   var sec = a.getSeconds();

  var min = "0" + a.getMinutes();
  // Seconds part from the timestamp
  var sec = "0" + a.getSeconds();
  var time =
    date +
    " " +
    month +
    " " +
    year +
    " " +
    hour +
    ":" +
    min.slice(-2) +
    ":" +
    sec.slice(-2);
  return time;
}
