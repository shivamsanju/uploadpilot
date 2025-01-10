import { useMemo, } from 'react';
import { Button, Title, Menu, Group, Box, Anchor } from '@mantine/core';
import { IconCirclePlus2, IconDots, IconEye, } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { useGetImportPolicies } from '../../apis/importPolicy';
import { ThemedAgGridReact } from '../../components/AgGrid/AgGrid';

TimeAgo.addDefaultLocale(en)
const timeAgo = new TimeAgo('en-US');

const ImportPoliciesPage = () => {
    const navigate = useNavigate();
    const { isPending, error, importPolicies } = useGetImportPolicies();

    const colDefs = useMemo(() => {
        return [
            {
                headerName: 'Name',
                field: 'name',
                cellRenderer: (params: any) => <Anchor size='sm' href={`/importPolicies/${params.data.id}`}>{params.value}</Anchor>,
                flex: 2
            },
            { headerName: 'Max File Size (KB)', field: 'maxFileSizeKb' },
            { headerName: 'Max File Count', field: 'maxFileCount' },
            { headerName: 'Created By', field: 'createdBy' },
            { headerName: 'Created', field: 'createdAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            { headerName: 'Last Updated', field: 'updatedAt', valueFormatter: (params: any) => params.value && timeAgo.format(new Date(params.value)) },
            {
                headerName: 'Actions',
                field: 'actions',
                flex: 0.6,
                cellRenderer: (params: any) => (
                    <Menu>
                        <Menu.Target>
                            <IconDots size={15} style={{ cursor: 'pointer' }} />
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                onClick={() => navigate(`/importPolicies/${params.data.id}`)}
                                leftSection={<IconEye size={15} />}
                            >
                                View Details
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>
                )
            },
        ]
    }, []);




    const handleNewImportPolicy = () => {
        navigate('/importPolicies/new');
    }

    return (
        <>
            <Group justify='space-between' align='flex-start' gap="lg">
                <Title order={3} mb="lg" opacity={0.8}>Import Policies</Title>
                <Button size="xs" leftSection={<IconCirclePlus2 size={16} />} onClick={handleNewImportPolicy}>Create</Button>
            </Group>
            <Box h="83vh">
                <ThemedAgGridReact
                    overlayNoRowsTemplate='No Import Policies Found'
                    loading={isPending && !error}
                    rowData={importPolicies}
                    columnDefs={colDefs}
                    rowHeight={40}
                    headerHeight={40}
                    defaultColDef={{ sortable: true, filter: true, flex: 1 }}
                    pagination
                    paginationAutoPageSize
                    paginationPageSizeSelector={false}

                />
            </Box>
        </>
    );
}

export default ImportPoliciesPage;