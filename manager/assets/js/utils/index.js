export {
  buildURLEncodedQueryString,
  formattedDate,
  formattedDateFromTimeStamp,
  formattedDateFromString
}

const formattedDateFromTimeStamp = timestamp => {
  let date = new Date(timestamp * 1e3)
  return formattedDate(date)
}

const formattedDateFromString = s => {
  let date = new Date(s)
  return formattedDate(date)
}

const formattedDate = date => {
  const y = date.getFullYear()
  const m = (() => {
    const _m = date.getMonth() + 1
    return _m < 10 ? `0${_m}` : _m
  })()
  const d = date.getDate() < 10 ? `0${date.getDate()}` : date.getDate()

  return `${y}-${m}-${d} ${date.toTimeString().substr(0, 8)}`
}

const buildURLEncodedQueryString = obj => {
  return Object.entries(obj)
    .filter(kv => {
      return kv[1]
    })
    .map(kv => {
      return `${kv[0]}=${encodeURIComponent(kv[1])}`
    })
    .join('&')
}
