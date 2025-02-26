import {
  ActionIcon,
  Box,
  Button,
  Group,
  Menu,
  Stack,
  Text,
  TextInput,
  ThemeIcon,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { showNotification } from '@mantine/notifications';
import {
  IconAlertCircle,
  IconCancel,
  IconCircleCheck,
  IconDots,
  IconKey,
  IconPlus,
  IconSearch,
} from '@tabler/icons-react';
import { DataTableColumn } from 'mantine-datatable';
import { useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useGetApiKeys, useRevokeApiKeyMutation } from '../../apis/apikeys';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import { UploadPilotDataTable } from '../../components/Table/Table';

const WorkspaceApiKeyList = () => {
  const { width } = useViewportSize();
  const { isPending, error, apikeys, invalidate } = useGetApiKeys();
  const { mutateAsync, isPending: revoking } = useRevokeApiKeyMutation();
  const navigate = useNavigate();

  const revokeApiKey = useCallback(
    async (id: string) => {
      if (!id || id === '') {
        showNotification({
          color: 'red',
          title: 'Error',
          message: 'ID is not available',
        });
        return;
      }

      showConfirmationPopup({
        message:
          'Are you sure you want to revoke this api key? this is irreversible.',
        onOk: async () => {
          try {
            await mutateAsync({ id });
          } catch (error) {
            console.error(error);
          }
        },
      });
    },
    [mutateAsync],
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
        title: 'Owner',
        hidden: width < 768,
      },
      {
        accessor: 'revoked',
        title: 'Status',
        render: (item: any) => {
          const expired = new Date(item.expiresAt).getTime() < Date.now();
          return (
            <Group align="center" gap="0">
              <ThemeIcon
                variant="subtle"
                color={item.revoked || expired ? 'red' : 'green'}
              >
                {item.revoked || expired ? (
                  <IconAlertCircle size={18} stroke={1.5} />
                ) : (
                  <IconCircleCheck size={18} stroke={1.5} />
                )}
              </ThemeIcon>
              {item.revoked ? 'Revoked' : expired ? 'Expired' : 'Active'}
            </Group>
            // <Badge
            //   variant="subtle"
            //   color={item.revoked || expired ? 'red' : 'green'}
            // >
            //   {item.revoked || expired
            //     ? 'Revoked'
            //     : expired
            //       ? 'Expired'
            //       : 'Active'}
            // </Badge>
          );
        },
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
                onClick={() => navigate('/api-keys/new')}
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
    </Box>
  );
};

export default WorkspaceApiKeyList;
