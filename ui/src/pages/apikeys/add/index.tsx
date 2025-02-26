import { CodeHighlight } from '@mantine/code-highlight';
import {
  Box,
  Button,
  Group,
  LoadingOverlay,
  Stack,
  Text,
  TextInput,
} from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { useForm } from '@mantine/form';
import { modals } from '@mantine/modals';
import { IconAlertTriangle } from '@tabler/icons-react';
import { useCreateApiKeyMutation } from '../../../apis/apikeys';
import { CreateApiKeyData } from '../../../types/apikey';

type Props = {
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
  workspaceId: string;
};

const CreateApiKeyForm: React.FC<Props> = ({ setOpened, workspaceId }) => {
  const form = useForm<CreateApiKeyData>({
    initialValues: {
      expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      name: '',
    },
    validate: {
      expiresAt: value => {
        if (value.getTime() < Date.now()) {
          return 'Expiry date must be in the future';
        }
        return null;
      },
      name: value => {
        if (!value || value.length < 3 || value.length > 25) {
          return 'Name must be between 3 and 25 characters';
        }
        return null;
      },
    },
  });

  const { mutateAsync, isPending } = useCreateApiKeyMutation();

  const handleAdd = async (data: any) => {
    try {
      const newKey = await mutateAsync({
        workspaceId,
        data,
      });
      setOpened(false);
      modals.open({
        title: 'Api key created successfully',
        centered: true,
        closeOnEscape: false,
        closeOnClickOutside: false,
        padding: 'lg',
        size: 'xl',
        children: (
          <Box>
            <Group align="flex-start" mb="md">
              <IconAlertTriangle color="orange" size="18" />
              <Text size="sm" c="orange">
                Please copy the api key and store it safely, as it will not be
                displayed again.
              </Text>
            </Group>

            <CodeHighlight code={newKey} />
          </Box>
        ),
      });
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <form onSubmit={form.onSubmit(handleAdd)}>
      <LoadingOverlay
        visible={isPending}
        overlayProps={{ backgroundOpacity: 0 }}
        zIndex={1000}
      />
      <Stack gap="xl">
        <TextInput
          label="Name"
          withAsterisk
          description="A name to identify the api key (between 3-25 characters)"
          placeholder="Enter a comment"
          {...form.getInputProps('name')}
        />
        <DateInput
          label="Expiry Date"
          withAsterisk
          description="The expiry date of the api key"
          placeholder="Enter an expiry date"
          {...form.getInputProps('expiresAt')}
        />
      </Stack>
      <Group justify="flex-end" mt={50}>
        <Button type="submit">Create</Button>
      </Group>
    </form>
  );
};

export default CreateApiKeyForm;
