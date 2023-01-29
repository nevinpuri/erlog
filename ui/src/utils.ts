export function toFormattedDate(_date: Date): string {
  // fixes getUTC functions not found
  let date = new Date(_date);
  const year = date.getUTCFullYear();
  const month = date.getUTCMonth() + 1;
  const curDate = date.getUTCDate();
  const hours = date.getUTCHours();
  const minutes = date.getUTCMinutes();
  const seconds = date.getUTCSeconds();
  const milliseconds = date.getUTCMilliseconds();

  return `${year}-${month}-${curDate}T${hours}:${minutes}:${seconds}.${milliseconds}`;
}
