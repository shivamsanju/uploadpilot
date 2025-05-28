import {
  ActionIcon,
  Box,
  Button,
  Group,
  Menu,
  Modal,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import {
  IconBolt,
  IconCancel,
  IconDots,
  IconDownload,
  IconLogs,
  IconSearch,
} from '@tabler/icons-react';
import { DataTableColumn } from 'mantine-datatable';
import { useCallback, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  useCancelWorkflowRun,
  useDownloadRunArtifacts,
  useGetProcessorRuns,
} from '../../apis/processors';
import { useTriggerProcessUpload } from '../../apis/upload';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import { UploadPilotDataTable } from '../../components/Table/Table';
import { formatMilliseconds } from '../../utils/datetime';
import { LogsModal } from './Logs';
import { WorkflowStatus } from './Status';

const ProcessorRunsList = () => {
  const [workflowId, setWorkflowId] = useState<string>('');
  const [runId, setRunId] = useState<string>('');
  const [finishedRun, setFinishedRun] = useState<boolean>(false);
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

  const { mutateAsync: downloadRunArtifacts } = useDownloadRunArtifacts(
    workspaceId || '',
    processorId || '',
  );

  const { mutateAsync: cancelWorkflowRun, isPending: isCancelling } =
    useCancelWorkflowRun(workspaceId || '', processorId || '');

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
      message: 'Are you sure you want to trigger processing for these uploads?',
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

  const handleViewLogs = useCallback(
    (runId: string, workflowId: string, finishedRun = false) => {
      setFinishedRun(finishedRun);
      setRunId(runId);
      setWorkflowId(workflowId);
      setOpened(true);
    },
    [],
  );

  const closeLogs = useCallback(() => {
    setRunId('');
    setWorkflowId('');
    setFinishedRun(false);
    setOpened(false);
  }, []);

  const handleDownloadArtifacts = useCallback(
    async (uploadId: string, runId: string) => {
      try {
        const url = await downloadRunArtifacts({ uploadId, runId });
        window.open(url, '_blank');
      } catch (error) {
        console.error('Error downloading artifacts:', error);
      }
    },
    [downloadRunArtifacts],
  );

  const handleCancelRun = useCallback(
    async (runId: string, workflowId: string) => {
      try {
        showConfirmationPopup({
          message: 'Are you sure you want to cancel this run?',
          onOk: async () => {
            await cancelWorkflowRun({ runId, workflowId });
            setTimeout(() => {
              invalidate();
            }, 2000);
          },
        });
      } catch (error) {
        console.error('Error cancelling run:', error);
      }
    },
    [cancelWorkflowRun, invalidate],
  );

  const columns: DataTableColumn[] = useMemo(
    () => [
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
      {
        title: 'Execution Time',
        accessor: 'executionTimeMillis',
        hidden: width < 768,
        render: (item: any) =>
          formatMilliseconds(item?.executionTimeMillis || 0),
      },
      {
        accessor: 'status',
        title: 'Status',
        render: (item: any) => (
          <Group align="center" gap="0">
            <WorkflowStatus status={item?.status} />
            {item?.status}
          </Group>
        ),
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
                {item.status === 'Completed' && (
                  <Menu.Item
                    leftSection={<IconDownload size={16} stroke={1.5} />}
                    onClick={() =>
                      handleDownloadArtifacts(item?.uploadId, item?.runId)
                    }
                  >
                    <Text>Download Artifacts</Text>
                  </Menu.Item>
                )}
                <Menu.Item
                  leftSection={<IconLogs size={16} stroke={1.5} />}
                  onClick={() =>
                    handleViewLogs(
                      item?.runId,
                      item?.workflowId,
                      item?.record?.status !== 'Running' &&
                        item?.record?.status !== 'Continued-As-New',
                    )
                  }
                >
                  <Text>View Logs</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconBolt size={16} stroke={1.5} />}
                  onClick={() => processUpload(item?.uploadId)}
                >
                  <Text>Re Trigger</Text>
                </Menu.Item>
                {item?.status === 'Running' && (
                  <Menu.Item
                    color="red"
                    leftSection={<IconCancel size={16} stroke={1.5} />}
                    onClick={() =>
                      handleCancelRun(item?.runId, item?.workflowId)
                    }
                  >
                    <Text>Cancel Run</Text>
                  </Menu.Item>
                )}
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ],
    [
      width,
      handleViewLogs,
      processUpload,
      handleCancelRun,
      handleDownloadArtifacts,
    ],
  );

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  return (
    <Box>
      <ContainerOverlay
        visible={isPending || isTriggeringProcess || isCancelling}
      />
      <UploadPilotDataTable
        minHeight="75vh"
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
        onRowDoubleClick={item =>
          handleViewLogs(
            item.record?.runId as string,
            item?.record?.workflowId as string,
            item?.record?.status !== 'Running' &&
              item?.record?.status !== 'Continued-As-New',
          )
        }
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
              finishedRun={finishedRun}
            />
          </div>
        )}
      </Modal>
    </Box>
  );
};

export default ProcessorRunsList;
