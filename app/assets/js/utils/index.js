export { formattedDate, formattedTimeStamp, formattedUTCString }

function formattedTimeStamp(timestamp) {
  let date = new Date(timestamp * 1e3)
  return formattedDate(date)
}

function formattedUTCString(s) {
  let date = new Date(s)
  return formattedDate(date)
}

function formattedDate(date) {
  const y = date.getFullYear()
  const m = date.getMonth() + 1
  const d = date.getDate() < 10 ? `0${date.getDate()}` : date.getDate()

  return `${y}-${m}-${d} ${date.toTimeString().substr(0, 8)}`
}
