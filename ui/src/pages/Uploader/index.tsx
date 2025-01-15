import { useMemo, } from 'react';
import { Button, Title, Menu, Group, Box, Anchor } from '@mantine/core';
import { IconCirclePlus2, IconDots, IconEye } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useGetUploaders } from '../../apis/uploader';
import { ThemedAgGridReact } from '../../components/AgGrid/AgGrid';
import { timeAgo } from '../../utils/datetime';


const UploadersListPage = () => {
    const navigate = useNavigate();

    const { isPending, error, uploaders } = useGetUploaders();

    const handleCreateNewUploader = () => {
        navigate('/uploaders/new');
    }

    const colDefs = useMemo(() => {
        return [
            {
                headerName: 'Name',
                field: 'name',
                cellRenderer: (params: any) => <Anchor size="sm" href={`/uploaders/${params.data.id}`}>{params.value}</Anchor>,
                flex: 1.7
            },
            {
                headerName: 'Created By', field: 'createdBy',
                filter: 'agSetColumnFilter',
                filterParams: {
                    suppressCloseOnClickOutside: true, // Keeps the filter open
                },
                editable: true,

            },
            { headerName: 'Created At', field: 'createdAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            { headerName: 'Updated At', field: 'updatedAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
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
                                onClick={() => navigate(`/uploaders/${params.data.id}`)}
                                leftSection={<IconEye size={15} />}
                            >
                                View
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>
                )
            },
        ]
    }, [navigate]);


    return (
        <>
            <Group justify='space-between' align='flex-start' gap="lg">
                <Title order={3} mb="lg" opacity={0.8}>Uploaders</Title>
                <Button size="xs" leftSection={<IconCirclePlus2 size={16} />} onClick={handleCreateNewUploader}>Create</Button>
            </Group>
            <Box h="85vh">
                <ThemedAgGridReact
                    overlayNoRowsTemplate='No uploaders found'
                    loading={isPending && !error}
                    rowData={uploaders}
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