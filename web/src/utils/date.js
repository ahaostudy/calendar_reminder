export function dateString(date) {
  let year = date.getFullYear()
  let mouth = date.getMonth() + 1
  if (mouth < 10) mouth = `0${mouth}`
  let day = date.getDate()
  if (day < 10) day = `0${day}`
  return `${year}-${mouth}-${day}`
}
