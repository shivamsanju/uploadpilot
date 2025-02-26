import {
  ActionIcon,
  Box,
  Button,
  Group,
  Menu,
  Pill,
  Stack,
  Text,
  TextInput,
  ThemeIcon,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { showNotification } from '@mantine/notifications';
import {
  IconActivity,
  IconAlertCircle,
  IconBolt,
  IconCancel,
  IconCircleCheck,
  IconDots,
  IconPlus,
  IconRoute,
  IconSearch,
  IconSettings,
  IconTrash,
} from '@tabler/icons-react';
import { DataTableColumn } from 'mantine-datatable';
import { useCallback, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  useDeleteProcessorMutation,
  useEnableDisableProcessorMutation,
  useGetProcessors,
} from '../../apis/processors';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import { UploadPilotDataTable } from '../../components/Table/Table';
import { timeAgo } from '../../utils/datetime';

const ProcessorList = () => {
  const { width } = useViewportSize();
  const [selectedRecords, setSelectedRecords] = useState<any[]>([]);
  const { workspaceId } = useParams();
  const { isPending, error, processors, invalidate } = useGetProcessors(
    workspaceId || '',
  );
  const { mutateAsync, isPending: isDeleting } = useDeleteProcessorMutation();
  const { mutateAsync: enableDisableProcessor, isPending: isEnabling } =
    useEnableDisableProcessorMutation();
  const navigate = useNavigate();

  const handleRemoveProcessor = useCallback(
    async (processorId: string) => {
      if (
        !workspaceId ||
        workspaceId === '' ||
        !processorId ||
        processorId === ''
      ) {
        showNotification({
          color: 'red',
          title: 'Error',
          message: 'Workspace ID or Processor ID is not available',
        });
        return;
      }

      try {
        await mutateAsync({ workspaceId, processorId });
      } catch (error) {
        console.error(error);
      }
    },
    [workspaceId, mutateAsync],
  );

  const handleRemoveProcessorWithConfirmation = useCallback(
    (processorId: string) => {
      showConfirmationPopup({
        message:
          'Are you sure you want to delete this processor? this is irreversible.',
        onOk: async () => {
          await handleRemoveProcessor(processorId);
        },
      });
    },
    [handleRemoveProcessor],
  );
  const handleBulkRemove = useCallback(async () => {
    if (selectedRecords.length > 0) {
      showConfirmationPopup({
        message:
          'Are you sure you want to delete these processors? this is irreversible.',
        onOk: async () => {
          await Promise.all(
            selectedRecords.map(record => handleRemoveProcessor(record.id)),
          );
          setSelectedRecords([]);
        },
      });
    }
  }, [selectedRecords, handleRemoveProcessor]);

  const handleEnableDisableProcessor = useCallback(
    async (processorId: string, enabled: boolean) => {
      if (
        !workspaceId ||
        workspaceId === '' ||
        !processorId ||
        processorId === ''
      ) {
        showNotification({
          color: 'red',
          title: 'Error',
          message: 'Workspace ID or Processor ID is not available',
        });
        return;
      }

      try {
        await enableDisableProcessor({ workspaceId, processorId, enabled });
      } catch (error) {
        console.error(error);
      }
    },
    [workspaceId, enableDisableProcessor],
  );

  const handleBulkEnableDisable = useCallback(
    async (enabled: boolean) => {
      if (selectedRecords.length > 0) {
        await Promise.all(
          selectedRecords.map(record =>
            handleEnableDisableProcessor(record.id, enabled),
          ),
        );
        setSelectedRecords([]);
      }
    },
    [selectedRecords, handleEnableDisableProcessor],
  );

  const columns: DataTableColumn[] = useMemo(
    () => [
      {
        accessor: 'id',
        title: '',
        width: 20,
        render: (item: any) => (
          <Stack justify="center">
            <IconRoute size={16} stroke={1.5} />
          </Stack>
        ),
      },
      {
        accessor: 'name',
        title: 'Name',
        render: (item: any) => (
          <Group gap="sm">
            {/* <ThemeIcon
              size={20}
              radius={20}
              variant="light"
              color={item?.enabled ? "appcolor" : "gray"}
            >
              {item?.enabled ? (
                <IconRoute size={16} />
              ) : (
                <IconRouteOff size={16} />
              )}
            </ThemeIcon> */}
            <div>{item.name}</div>
          </Group>
        ),
      },
      {
        accessor: 'triggers',
        title: 'Triggers',
        hidden: width < 768,
        render: (item: any) => (
          <div>
            <Text fz="sm" fw={500}>
              {item.triggers && item.triggers.length > 0 ? (
                <>
                  {item.triggers
                    .slice(0, 5)
                    .map((trigger: any, index: number) => (
                      <Pill mr="xs" size="xs" key={index}>
                        {trigger}
                      </Pill>
                    ))}
                  {item.triggers.length > 5 && (
                    <Pill size="xs" color="blue">
                      +{item.triggers.length - 5}
                    </Pill>
                  )}
                </>
              ) : (
                'No Triggers'
              )}
            </Text>
            {/* <Text c="dimmed" fz="xs">
              Triggers
            </Text> */}
          </div>
        ),
      },
      {
        title: 'Updated At',
        accessor: 'updatedAt',
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">
              {params?.updatedAt && timeAgo.format(new Date(params?.updatedAt))}
            </Text>
            {/* <Text fz="xs" c="dimmed">
              Last Updated
            </Text> */}
          </>
        ),
      },
      {
        accessor: 'enabled',
        title: 'Status',
        hidden: width < 768,
        render: (item: any) => (
          <Group align="center" gap="0">
            <ThemeIcon variant="subtle" color={item.enabled ? 'green' : 'red'}>
              {!item.enabled ? (
                <IconAlertCircle size={18} stroke={1.5} />
              ) : (
                <IconCircleCheck size={18} stroke={1.5} />
              )}
            </ThemeIcon>
            {item.enabled ? 'Enabled' : 'Disabled'}
          </Group>
        ),
      },
      // {
      //     accessor: 'goto',
      //     title: 'goto',
      //     textAlign: 'right',
      //     hidden: width < 768,
      //     render: (item: any) => (
      //         <ActionIcon
      //             variant="light"
      //             size="lg"
      //             onClick={() => navigate(`/workspace/${workspaceId}/processors/${item?.id}`)}
      //         >
      //             <IconChevronRight />
      //         </ActionIcon>
      //     ),
      // },
      {
        accessor: 'actions',
        title: 'Actions',
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
                  leftSection={<IconActivity size={16} stroke={1.5} />}
                  disabled={!item?.enabled}
                  onClick={() =>
                    navigate(
                      `/workspace/${workspaceId}/processors/${item?.id}/workflow`,
                    )
                  }
                >
                  <Text>Workflow</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconBolt size={16} stroke={1.5} />}
                  disabled={!item?.enabled}
                  onClick={() =>
                    navigate(
                      `/workspace/${workspaceId}/processors/${item?.id}/runs`,
                    )
                  }
                >
                  <Text>View Runs</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconSettings size={16} stroke={1.5} />}
                  disabled={!item?.enabled}
                  onClick={() =>
                    navigate(
                      `/workspace/${workspaceId}/processors/${item?.id}/settings`,
                    )
                  }
                >
                  <Text>Settings</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={
                    item?.enabled ? (
                      <IconCancel size={16} stroke={1.5} />
                    ) : (
                      <IconCircleCheck size={16} stroke={1.5} />
                    )
                  }
                  onClick={() =>
                    handleEnableDisableProcessor(item.id, !item?.enabled)
                  }
                >
                  <Text>{item?.enabled ? 'Disable' : 'Enable'}</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconTrash size={16} stroke={1.5} />}
                  color="red"
                  onClick={() => handleRemoveProcessorWithConfirmation(item.id)}
                >
                  <Text>Delete</Text>
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ],
    [
      handleRemoveProcessorWithConfirmation,
      handleEnableDisableProcessor,
      width,
      navigate,
      workspaceId,
    ],
  );

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  return (
    <Box mr="md">
      <ContainerOverlay visible={isPending} />
      <UploadPilotDataTable
        fetching={isDeleting || isEnabling}
        minHeight={500}
        columns={columns}
        records={processors}
        verticalSpacing="xs"
        horizontalSpacing="md"
        noHeader={false}
        noRecordsText="No processors found"
        highlightOnHover
        menuBar={
          <Group gap="sm" align="center" justify="space-between">
            <Group gap="sm">
              <Button
                variant="subtle"
                onClick={() =>
                  navigate(`/workspace/${workspaceId}/processors/new`)
                }
                leftSection={<IconPlus size={18} />}
              >
                Add
              </Button>
              <Button
                variant="subtle"
                leftSection={<IconCircleCheck size={18} />}
                onClick={() => handleBulkEnableDisable(true)}
              >
                Enable
              </Button>
              <Button
                variant="subtle"
                leftSection={<IconCancel size={18} />}
                onClick={() => handleBulkEnableDisable(false)}
              >
                Disable
              </Button>
              <Button
                variant="subtle"
                leftSection={<IconTrash size={18} />}
                onClick={handleBulkRemove}
              >
                Delete
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
        onRowDoubleClick={row => {
          navigate(
            `/workspace/${workspaceId}/processors/${row?.record?.id}/runs`,
          );
        }}
        selectedRecords={selectedRecords}
        onSelectedRecordsChange={setSelectedRecords}
      />
    </Box>
  );
};

export default ProcessorList;
