import { Loader, ThemeIcon } from '@mantine/core';
import {
  IconAlertCircle,
  IconCancel,
  IconCircleCheck,
  IconCircleDot,
  IconClockExclamation,
  IconHelpCircle,
} from '@tabler/icons-react';
import React from 'react';

const statusConfig: Record<string, { color: string; icon: React.ReactNode }> = {
  Running: { color: 'gray', icon: <Loader size={16} /> },
  Completed: { color: 'green', icon: <IconCircleCheck size={18} /> },
  Terminated: { color: 'red', icon: <IconAlertCircle size={18} /> },
  Failed: { color: 'red', icon: <IconAlertCircle size={18} /> },
  Canceled: { color: '#808080', icon: <IconCancel size={18} /> },
  'Continued-As-New': { color: 'gray', icon: <IconCircleDot size={18} /> },
  'Timed Out': { color: 'red', icon: <IconClockExclamation size={18} /> },
  Unknown: { color: 'gray', icon: <IconHelpCircle size={18} /> },
};

export const WorkflowStatus = ({ status = 'Unknown' }: { status?: string }) => {
  const { color, icon } = statusConfig[status] || statusConfig.Unknown;

  return (
    <ThemeIcon color={color} variant="subtle">
      {icon}
    </ThemeIcon>
  );
};
