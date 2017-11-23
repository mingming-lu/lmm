export { formattedTime }

function formattedTime (timestamp) {
  let date = new Date(timestamp * 1e3).toISOString().slice(0, 10)
  return date.split('-').join('/')
}
