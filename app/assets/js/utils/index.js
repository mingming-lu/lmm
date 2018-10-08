export {formattedDate,formattedTimeStamp,formattedUTCString}

function formattedTimeStamp(timestamp) {
  let date = new Date(timestamp * 1e3)
  return formattedDate(date)
}

function formattedUTCString(s) {
  let date = new Date(s)
  return formattedDate(date)
}

function formattedDate(date) {
  let d = {
    y: date.getFullYear(),
    m: date.getMonth(),
    d: date.getDate()
  }
  let s = ''
  s += d.y + '-'
  s += (d.m + 1) + '-'
  s += (d.d < 10 ? '0' + d.d : d.d)
  return s
}
