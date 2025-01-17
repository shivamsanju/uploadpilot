import { useMemo } from 'react';
import { Button, Title, Menu, Group, Anchor } from '@mantine/core';
import { IconCirclePlus2, IconDots, IconEye } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useGetUploaders } from '../../apis/uploader';
import { timeAgo } from '../../utils/datetime';
import { DataTableColumn, } from 'mantine-datatable';
import { UploadPilotDataTable, useUploadPilotDataTable } from '../../components/Table/Table';
import ErrorCard from '../../components/ErrorCard/ErrorCard';

const recordsPerPage = 10;

const UploadersListPage = () => {
    const navigate = useNavigate();

    const { searchFilter, onSearchFilterChange, page, onPageChange } = useUploadPilotDataTable();

    const { isPending, error, uploaders, totalRecords, invalidate } = useGetUploaders({
        skip: (page - 1) * recordsPerPage,
        limit: recordsPerPage,
        search: searchFilter
    });


    const handleCreateNewUploader = () => {
        navigate('/uploaders/new');
    }


    const colDefs: DataTableColumn[] = useMemo(() => {
        return [
            {
                title: 'Name',
                accessor: 'name',
                render: ({ id, name }) => <Anchor size="sm" href={`/uploaders/${id}`}>{`${name}`}</Anchor>,
                flex: 1.7
            },
            {
                title: 'Created By', accessor: 'createdBy',
                filtering: true,
                editable: true,
                textAlign: 'center',
            },
            {
                title: 'Created At',
                accessor: 'createdAt',
                textAlign: 'center',
                render: ({ createdAt }) => createdAt ? timeAgo.format(new Date(createdAt as string)) : ""
            },
            {
                title: 'Updated At',
                textAlign: 'center',
                accessor: 'updatedAt',
                render: ({ updatedAt }) => updatedAt ? timeAgo.format(new Date(updatedAt as string)) : ""
            },
            {
                title: 'Actions',
                accessor: 'actions',
                flex: 0.4,
                textAlign: 'center',
                render: ({ id }) => (
                    <Menu>
                        <Menu.Target>
                            <IconDots size={15} style={{ cursor: 'pointer' }} />
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => navigate(`/uploaders/${id}`)}
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
            {error ?
                <ErrorCard title={error.name} message={error.message} h="70vh" /> :
                <UploadPilotDataTable
                    fetching={isPending}
                    showSearch
                    onSearchFilterChange={onSearchFilterChange}
                    showRefresh
                    onRefresh={invalidate}
                    height={"80vh"}
                    verticalSpacing={"sm"}
                    columns={colDefs}
                    records={uploaders}
                    page={page}
                    onPageChange={onPageChange}
                    recordsPerPage={10}
                    totalRecords={totalRecords}
                />
            }
        </>
    );
}

export default UploadersListPage;