import { Group, Text } from '@mantine/core';
import { modals } from '@mantine/modals';
import { IconAlertTriangle } from '@tabler/icons-react';

type AlertPopupProps = {
  message?: string;
  onOk: () => Promise<void> | void;
  onCancel?: () => void;
};
export const showConfirmationPopup = ({
  message,
  onOk,
  onCancel = () => {},
}: AlertPopupProps) => {
  modals.openConfirmModal({
    title: (
      <Group c="orange">
        <IconAlertTriangle stroke={2} size={20} />
        Confirm action
      </Group>
    ),
    children: (
      <Text c="dimmed" size="sm">
        {message || 'Are you sure you want to perform this action?'}
      </Text>
    ),
    labels: { confirm: 'OK', cancel: 'Cancel' },
    onConfirm: () => onOk(),
    onCancel: () => onCancel(),
    withCloseButton: false,
    closeOnClickOutside: false,
  });
};
