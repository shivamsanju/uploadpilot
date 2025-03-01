import { CodeHighlight } from '@mantine/code-highlight';
import {
  Box,
  Button,
  Checkbox,
  Container,
  Group,
  Loader,
  MultiSelect,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { useForm } from '@mantine/form';
import { modals } from '@mantine/modals';
import { IconAlertTriangle } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useCreateApiKeyMutation } from '../../../apis/apikeys';
import { useGetWorkspaces } from '../../../apis/workspace';
import { ErrorCard } from '../../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../../components/Overlay';
import { CreateApiKeyData } from '../../../types/apikey';

type Props = {};

const CreateApiKeyPage: React.FC<Props> = () => {
  const { isPending, error, workspaces } = useGetWorkspaces();
  const wsMap = new Map(workspaces?.map((ws: any) => [ws.id, ws.name]));

  const navigate = useNavigate();

  const form = useForm<CreateApiKeyData>({
    initialValues: {
      expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      name: '',
      tenantRead: false,
      tenantManage: false,
      workspacePerms: [],
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

  const { mutateAsync, isPending: isCreating } = useCreateApiKeyMutation();

  const handleAdd = async (data: any) => {
    try {
      const newKey = await mutateAsync({
        data,
      });
      modals.open({
        title: 'Api key created successfully',
        centered: true,
        closeOnEscape: false,
        closeOnClickOutside: false,
        padding: 'lg',
        size: 'xl',
        onClose: () => {
          navigate('/api-keys');
        },
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

  const onChangeWs = (wsIds: string[]) => {
    const oldPermissionsMap = new Map(
      form.values.workspacePerms.map((p: any) => [p.id, p]),
    );

    form.setFieldValue(
      'workspacePerms',
      wsIds.map(
        id =>
          oldPermissionsMap.get(id) || {
            id: id,
            read: false,
            manage: false,
            upload: false,
          },
      ),
    );
  };

  if (error) {
    return <ErrorCard message={error?.message} title={error?.name} />;
  }

  return (
    <Box mb={50}>
      <ContainerOverlay visible={isCreating || isPending} />
      <Container mt="md">
        <form onSubmit={form.onSubmit(handleAdd)}>
          <Group justify="space-between" align="center" mb="lg">
            <Title order={3}>Create new api key</Title>
            <Button type="submit">Create</Button>
          </Group>

          <Stack gap="md">
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
            <Text size="md" fw={500} mt="xl">
              Tenant level permissions
            </Text>
            <Checkbox
              label="Read access"
              {...form.getInputProps('tenantRead')}
            />
            <Checkbox
              label="Manage access"
              {...form.getInputProps('tenantManage')}
            />
            <Text size="md" fw={500} mt="xl">
              Workspace level permissions
            </Text>
            <MultiSelect
              searchable
              leftSection={isPending ? <Loader size="xs" type="oval" /> : null}
              disabled={isPending}
              data={
                workspaces?.map((ws: any) => ({
                  value: ws.id,
                  label: ws.name,
                })) || []
              }
              label="Workspaces"
              placeholder="Select workspaces"
              onChange={onChangeWs}
            />
            <Stack mt="sm">
              <Table striped highlightOnHover>
                <Table.Thead>
                  <tr>
                    <th style={{ textAlign: 'left', padding: '10px' }}>
                      Workspace
                    </th>
                    {['Can upload', 'Can read', 'Can manage'].map(
                      (label, i) => (
                        <th
                          key={i}
                          style={{ textAlign: 'center', padding: '10px' }}
                        >
                          <Checkbox
                            label={label}
                            checked={form.values.workspacePerms.every(
                              (p: any) => p[label.replace('Can ', '')],
                            )}
                            onChange={e => {
                              const checked = e.currentTarget.checked;
                              form.setFieldValue(
                                'workspacePerms',
                                form.values.workspacePerms.map((p: any) => ({
                                  ...p,
                                  [label.replace('Can ', '')]: checked,
                                })),
                              );
                            }}
                          />
                        </th>
                      ),
                    )}
                  </tr>
                </Table.Thead>
                <Table.Tbody>
                  {form.values.workspacePerms.map(
                    (perm: any, index: number) => (
                      <tr key={perm.id}>
                        <td style={{ padding: '10px' }}>
                          {wsMap.get(perm.id) as string}
                        </td>
                        {['upload', 'read', 'manage'].map((field, i) => (
                          <td
                            key={i}
                            style={{ textAlign: 'center', padding: '10px' }}
                          >
                            <Checkbox
                              {...form.getInputProps(
                                `workspacePerms.${index}.${field}`,
                              )}
                              checked={
                                (form.values.workspacePerms[index] as any)[
                                  field
                                ]
                              }
                            />
                          </td>
                        ))}
                      </tr>
                    ),
                  )}
                </Table.Tbody>
              </Table>
            </Stack>
          </Stack>
        </form>
      </Container>
    </Box>
  );
};

export default CreateApiKeyPage;
