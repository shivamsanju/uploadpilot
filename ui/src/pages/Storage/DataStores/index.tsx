import { useMemo, } from 'react';
import { Button, Group, Title, Menu, Box, Anchor } from '@mantine/core';
import { IconCirclePlus2, IconDots, IconTrash } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { useNavigate } from 'react-router-dom';
import { useGetDataStores } from '../../../apis/storage';
import { ThemedAgGridReact } from '../../../components/AgGrid/AgGrid';
import { timeAgo } from '../../../utils/datetime';

const DataStoresPage = () => {
    const navigate = useNavigate();

    const { isPending, error, connectors } = useGetDataStores();

    const handleNewConnector = () => {
        navigate('/storage/connectors/new');
    }



    const handleDeleteDataStore = async (id: string) => {
        notifications.show({
            title: 'Coming Soon',
            message: 'This feature is coming soon',
            color: 'yellow',
        })

    }


    const colDefs = useMemo(() => {
        return [
            {
                headerName: 'Name',
                field: 'name',
                flex: 1.5
            },
            {
                headerName: 'Connector',
                field: 'connectorName',
                cellRenderer: (params: any) => <Anchor size='xs' href={`/storage/connectors/${params.data.connectorId}`}>{params.value}</Anchor>,
            },
            {
                headerName: 'Bucket',
                field: 'bucket',
            },
            { headerName: 'Created By', field: 'createdBy' },
            { headerName: 'Created At', field: 'createdAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            { headerName: 'Updated At', field: 'updatedAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            {
                headerName: 'Actions',
                field: 'actions',
                flex: 0.5,
                cellRenderer: (params: any) => (
                    <Menu>
                        <Menu.Target>
                            <IconDots size={15} style={{ cursor: 'pointer' }} />
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => handleDeleteDataStore(params.data.id)}
                                leftSection={<IconTrash size={13} />}
                            >
                                Delete
                            </Menu.Item>

                        </Menu.Dropdown>
                    </Menu>
                )
            },
        ]
    }, []);


    return (
        <>
            <Group justify='space-between' align='flex-start' gap="lg">
                <Title order={3} mb="lg" opacity={0.8}>Storage Connectors</Title>
                <Button size="xs" leftSection={<IconCirclePlus2 size={16} />} onClick={handleNewConnector}>Create</Button>
            </Group>
            <Box h="85vh">
                <ThemedAgGridReact
                    overlayNoRowsTemplate='No storage connectors found'
                    loading={isPending && !error}
                    rowData={connectors}
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

export default DataStoresPage;
