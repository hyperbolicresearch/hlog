export const timeAgo = (interval: number) : string => {
  let str = "("
  const days = Math.floor(interval / 84600)
  let remaining = Math.floor(interval % 84600)
  const hours = Math.floor(remaining / 3600)
  remaining = Math.floor(remaining % 3600)
  const minutes = Math.floor(remaining / 60)
  remaining = Math.floor(remaining % 60)
  const seconds = Math.floor(remaining)

  if (days != 0) { str += `${days} days `}
  if (hours != 0) { str += `${hours} hours `}
  if (minutes != 0) { str += `${minutes} minutes `}
  if (seconds != 0) { str += `${seconds} seconds`} else {
      str += "just now"
  }
  return str + ")"
}