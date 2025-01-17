import { useMemo, } from 'react';
import { Button, Group, Title, Menu } from '@mantine/core';
import { IconBrandAws, IconBrandAzure, IconBrandGoogle, IconCirclePlus2, IconDots, IconEdit, IconFolder, IconTrash } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { useNavigate } from 'react-router-dom';
import { useGetStorageConnectors } from '../../apis/storage';
import { timeAgo } from '../../utils/datetime';
import { UploadPilotDataTable, useUploadPilotDataTable } from '../../components/Table/Table';
import ErrorCard from '../../components/ErrorCard/ErrorCard';
import { DataTableColumn } from 'mantine-datatable';

const recordsPerPage = 10;

const getConnectorIcon = (type: string) => {
    switch (type) {
        case 's3':
            return <IconBrandAws size={15} />;
        case 'gcs':
            return <IconBrandGoogle size={15} />;
        case 'azure':
            return <IconBrandAzure size={15} />;
        default:
            return <IconFolder size={15} />;
    }
}
const ConnectorsPage = () => {
    const navigate = useNavigate();

    const { searchFilter, onSearchFilterChange, page, onPageChange } = useUploadPilotDataTable();

    const { isPending, error, connectors, totalRecords, invalidate } = useGetStorageConnectors({
        skip: (page - 1) * recordsPerPage,
        limit: recordsPerPage,
        search: searchFilter
    });

    const handleNewConnector = () => {
        navigate('/storageConnectors/new');
    }

    const handleDeleteConnector = async (id: string) => {
        notifications.show({
            title: 'Coming Soon',
            message: 'This feature is coming soon',
            color: 'yellow',
        })

    }


    const colDefs: DataTableColumn[] = useMemo(() => {
        return [
            {
                title: 'Name',
                accessor: 'name',
                flex: 1.5
            },
            {
                title: 'Type',
                accessor: 'type',
                textAlign: 'center',
                render: (params: any) => <Group align='center' justify='center' gap='sm'>{getConnectorIcon(params.type)} {params.type}</Group>,
            },
            {
                title: 'Created By',
                accessor: 'createdBy',
                textAlign: 'center',
            },
            {
                title: 'Created At',
                accessor: 'createdAt',
                textAlign: 'center',
                render: (params: any) => params.createdAt && timeAgo.format(new Date(params.createdAt))
            },
            {
                title: 'Updated At',
                accessor: 'updatedAt',
                textAlign: 'center',
                render: (params: any) => params.updatedAt && timeAgo.format(new Date(params.updatedAt))
            },
            {
                title: 'Actions',
                accessor: 'actions',
                flex: 0.5,
                textAlign: 'center',
                render: (params: any) => (
                    <Menu>
                        <Menu.Target>
                            <IconDots size={15} style={{ cursor: 'pointer' }} />
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => handleDeleteConnector(params.id)}
                                leftSection={<IconEdit size={13} />}
                            >
                                Edit
                            </Menu.Item>
                            <Menu.Item
                                onClick={() => handleDeleteConnector(params.id)}
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
            {error ?
                <ErrorCard title={error.name} message={error.message} h="80vh" /> :
                <UploadPilotDataTable
                    fetching={isPending}
                    showSearch
                    onSearchFilterChange={onSearchFilterChange}
                    showRefresh
                    onRefresh={invalidate}
                    height={"80vh"}
                    verticalSpacing={"sm"}
                    columns={colDefs}
                    records={connectors}
                    page={page}
                    onPageChange={onPageChange}
                    recordsPerPage={10}
                    totalRecords={totalRecords}
                />
            }
        </>
    );
}

export default ConnectorsPage;
