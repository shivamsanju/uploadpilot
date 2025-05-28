import {
  Box,
  Button,
  Container,
  Group,
  SimpleGrid,
  Stack,
  TagsInput,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconSearch } from '@tabler/icons-react';
import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  useCreateProcessorMutation,
  useGetAllWorkflowTemplates,
} from '../../../apis/processors';
import { ErrorCard } from '../../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../../components/Overlay';
import { Template } from './Template';

const NewprocessorPage = () => {
  const [activePage, setActivePage] = useState(0);
  const { workspaceId } = useParams();
  const navigate = useNavigate();

  const form = useForm({
    initialValues: {
      name: '',
      triggers: [],
      enabled: true,
      templateKey: '',
    },
    validate: {
      name: value => (value ? null : 'Name is required'),
    },
  });

  const { mutateAsync: createMutateAsync, isPending: isCreating } =
    useCreateProcessorMutation();
  const { isPending, templates, error } = useGetAllWorkflowTemplates(
    workspaceId || '',
  );

  const handleAdd = async (values: any) => {
    try {
      const procId = await createMutateAsync({
        workspaceId: workspaceId || '',
        processor: {
          workspaceId,
          name: values.name,
          triggers: values.triggers,
          templateKey: values.templateKey,
        },
      });
      console.log(procId);
      form.reset();
      navigate(`/workspace/${workspaceId}/processors/${procId}/workflow`);
    } catch (error) {
      console.error(error);
    }
  };

  const selectTemplate = (templateKey: string) => {
    if (templateKey === form.values.templateKey) {
      form.setFieldValue('templateKey', '');
    } else {
      form.setFieldValue('templateKey', templateKey);
    }
  };

  if (error) {
    return <ErrorCard message={error?.message} title={'Error'} />;
  }

  return (
    <Box mb={50}>
      <ContainerOverlay visible={isCreating || isPending} />

      <Container mt="md">
        <form onSubmit={form.onSubmit(handleAdd)}>
          <Group justify="space-between" align="center" mb="lg">
            <Box>
              <Title order={3}>Create new processor</Title>
            </Box>
            <Group justify="flex-end">
              {activePage === 0 ? (
                <Button
                  variant="default"
                  onClick={() =>
                    navigate(`/workspace/${workspaceId}/processors`)
                  }
                >
                  Discard
                </Button>
              ) : (
                <Button
                  onClick={() => setActivePage(activePage - 1)}
                  variant="default"
                >
                  Back
                </Button>
              )}
              {activePage === 0 && (
                <Button onClick={() => setActivePage(activePage + 1)}>
                  Next
                </Button>
              )}
              {activePage === 1 && (
                <Button
                  type="submit"
                  loading={isCreating}
                  disabled={!form.isValid()}
                >
                  Create
                </Button>
              )}
            </Group>
          </Group>
          {activePage === 0 ? (
            <Stack>
              <TextInput
                withAsterisk
                label="Name"
                description="Name of the processor"
                type="name"
                placeholder="Enter a name"
                {...form.getInputProps('name')}
              />
              <TagsInput
                label="Triggers"
                description="Mime type to trigger the processor"
                placeholder="Add comma separated mime types"
                {...form.getInputProps('triggers')}
              />
            </Stack>
          ) : (
            <Stack>
              <Text size="md" mt="md" fw={500}>
                Choose workflow from an prebuilt templates
              </Text>
              <TextInput
                placeholder="Search for templates"
                leftSection={<IconSearch size={18} />}
              />

              <Box mb="sm" mt="sm">
                <Text size="sm" mb="sm">
                  Suggested templates
                </Text>
                <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="xl">
                  {templates &&
                    templates.length > 0 &&
                    templates.map((t: any) => (
                      <Box onClick={() => selectTemplate(t.key)}>
                        <Template
                          template={t}
                          selected={t.key === form.values.templateKey}
                        />
                      </Box>
                    ))}
                </SimpleGrid>
              </Box>

              {/* <Box mb="sm">
                <Text size="sm" mt="md" mb="sm">
                  Document templates
                </Text>
                <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="xl">
                  {templates &&
                    templates.length > 0 &&
                    [].map((t: any) => (
                      <Box onClick={() => selectTemplate(t.key)}>
                        <Template
                          template={t}
                          selected={t.key === form.values.templateKey}
                        />
                      </Box>
                    ))}
                </SimpleGrid>
              </Box> */}
            </Stack>
          )}
        </form>
      </Container>
    </Box>
  );
};

export default NewprocessorPage;
