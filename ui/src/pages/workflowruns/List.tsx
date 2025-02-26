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
  Title,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconBolt, IconDots, IconLogs, IconSearch } from '@tabler/icons-react';
import { DataTableColumn } from 'mantine-datatable';
import { useCallback, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useGetProcessorRuns } from '../../apis/processors';
import { useTriggerProcessUpload } from '../../apis/upload';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import { UploadPilotDataTable } from '../../components/Table/Table';
import { formatMilliseconds } from '../../utils/datetime';
import { LogsModal } from './Logs';
import { statusConfig } from './status';

const ProcessorRunsList = () => {
  const [workflowId, setWorkflowId] = useState<string>('');
  const [runId, setRunId] = useState<string>('');
  const [opened, setOpened] = useState(false);
  const [selectedRecords, setSelectedRecords] = useState<any[]>([]);

  const { width } = useViewportSize();
  const { workspaceId, processorId } = useParams();
  const { isPending, error, runs, invalidate } = useGetProcessorRuns(
    workspaceId || '',
    processorId || '',
  );

  const { mutateAsync: triggerProcessUpload, isPending: isTriggeringProcess } =
    useTriggerProcessUpload(workspaceId || '');

  const processUpload = useCallback(
    async (uploadId: string) => {
      try {
        await triggerProcessUpload({ uploadId });
        setTimeout(() => {
          invalidate();
        }, 2000);
      } catch (error) {
        console.error('Error processing upload:', error);
      }
    },
    [triggerProcessUpload, invalidate],
  );

  const handleBulkProcess = async () => {
    showConfirmationPopup({
      message: 'Are you sure you want to start processing for these uploads?',
      onOk: async () => {
        try {
          await Promise.all(
            selectedRecords.map(record =>
              triggerProcessUpload({
                uploadId: record?.uploadId,
              }),
            ),
          );
          setSelectedRecords([]);
          setTimeout(() => {
            invalidate();
          }, 2000);
        } catch (error) {
          console.error('Error processing upload:', error);
        }
      },
    });
  };

  const handleViewLogs = useCallback((runId: string, workflowId: string) => {
    setRunId(runId);
    setWorkflowId(workflowId);
    setOpened(true);
  }, []);

  const closeLogs = useCallback(() => {
    setRunId('');
    setWorkflowId('');
    setOpened(false);
  }, []);

  const columns: DataTableColumn[] = useMemo(
    () => [
      {
        accessor: 'id',
        title: '',
        width: 20,
        render: (item: any) => (
          <Stack justify="center">
            <IconBolt size={16} stroke={1.5} />
          </Stack>
        ),
      },
      {
        accessor: 'status',
        title: 'Status',
        render: (item: any) => (
          <Badge
            variant="outline"
            color={statusConfig[item?.status?.toLowerCase() || '']}
            size="sm"
          >
            {item?.status}
          </Badge>
        ),
      },
      {
        accessor: 'runId',
        title: 'Run ID',
      },
      {
        accessor: 'uploadId',
        title: 'Upload ID',
      },
      {
        title: 'Started At',
        accessor: 'startTime',
        hidden: width < 768,
        render: (item: any) =>
          new Date(item?.startTime).toLocaleString('en-US'),
      },
      {
        title: 'Ended At',
        accessor: 'endTime',
        hidden: width < 768,
        render: (item: any) => {
          if (item?.endTime !== '0001-01-01T00:00:00Z') {
            return new Date(item?.endTime).toLocaleString('en-US');
          }
          return '-';
        },
      },
      // {
      //   title: "Workflow Time",
      //   accessor: "workflowTimeMillis",
      //   hidden: width < 768,
      //   render: (item: any) =>
      //     formatMilliseconds(item?.workflowTimeMillis || 0),
      // },
      {
        title: 'Execution Time',
        accessor: 'executionTimeMillis',
        hidden: width < 768,
        render: (item: any) =>
          formatMilliseconds(item?.executionTimeMillis || 0),
      },
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
                  leftSection={<IconLogs size={16} stroke={1.5} />}
                  onClick={() => handleViewLogs(item?.runId, item?.workflowId)}
                >
                  <Text>View Logs</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconBolt size={16} stroke={1.5} />}
                  onClick={() => processUpload(item?.uploadId)}
                >
                  <Text>Re Trigger</Text>
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ],
    [width, handleViewLogs, processUpload],
  );

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  return (
    <Box>
      <ContainerOverlay visible={isPending || isTriggeringProcess} />
      <UploadPilotDataTable
        minHeight={500}
        columns={columns}
        records={runs}
        verticalSpacing="xs"
        horizontalSpacing="md"
        noHeader={false}
        noRecordsText="No runs yet"
        highlightOnHover
        menuBar={
          <Group gap="sm" align="center" justify="space-between">
            <Group gap="sm">
              <Button
                variant="subtle"
                onClick={handleBulkProcess}
                leftSection={<IconBolt size={18} />}
              >
                Re Trigger
              </Button>

              <RefreshButton onClick={invalidate} />
            </Group>
            <TextInput
              value={''}
              placeholder="Search (Upload ID, Run ID)"
              rightSection={<IconSearch size={18} />}
              variant="subtle"
            />
          </Group>
        }
        selectedRecords={selectedRecords}
        onSelectedRecordsChange={setSelectedRecords}
      />
      <Modal
        padding="xl"
        transitionProps={{ transition: 'pop' }}
        opened={opened}
        onClose={() => {
          setOpened(false);
        }}
        title={
          <Title order={5} opacity={0.7}>
            Logs
          </Title>
        }
        closeOnClickOutside={false}
        size="lg"
      >
        {workspaceId && processorId && runId && workflowId && (
          <div>
            <LogsModal
              open={opened}
              onClose={closeLogs}
              workspaceId={workspaceId}
              processorId={processorId}
              workflowId={workflowId}
              runId={runId}
            />
          </div>
        )}
      </Modal>
    </Box>
  );
};

export default ProcessorRunsList;
