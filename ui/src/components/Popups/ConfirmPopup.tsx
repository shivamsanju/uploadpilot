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
          color="orange"
          variant="light"
          icon={<IconInfoCircle size={30} />}
        >
          <Text size="sm" c="orange">
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
    confirmProps: { variant: 'white' },
  });
};
