import { Alert, Stack, Text } from '@mantine/core';
import { modals } from '@mantine/modals';
import { IconInfoCircle } from '@tabler/icons-react';

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
    children: (
      <Stack align="center">
        <Alert
          color="#E4A11B"
          variant="light"
          icon={<IconInfoCircle size={30} />}
        >
          <Text size="sm" c="#E4A11B">
            {message || 'Are you sure you want to perform this action?'}
          </Text>
        </Alert>
      </Stack>
    ),
    labels: { confirm: 'OK', cancel: 'Cancel' },
    onConfirm: () => onOk(),
    onCancel: () => onCancel(),
    withCloseButton: false,
    closeOnClickOutside: false,
  });
};
