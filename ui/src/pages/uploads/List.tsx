import {
  ActionIcon,
  Box,
  Button,
  Group,
  Menu,
  MultiSelect,
  Stack,
  Text,
  TextInput,
} from '@mantine/core';
import { useDebouncedValue, useViewportSize } from '@mantine/hooks';
import {
  IconBolt,
  IconBraces,
  IconChevronsDown,
  IconCopy,
  IconDots,
  IconDownload,
  IconSearch,
} from '@tabler/icons-react';
import { DataTableColumn } from 'mantine-datatable';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  useDownloadUploadedFile,
  useGetUploads,
  useTriggerProcessUpload,
} from '../../apis/upload';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { ContainerOverlay } from '../../components/Overlay';
import { showConfirmationPopup } from '../../components/Popups/ConfirmPopup';
import { UploadPilotDataTable } from '../../components/Table/Table';
import { timeAgo } from '../../utils/datetime';
import { getFileIcon } from '../../utils/fileicons';
import { formatBytes } from '../../utils/utility';
import { MetadataModal } from './Metadata';
import { UploadStatus } from './Status';

const batchSize = 20;

const UploadList = ({ setTotalRecords }: any) => {
  const scrollViewportRef = useRef<HTMLDivElement>(null);
  const { width } = useViewportSize();
  const [openModal, setOpenModal] = useState(false);
  const [modalVariant, setModalVariant] = useState<'logs' | 'metadata'>('logs');
  const [metadata, setMetadata] = useState({});
  const [selectedRecords, setSelectedRecords] = useState<any[]>([]);
  const [search, setSearch] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string[]>([]);
  const [debouncedSearch] = useDebouncedValue(search, 1000);
  const [debouncedStatusFilter] = useDebouncedValue(statusFilter, 1000);

  const { workspaceId } = useParams();

  const {
    isPending,
    error,
    isFetchNextPageError,
    uploads,
    fetchNextPage,
    totalRecords,
    isFetchingNextPage,
    invalidate,
    hasNextPage,
  } = useGetUploads({
    workspaceId: workspaceId || '',
    batchSize,
    search: debouncedSearch,
    filter: {
      status: debouncedStatusFilter,
    },
  });

  const { mutateAsync: downloadFile } = useDownloadUploadedFile(
    workspaceId || '',
  );

  const { mutateAsync: triggerProcessUpload, isPending: isTriggeringProcess } =
    useTriggerProcessUpload(workspaceId || '');

  const getFileUrl = useCallback(
    async (uploadId: string) => {
      if (!workspaceId) {
        return;
      }
      try {
        const url = await downloadFile({
          uploadId,
        });
        return url;
      } catch (error) {
        console.error('Error downloading file:', error);
      }
    },
    [downloadFile, workspaceId],
  );

  const processUpload = useCallback(
    async (uploadId: string) => {
      try {
        await triggerProcessUpload({ uploadId });
      } catch (error) {
        console.error('Error processing upload:', error);
      }
    },
    [triggerProcessUpload],
  );

  const handleBulkProcess = useCallback(async () => {
    showConfirmationPopup({
      message: 'Are you sure you want to start processing for these uploads?',
      onOk: async () => {
        try {
          await Promise.all(
            selectedRecords.map(record => processUpload(record.id)),
          );
          setSelectedRecords([]);
        } catch (error) {
          console.error('Error processing upload:', error);
        }
      },
    });
  }, [selectedRecords, processUpload]);

  const handleDownload = useCallback(
    async (uploadId: string) => {
      const url = await getFileUrl(uploadId);
      window.open(url, '_blank');
    },
    [getFileUrl],
  );

  const handleCopyLink = useCallback(
    async (uploadId: string) => {
      const url = await getFileUrl(uploadId);
      navigator.clipboard.writeText(url);
    },
    [getFileUrl],
  );

  const handleViewMetadata = useCallback((metadata: any) => {
    setOpenModal(true);
    setMetadata(metadata || {});
    setModalVariant('metadata');
  }, []);

  const handleRefresh = () => {
    invalidate();
    scrollViewportRef.current?.scrollTo(0, 0);
  };

  const colDefs: DataTableColumn[] = useMemo(() => {
    return [
      {
        accessor: 'id',
        title: '',
        width: 20,
        render: (params: any) => (
          <Stack justify="center">{getFileIcon(params?.fileType, 16)}</Stack>
        ),
      },
      {
        accessor: 'id',
        title: 'Upload ID',
      },
      {
        title: 'Name',
        accessor: 'fileName',
        elipsis: true,
        render: (params: any) => <Text fz="sm">{params?.fileName}</Text>,
      },
      {
        title: 'Content Type',
        accessor: 'contentType',
        hidden: width < 768,
        render: (params: any) => <Text fz="sm">{params?.contentType}</Text>,
      },
      {
        title: 'Size',
        accessor: 'contentLength',
        hidden: width < 768,
        render: (params: any) => (
          <Text fz="sm">{formatBytes(Number(params?.contentLength))}</Text>
        ),
      },
      {
        title: 'Started At',
        accessor: 'startedAt',
        hidden: width < 768,
        render: (params: any) => (
          <Text fz="sm">
            {params?.startedAt && !params.startedAt.startsWith('0001-01-01')
              ? timeAgo.format(new Date(params?.startedAt))
              : '-'}
          </Text>
        ),
      },
      {
        title: 'Finished At',
        accessor: 'finishedAt',
        hidden: width < 768,
        render: (params: any) => (
          <Text fz="sm">
            {params?.finishedAt && !params.finishedAt.startsWith('0001-01-01')
              ? timeAgo.format(new Date(params?.finishedAt))
              : '-'}
          </Text>
        ),
      },
      {
        title: 'Status',
        accessor: 'status',
        render: (params: any) => (
          <Group align="center" gap="0">
            <UploadStatus status={params?.status} />
            {params?.status}
          </Group>
        ),
        filter: (
          <MultiSelect
            w={200}
            data={[
              'Finished',
              'Failed',
              'Created',
              'Queued',
              'Skipped',
              'Cancelled',
            ]}
            value={statusFilter}
            placeholder="Filter by status"
            onChange={setStatusFilter}
            comboboxProps={{ withinPortal: false }}
            clearable
            searchable
          />
        ),
        filtering: statusFilter.length > 0,
      },
      {
        title: 'Actions',
        accessor: 'actions',
        textAlign: 'right',
        width: 100,
        render: (params: any) => (
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
                {params?.status === 'Finished' && (
                  <Menu.Item
                    onClick={() => handleDownload(params?.id)}
                    leftSection={<IconDownload size={18} />}
                  >
                    <Text>Download</Text>
                  </Menu.Item>
                )}
                {params?.status === 'Finished' && (
                  <Menu.Item
                    onClick={() => handleCopyLink(params?.id)}
                    leftSection={<IconCopy size={18} />}
                  >
                    <Text>Copy URL</Text>
                  </Menu.Item>
                )}
                <Menu.Item
                  onClick={() => handleViewMetadata(params?.metadata)}
                  leftSection={<IconBraces size={18} />}
                >
                  <Text>View Metadata</Text>
                </Menu.Item>
                {params?.status === 'Finished' && (
                  <Menu.Item
                    onClick={() => processUpload(params?.id)}
                    leftSection={<IconBolt size={18} />}
                  >
                    <Text>Process</Text>
                  </Menu.Item>
                )}
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ];
  }, [
    handleViewMetadata,
    width,
    handleDownload,
    handleCopyLink,
    statusFilter,
    setStatusFilter,
    processUpload,
  ]);

  useEffect(() => {
    setTotalRecords(totalRecords);
  }, [setTotalRecords, totalRecords]);

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="70vh" />;
  }

  if (isFetchNextPageError) {
    return (
      <ErrorCard title="Error" message={'Failed to fetch next page'} h="70vh" />
    );
  }

  return (
    <Box mt="lg">
      {/* Loading overlay only while pending, not on refetch*/}
      <ContainerOverlay visible={isPending || isTriggeringProcess} />{' '}
      <MetadataModal
        open={openModal && modalVariant === 'metadata'}
        onClose={() => setOpenModal(false)}
        metadata={metadata || {}}
      />
      <Box mr="md">
        <UploadPilotDataTable
          minHeight={500}
          verticalSpacing="xs"
          horizontalSpacing="lg"
          noHeader={false}
          onScrollToBottom={fetchNextPage}
          columns={colDefs}
          records={uploads}
          scrollViewportRef={scrollViewportRef}
          noRecordsText="No imports yet"
          selectedRecords={selectedRecords}
          onSelectedRecordsChange={setSelectedRecords}
          menuBar={
            <Group gap="sm" align="center" justify="space-between">
              <Group gap="sm">
                <Button
                  variant="subtle"
                  leftSection={<IconBolt size={18} />}
                  onClick={handleBulkProcess}
                >
                  Process
                </Button>
                <RefreshButton variant="subtle" onClick={handleRefresh} />
              </Group>
              <TextInput
                value={search}
                onChange={e => setSearch(e.target.value)}
                placeholder="Search (min 3 characters)"
                rightSection={<IconSearch size={18} />}
                variant="outline"
              />
            </Group>
          }
        />
      </Box>
      <Stack align="center" justify="center" p="md" key="selectedRecords">
        <Button
          display={hasNextPage ? 'block' : 'none'}
          leftSection={<IconChevronsDown size={16} />}
          variant="subtle"
          disabled={isPending || isFetchingNextPage || !hasNextPage}
          loading={isFetchingNextPage}
          onClick={() => fetchNextPage({})}
        >
          Load More
        </Button>
      </Stack>
    </Box>
  );
};

export default UploadList;
