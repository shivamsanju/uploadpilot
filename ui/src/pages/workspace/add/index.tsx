import {
  Box,
  Button,
  Container,
  Group,
  Stack,
  TagsInput,
  Textarea,
  TextInput,
  Title,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useNavigate } from 'react-router-dom';
import { useCreateWorkspaceMutation } from '../../../apis/workspace';
import { ContainerOverlay } from '../../../components/Overlay';
import { CreateWorkspaceData } from '../../../types/workspace';

type Props = {};

const CreateWorkspacePage: React.FC<Props> = () => {
  const navigate = useNavigate();

  const form = useForm<CreateWorkspaceData>({
    initialValues: {
      name: '',
      description: '',
      tags: [],
    },
    validate: {
      name: value => {
        if (!value || value.length < 3 || value.length > 25) {
          return 'Name must be between 3 and 25 characters';
        }
        return null;
      },
      description: value => {
        if (value && value.length > 300) {
          return 'Description must be less than 300 characters';
        }
        return null;
      },
      tags: value => {
        if (value && value.length > 10) {
          return 'You can select a maximum of 10 tags';
        }
        return null;
      },
    },
  });

  const { mutateAsync, isPending: isCreating } = useCreateWorkspaceMutation();

  const handleAdd = async (data: any) => {
    try {
      const wsId = await mutateAsync(data);
      navigate(`/workspace/${wsId}`);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Box mb={50}>
      <ContainerOverlay visible={isCreating} />
      <Container mt="md">
        <form onSubmit={form.onSubmit(handleAdd)}>
          <Group justify="space-between" align="center" mb="lg">
            <Title order={3}>Create new workspace</Title>
            <Button type="submit">Create</Button>
          </Group>

          <Stack gap="lg">
            <TextInput
              label="Name"
              required
              placeholder="Enter a name to identify your workspace (between 3-25 characters)"
              {...form.getInputProps('name')}
            />

            <Textarea
              rows={4}
              label="Description"
              placeholder="Enter a description for your workspace (optional, max 300 characters)"
              {...form.getInputProps('description')}
            />

            <TagsInput
              label="Tags"
              placeholder="Enter a description for your workspace (optional)"
              {...form.getInputProps('tags')}
            />
          </Stack>
        </form>
      </Container>
    </Box>
  );
};

export default CreateWorkspacePage;
