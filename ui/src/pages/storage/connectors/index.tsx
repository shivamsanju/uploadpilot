import { useMemo, } from 'react';
import { Button, Group, Title, Menu, Box } from '@mantine/core';
import { IconBrandAws, IconBrandAzure, IconBrandGoogle, IconCirclePlus2, IconDots, IconFolder, IconTrash } from '@tabler/icons-react';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../../utils/config';
import { useNavigate } from 'react-router-dom';
import { useGetStorageConnectors } from '../../../apis/storage';
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { ThemedAgGridReact } from '../../../components/AgGrid/AgGrid';


TimeAgo.addDefaultLocale(en)
const timeAgo = new TimeAgo('en-US');

const ConnectorsPage = () => {
    const navigate = useNavigate();

    const { isPending, error, connectors } = useGetStorageConnectors();

    const handleNewConnector = () => {
        navigate('/storage/connectors/new');
    }

    const handleDeleteConnector = async (id: string) => {
        notifications.show({
            title: 'Comming Soon',
            message: 'This feature is comming soon',
            color: 'yellow',
        })

    }


    const colDefs = useMemo(() => {
        return [

            {
                headerName: '',
                field: 'type',
                flex: 0.5,
                cellRenderer: (params: any) => {
                    switch (params.value) {
                        case 's3':
                            return <IconBrandAws size={15} />;
                        case 'gcs':
                            return <IconBrandGoogle size={15} />;
                        case 'azure':
                            return <IconBrandAzure size={15} />;
                        default:
                            return <IconFolder size={15} />;
                    }
                },
                sortable: false,
                filter: false,
                cellStyle: { textAlign: 'center' }
            },
            {
                headerName: 'Name',
                field: 'name',
                flex: 1.5
            },
            {
                headerName: 'Type',
                field: 'type',
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
                                onClick={() => handleDeleteConnector(params.data.id)}
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
            <Box h="83vh">
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

export default ConnectorsPage;
