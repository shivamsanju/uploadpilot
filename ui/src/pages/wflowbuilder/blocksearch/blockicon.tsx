import {
  IconCubeSend,
  IconFileTextSpark,
  IconWebhook,
} from '@tabler/icons-react';
import { ReactNode } from 'react';

export const BLOCK_ICONS: { [key: string]: any } = {
  Webhook: IconWebhook,
  ExtractPDFContent: IconFileTextSpark,
};

export const getBlockIcon = (key: string, size?: number): ReactNode => {
  const IconComponent = BLOCK_ICONS[key];
  return IconComponent ? (
    <IconComponent size={size} />
  ) : (
    <IconCubeSend size={size} />
  );
};
