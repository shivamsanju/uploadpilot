import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en";

TimeAgo.addDefaultLocale(en);
export const timeAgo = new TimeAgo("en-US");

export const formatMilliseconds = (ms: number = 0) => {
  if (ms < 1000) return `${ms} ms`;

  const seconds = Math.floor(ms / 1000);
  if (seconds < 60) return `${seconds} sec`;

  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes} min`;

  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hr`;

  const days = Math.floor(hours / 24);
  return `${days} day${days > 1 ? "s" : ""}`;
};
