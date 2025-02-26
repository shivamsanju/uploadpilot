import {
  ActionIcon,
  Badge,
  Box,
  Button,
  Group,
  Menu,
  Modal,
  Stack,
  Text,
  TextInput,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { showNotification } from '@mantine/notifications';
import {
  IconCancel,
  IconDots,
  IconKey,
  IconPlus,
  IconSearch,
} from '@tabler/icons-react';
import { DataTableColumn } from 'mantine-datatable';
import { useCallback, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  useGetApiKeysInWorkspace,
  useRevokeApiKeyMutation,
} from '../../apis/apikeys';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import { UploadPilotDataTable } from '../../components/Table/Table';
import CreateApiKeyForm from './add';

const WorkspaceApiKeyList = () => {
  const { width } = useViewportSize();
  const [opened, setOpened] = useState(false);
  const { workspaceId } = useParams();
  const { isPending, error, apikeys, invalidate } = useGetApiKeysInWorkspace(
    workspaceId || '',
  );
  const { mutateAsync, isPending: revoking } = useRevokeApiKeyMutation();

  const revokeApiKey = useCallback(
    async (id: string) => {
      if (!workspaceId || workspaceId === '' || !id || id === '') {
        showNotification({
          color: 'red',
          title: 'Error',
          message: 'Workspace ID or ID is not available',
        });
        return;
      }

      showConfirmationPopup({
        message:
          'Are you sure you want to revoke this api key? this is irreversible.',
        onOk: async () => {
          try {
            await mutateAsync({ workspaceId, id });
          } catch (error) {
            console.error(error);
          }
        },
      });
    },
    [workspaceId, mutateAsync],
  );

  const columns: DataTableColumn[] = useMemo(
    () => [
      {
        accessor: 'id',
        title: '',
        width: 20,
        render: (item: any) => (
          <Stack justify="center">
            <IconKey size={16} stroke={1.5} />
          </Stack>
        ),
      },
      {
        accessor: 'name',
        title: 'Name',
        width: '30%',
      },
      {
        accessor: 'expiresAt',
        title: 'Expires At',
        render: (item: any) => new Date(item.expiresAt).toLocaleDateString(),
      },
      {
        accessor: 'createdAt',
        title: 'Created At',
        render: (item: any) => new Date(item.createdAt).toLocaleString(),
      },
      {
        accessor: 'createdBy',
        title: 'Created By',
        hidden: width < 768,
      },
      {
        accessor: 'revoked',
        title: 'Status',
        render: (item: any) => (
          <Badge color={item.revoked ? 'red' : 'green'} variant="outline">
            {item.revoked ? 'Revoked' : 'Active'}
          </Badge>
        ),
      },
      {
        accessor: 'actions',
        title: 'Actions',
        width: 100,
        textAlign: 'right',
        render: (item: any) => (
          <Group gap={0} justify="flex-end">
            <Menu
              transitionProps={{ transition: 'pop' }}
              withArrow
              position="bottom-end"
              withinPortal
            >
              <Menu.Target>
                <ActionIcon variant="subtle" color="dimmed">
                  <IconDots size={16} stroke={1.5} />
                </ActionIcon>
              </Menu.Target>
              <Menu.Dropdown>
                <Menu.Item
                  disabled={item.revoked}
                  leftSection={<IconCancel size={16} stroke={1.5} />}
                  color="red"
                  onClick={() => revokeApiKey(item.id)}
                >
                  <Text>Revoke</Text>
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ],
    [revokeApiKey, width],
  );

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  return (
    <Box mr="md">
      <ContainerOverlay visible={isPending} />
      <UploadPilotDataTable
        fetching={revoking}
        minHeight={500}
        columns={columns}
        records={apikeys}
        verticalSpacing="xs"
        horizontalSpacing="md"
        noRecordsText="No api key created"
        menuBar={
          <Group gap="sm" align="center" justify="space-between">
            <Group gap="sm">
              <Button
                variant="subtle"
                onClick={() => setOpened(true)}
                leftSection={<IconPlus size={18} />}
              >
                Add
              </Button>
              <RefreshButton onClick={invalidate} />
            </Group>
            <TextInput
              value={''}
              // onChange={(e) => onSearchFilterChange(e.target.value)}
              placeholder="Search by name or status"
              leftSection={<IconSearch size={18} />}
              variant="subtle"
            />
          </Group>
        }
      />
      <Modal
        centered
        padding="xl"
        transitionProps={{ transition: 'pop' }}
        opened={opened}
        onClose={() => setOpened(false)}
        title="Create APIKey"
        closeOnClickOutside={false}
        size="lg"
      >
        <CreateApiKeyForm
          setOpened={setOpened}
          workspaceId={workspaceId || ''}
        />
      </Modal>
    </Box>
  );
};

export default WorkspaceApiKeyList;
