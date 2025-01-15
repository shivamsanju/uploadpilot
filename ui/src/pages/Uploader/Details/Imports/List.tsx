import { useMemo, useState, } from 'react';
import { Title, Menu, Group, Box, Badge, Button, Loader, Tooltip, Text } from '@mantine/core';
import { IconCircleCheck, IconDots, IconExclamationCircle, IconFile, IconBraces, IconRefresh } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { useGetImports } from '../../../../apis/import';
import { ThemedAgGridReact } from '../../../../components/AgGrid/AgGrid';
import { timeAgo } from '../../../../utils/datetime';
import { useQueryClient } from '@tanstack/react-query';
import MetadataLogsModal from './MetadataLogs';
import { formatBytes } from '../../../../utils/utility';


const UploadersListPage = () => {
    const [openModal, setOpenModal] = useState(false)
    const [modalVariant, setModalVariant] = useState<'logs' | 'metadata'>('logs')
    const [logs, setLogs] = useState([])
    const [metadata, setMetadata] = useState({})

    const { uploaderId } = useParams();
    const queryClient = useQueryClient();


    const { isPending, error, imports } = useGetImports(uploaderId as string);

    const handleViewLogs = (importId: string) => {
        const importItem = imports?.find((item: any) => item.id === importId);
        if (importItem) {
            setOpenModal(true)
            setLogs(importItem.logs || [])
            setModalVariant('logs')
        }
    }

    const handleViewMetadata = (importId: string) => {
        const importItem = imports?.find((item: any) => item.id === importId);
        if (importItem) {
            setOpenModal(true)
            setMetadata(importItem.metadata || {})
            setModalVariant('metadata')
        }
    }

    const handleRefresh = () => {
        queryClient.invalidateQueries({ queryKey: ['imports', uploaderId] });
    }

    const colDefs = useMemo(() => {
        return [
            {
                headerName: 'Name',
                field: 'metadata.filename',
            },
            {
                headerName: 'File Type',
                field: 'metadata.filetype',
                cellRenderer: (params: any) => <Badge size="xs" p="sm" variant='default'>{params.value}</Badge>
            },
            {
                headerName: 'Size',
                field: 'size',
                cellRenderer: (params: any) => formatBytes(Number(params.value))
            },
            {
                headerName: 'Stored File Name',
                field: 'storedFileName',
                flex: 1.4
            },
            { headerName: 'Started At', field: 'startedAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            { headerName: 'Finished At', field: 'finishedAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            {
                headerName: 'Status',
                field: 'status',
                flex: 0.7,
                cellStyle: {
                    textAlign: 'center'
                },
                cellRenderer: (params: any) => (
                    < Tooltip label={params.value} >
                        <div>
                            {params.value === 'Success' && <IconCircleCheck size={18} style={{ color: 'green' }} />}
                            {params.value === 'Failed' && <IconExclamationCircle size={18} style={{ color: 'red' }} />}
                            {params.value !== 'Success' && params.value !== 'Failed' && <Loader size={18} />}
                        </div>
                    </Tooltip >

                )

            },
            {
                headerName: 'Actions',
                field: 'actions',
                flex: 0.4,
                cellRenderer: (params: any) => (
                    <Menu>
                        <Menu.Target>
                            <IconDots size={15} style={{ cursor: 'pointer' }} />
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => handleViewLogs(params.data.id)}
                                leftSection={<IconFile size={18} />}
                            >
                                View Logs
                            </Menu.Item>
                            <Menu.Item
                                onClick={() => handleViewMetadata(params.data.id)}
                                leftSection={<IconBraces size={18} />}
                            >
                                View Metadata
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>
                )
            },
        ]
    }, [imports]);


    return (
        <>
            <MetadataLogsModal
                open={openModal}
                onClose={() => setOpenModal(false)}
                variant={modalVariant}
                logs={logs}
                metadata={metadata}
            />
            <Group justify='space-between' align='flex-start' gap="lg">
                {/* <Title order={3} opacity={0.8}>Imports</Title> */}
                <Box />
                <Button variant="default" onClick={handleRefresh} leftSection={<IconRefresh size={18} />}>Refresh</Button>
            </Group>
            <Box h="71.5vh">
                <ThemedAgGridReact
                    overlayNoRowsTemplate='No imports found'
                    loading={isPending && !error}
                    rowData={imports}
                    columnDefs={colDefs}
                    defaultColDef={{
                        sortable: true,
                        filter: true,
                        flex: 1,
                    }}
                    pagination={true}
                    paginationAutoPageSize={true}
                    paginationPageSizeSelector={false}
                />
            </Box>
        </>
    );
}

export default UploadersListPage;