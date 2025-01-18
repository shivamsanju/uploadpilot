import { useCallback, useMemo, useState, } from 'react';
import { Menu, Badge, Loader, Tooltip } from '@mantine/core';
import { IconCircleCheck, IconDots, IconExclamationCircle, IconFile, IconBraces } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { useGetImports } from '../../../../apis/import';
import { timeAgo } from '../../../../utils/datetime';
import MetadataLogsModal from './MetadataLogs';
import { formatBytes } from '../../../../utils/utility';
import { UploadPilotDataTable, useUploadPilotDataTable } from '../../../../components/Table/Table';
import ErrorCard from '../../../../components/ErrorCard/ErrorCard';
import { DataTableColumn } from 'mantine-datatable';

const recordsPerPage = 11;

const ImportsList = () => {
    const [openModal, setOpenModal] = useState(false)
    const [modalVariant, setModalVariant] = useState<'logs' | 'metadata'>('logs')
    const [logs, setLogs] = useState([])
    const [metadata, setMetadata] = useState({})
    const { uploaderId } = useParams();

    const { searchFilter, onSearchFilterChange, page, onPageChange } = useUploadPilotDataTable();

    const { isPending, error, imports, totalRecords, invalidate } = useGetImports({
        uploaderId: uploaderId || '',
        skip: (page - 1) * recordsPerPage,
        limit: recordsPerPage,
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


    const colDefs: DataTableColumn[] = useMemo(() => {
        return [
            {
                title: 'Name',
                accessor: 'metadata.filename',
            },
            {
                title: 'File Type',
                accessor: 'metadata.filetype',
                textAlign: 'center',
                render: (params: any) => <Badge size="xs" p="sm" variant='default'>{params?.metadata?.filetype}</Badge>
            },
            {
                title: 'Size',
                accessor: 'size',
                textAlign: 'center',
                render: (params: any) => formatBytes(Number(params?.size))
            },
            {
                title: 'Stored File Name',
                accessor: 'storedFileName',
                textAlign: 'center',
                flex: 1.4
            },
            {
                title: 'Started At',
                accessor: 'startedAt',
                textAlign: 'center',
                render: (params: any) => params.startedAt && timeAgo.format(new Date(params.startedAt))
            },
            {
                title: 'Finished At',
                accessor: 'finishedAt',
                textAlign: 'center',
                render: (params: any) => params.finishedAt && timeAgo.format(new Date(params.finishedAt))
            },
            {
                title: 'Status',
                accessor: 'status',
                flex: 0.7,
                textAlign: 'center',
                render: (params: any) => (
                    <Tooltip label={params.status} >
                        <div>
                            {params.status === 'Success' && <IconCircleCheck size={18} style={{ color: 'green' }} />}
                            {params.status === 'Failed' && <IconExclamationCircle size={18} style={{ color: 'red' }} />}
                            {params.status !== 'Success' && params.status !== 'Failed' && <Loader size={18} />}
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
                    <Menu>
                        <Menu.Target>
                            <IconDots size={15} style={{ cursor: 'pointer' }} />
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => handleViewLogs(params.id)}
                                leftSection={<IconFile size={18} />}
                            >
                                View Logs
                            </Menu.Item>
                            <Menu.Item
                                onClick={() => handleViewMetadata(params.id)}
                                leftSection={<IconBraces size={18} />}
                            >
                                View Metadata
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>
                )
            },
        ]
    }, [handleViewLogs, handleViewMetadata]);


    return (
        <>
            <MetadataLogsModal
                open={openModal}
                onClose={() => setOpenModal(false)}
                variant={modalVariant}
                logs={logs}
                metadata={metadata}
            />
            {error ?
                <ErrorCard title={error.name} message={error.message} h="65vh" /> :
                <UploadPilotDataTable
                    fetching={isPending}
                    showSearch
                    onSearchFilterChange={onSearchFilterChange}
                    showRefresh
                    onRefresh={invalidate}
                    height={"74vh"}
                    verticalSpacing={"sm"}
                    columns={colDefs}
                    records={imports}
                    page={page}
                    onPageChange={onPageChange}
                    recordsPerPage={recordsPerPage}
                    totalRecords={totalRecords}
                />
            }
        </>
    );
}

export default ImportsList;