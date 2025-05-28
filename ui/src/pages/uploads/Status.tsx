import { ThemeIcon } from '@mantine/core';
import {
  IconAlertHexagon,
  IconCancel,
  IconCircleCheck,
  IconClockExclamation,
  IconHelpCircle,
  IconProgressDown,
  IconProgressX,
  IconTimeDuration0,
  IconTrash,
} from '@tabler/icons-react';
import React from 'react';

const statusConfig: Record<string, { color: string; icon: React.ReactNode }> = {
  Started: { color: 'blue', icon: <IconTimeDuration0 size={18} /> },
  Skipped: { color: 'gray', icon: <IconProgressX size={18} /> },
  'In Progress': { color: 'blue', icon: <IconProgressDown size={18} /> },
  Finished: { color: 'green', icon: <IconCircleCheck size={18} /> },
  Failed: { color: 'red', icon: <IconAlertHexagon size={18} /> },
  Cancelled: { color: 'gray', icon: <IconCancel size={18} /> },
  Deleted: { color: 'gray', icon: <IconTrash size={18} /> },
  Unknown: { color: 'gray', icon: <IconHelpCircle size={18} /> },
  'Timed Out': { color: 'red', icon: <IconClockExclamation size={18} /> },
};

export const UploadStatus = ({ status = 'Unknown' }: { status?: string }) => {
  const { color, icon } = statusConfig[status] || statusConfig.Unknown;

  return (
    <ThemeIcon color={color} variant="subtle">
      {icon}
    </ThemeIcon>
  );
};
