export const statusConfig: Record<string, string> = {
  pending: "#A9A9A9", // Dark Gray - Waiting to start
  started: "#1E90FF", // Dodger Blue - Actively running
  in_progress: "#007BFF", // Blue - Ongoing process
  waiting: "#FFD700", // Gold - Paused, waiting for input/approval
  on_hold: "#FF8C00", // Dark Orange - Temporarily halted
  completed: "#32CD32", // Lime Green - Successfully finished
  failed: "#FF0000", // Red - Process encountered an error
  cancelled: "#808080", // Gray - Process was terminated
  expired: "#B22222", // Firebrick - Process timed out
  retrying: "#FFA500", // Orange - Attempting to redo
  escalated: "#FF4500", // Orange-Red - Raised to a higher authority
  reviewing: "#8A2BE2", // Blue Violet - Under evaluation
  approved: "#228B22", // Forest Green - Successfully authorized
  rejected: "#DC143C", // Crimson - Declined
  scheduled: "#20B2AA", // Light Sea Green - Planned for execution
  deferred: "#696969", // Dim Gray - Delayed for later processing
  aborted: "#800000", // Maroon - Forcefully stopped
  terminated: "#808080", // Gray - Process was terminated
};
