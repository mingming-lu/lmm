export {
  buildPageNumbers,
  buildURLEncodedString,
  formattedDate,
  formattedTimeStamp,
  formattedUTCString,
  range,
}

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

const buildURLEncodedString = obj => {
  return Object.entries(obj).filter(kv => {
    return kv[1]
  }).map(kv => {
    return `${kv[0]}=${kv[1]}`
  }).join('&')
}

const range = (from, to, interval=1) => {
  if (from > to) {
    return []
  }

  return Array(to - from + 1).fill(from).map((v, i) => {
    return v + (i * interval)
  })
}

const buildPageNumbers = (page, total, maxItem=5) => {
  if (page <= Math.ceil(maxItem / 2)) {
    return range(1, Math.min(total, maxItem))
  }

  if (total - page < maxItem / 2) {
    return range(total - maxItem + 1, total)
  }

  let growLeft = Math.floor(maxItem / 2)
  let left = page - growLeft
  if (left < 1) {
    left = 1
  }

  let growRight = Math.floor(maxItem / 2)
  let right = page + growRight
  if (right > total) {
    right = total
  }

  return range(left, right)
}
