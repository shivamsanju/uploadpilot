import { useCallback, useMemo, useRef, useState, } from 'react';
import { Menu, Badge, Loader, Tooltip, Stack, Box, Button, Group, ActionIcon, Text } from '@mantine/core';
import { IconCircleCheck, IconDots, IconExclamationCircle, IconBraces, IconChevronsDown, IconLogs } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { useGetImports } from '../../../apis/import';
import { timeAgo } from '../../../utils/datetime';
import MetadataLogsModal from './MetadataLogs';
import { formatBytes } from '../../../utils/utility';
import { UploadPilotDataTable, useUploadPilotDataTable } from '../../../components/Table/Table';
import { ErrorCard } from '../../../components/ErrorCard/ErrorCard';
import { DataTableColumn } from 'mantine-datatable';

const batchSize = 20;

const ImportsList = () => {
    const scrollViewportRef = useRef<HTMLDivElement>(null);

    const [openModal, setOpenModal] = useState(false)
    const [modalVariant, setModalVariant] = useState<'logs' | 'metadata'>('logs')
    const [logs, setLogs] = useState([])
    const [metadata, setMetadata] = useState({})
    const [selectedRecords, setSelectedRecords] = useState<any[]>([])

    const { workspaceId } = useParams();
    const { searchFilter, onSearchFilterChange } = useUploadPilotDataTable();

    const { isPending, error, isFetchNextPageError, imports, fetchNextPage, isFetchingNextPage, isFetching, invalidate, hasNextPage } = useGetImports({
        workspaceId: workspaceId || '',
        batchSize,
        search: searchFilter
    });

    const handleViewLogs = useCallback((importId: string) => {
        const importItem = imports?.find((item: any) => item.id === importId);
        if (importItem) {
            setOpenModal(true)
            setLogs(importItem.logs || [])
            setModalVariant('logs')
        }
    }, [imports])

    const handleViewMetadata = useCallback((importId: string) => {
        const importItem = imports?.find((item: any) => item.id === importId);
        if (importItem) {
            setOpenModal(true)
            setMetadata(importItem.metadata || {})
            setModalVariant('metadata')
        }
    }, [imports])

    const handleRefresh = () => {
        invalidate();
        scrollViewportRef.current?.scrollTo(0, 0);
    }


    const colDefs: DataTableColumn[] = useMemo(() => {
        return [
            {
                title: 'Name',
                accessor: 'storedFileName',
                elipsis: true,
                width: 500,
                render: (params: any) => (
                    <>
                        <Text fz="sm">{params?.storedFileName}</Text>
                        <Text fz="xs" c="dimmed">
                            Filename
                        </Text>
                    </>
                ),
            },
            {
                title: 'File Type',
                accessor: 'metadata.filetype',
                textAlign: 'center',
                render: (params: any) => (
                    <>
                        <Badge size="xs" p="sm" variant='outline'>{params?.metadata?.filetype}</Badge>
                        <Text fz="xs" c="dimmed">
                            Mime Type
                        </Text>
                    </>
                ),
            },
            {
                title: 'Size',
                accessor: 'size',
                textAlign: 'center',
                render: (params: any) => (
                    <>
                        <Text fz="sm">{formatBytes(Number(params?.size))}</Text>
                        <Text fz="xs" c="dimmed">
                            Size
                        </Text>
                    </>
                )
            },
            {
                title: 'Finished At',
                accessor: 'finishedAt',
                textAlign: 'center',
                render: (params: any) => (
                    <>
                        <Text fz="sm">{params?.finishedAt && timeAgo.format(new Date(params?.finishedAt))}</Text>
                        <Text fz="xs" c="dimmed">
                            Uploaded
                        </Text>
                    </>
                )
            },
            {
                title: 'Status',
                accessor: 'status',
                flex: 0.7,
                textAlign: 'center',
                render: (params: any) => (
                    <Tooltip label={params?.status} >
                        <div>
                            {params?.status === 'Success' && <IconCircleCheck size={18} style={{ color: 'green' }} />}
                            {params?.status === 'Failed' && <IconExclamationCircle size={18} style={{ color: 'red' }} />}
                            {params?.status !== 'Success' && params?.status !== 'Failed' && <Loader size={18} />}
                        </div>
                    </Tooltip >

                )

            },
            {
                title: 'Actions',
                accessor: 'actions',
                flex: 0.4,
                textAlign: 'center',
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
                                <Menu.Item
                                    onClick={() => handleViewLogs(params?.id)}
                                    leftSection={<IconLogs size={18} />}
                                >
                                    View Logs
                                </Menu.Item>
                                <Menu.Item
                                    onClick={() => handleViewMetadata(params?.id)}
                                    leftSection={<IconBraces size={18} />}
                                >
                                    View Metadata
                                </Menu.Item>
                            </Menu.Dropdown>
                        </Menu>
                    </Group>
                )
            },
        ]
    }, [handleViewLogs, handleViewMetadata]);

    if (error) {
        return <ErrorCard title="Error" message={error.message} h="70vh" />
    }

    if (isFetchNextPageError) {
        return <ErrorCard title="Error" message={"Failed to fetch next page"} h="70vh" />
    }


    return (
        <>
            <MetadataLogsModal
                open={openModal}
                onClose={() => setOpenModal(false)}
                variant={modalVariant}
                logs={logs}
                metadata={metadata}
            />
            <Box mr="md">
                <UploadPilotDataTable
                    verticalSpacing="lg"
                    horizontalSpacing="lg"
                    fetching={isPending || isFetchingNextPage || isFetching}
                    noHeader={true}
                    showSearch={true}
                    showExport={true}
                    showRefresh={true}
                    onRefresh={handleRefresh}
                    onSearchFilterChange={onSearchFilterChange}
                    columns={colDefs}
                    records={imports}
                    selectedRecords={selectedRecords}
                    onSelectedRecordsChange={setSelectedRecords}
                    selectionCheckboxProps={{ style: { cursor: 'pointer' } }}
                    onScrollToBottom={fetchNextPage}
                    scrollViewportRef={scrollViewportRef}
                />

                <Stack align='center' justify='center' p="md" key="selectedRecords">
                    <Button
                        display={hasNextPage ? 'block' : 'none'}
                        size="xs"
                        leftSection={<IconChevronsDown size={16} />}
                        variant="subtle"
                        disabled={isPending || isFetchingNextPage || isFetching || !hasNextPage}
                        onClick={() => fetchNextPage({})}
                    >
                        Load More
                    </Button>
                </Stack>
            </Box>
        </>
    );
}

export default ImportsList;