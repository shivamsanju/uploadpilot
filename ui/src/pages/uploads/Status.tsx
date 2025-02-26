import { Loader, ThemeIcon } from '@mantine/core';
import {
  IconAlertHexagon,
  IconCancel,
  IconCircleCheck,
  IconHelpCircle,
  IconProgressX,
  IconRosetteDiscountCheckFilled,
  IconTimeDuration0,
  IconTrash,
} from '@tabler/icons-react';
import React from 'react';

const statusConfig: Record<string, { color: string; icon: React.ReactNode }> = {
  Started: { color: 'blue', icon: <IconTimeDuration0 size={20} /> },
  Skipped: { color: 'gray', icon: <IconProgressX size={20} /> },
  'In Progress': { color: 'teal', icon: <Loader size={16} type="oval" /> },
  Uploaded: { color: 'green', icon: <IconCircleCheck size={20} /> },
  Failed: { color: 'red', icon: <IconAlertHexagon size={20} /> },
  Cancelled: { color: 'gray', icon: <IconCancel size={20} /> },
  Processing: { color: 'teal', icon: <Loader size={20} /> },
  'Processing Failed': { color: 'red', icon: <IconAlertHexagon size={20} /> },
  'Processing Complete': {
    color: 'green',
    icon: <IconRosetteDiscountCheckFilled size={20} />,
  },
  'Processing Cancelled': { color: 'gray', icon: <IconCancel size={20} /> },
  Deleted: { color: 'gray', icon: <IconTrash size={20} /> },
  Unknown: { color: 'gray', icon: <IconHelpCircle size={20} /> },
};

export const UploadStatus = ({ status = 'Unknown' }: { status?: string }) => {
  const { color, icon } = statusConfig[status] || statusConfig.Unknown;

  return (
    <ThemeIcon color={color} variant="subtle">
      {icon}
    </ThemeIcon>
  );
};
