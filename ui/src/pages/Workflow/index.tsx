import { useMemo, } from 'react';
import { Button, Title, Menu, Group, Box, Anchor } from '@mantine/core';
import { IconCirclePlus2, IconDots, IconEye } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { useGetWorkflows } from '../../apis/workflow';
import { ThemedAgGridReact } from '../../components/AgGrid/AgGrid';


TimeAgo.addDefaultLocale(en)
const timeAgo = new TimeAgo('en-US');

const WorkflowsPage = () => {
    const navigate = useNavigate();

    const { isPending, error, workflows } = useGetWorkflows();

    const handleNewWorkflow = () => {
        navigate('/workflows/new');
    }

    const colDefs = useMemo(() => {
        return [
            {
                headerName: 'Name',
                field: 'name',
                cellRenderer: (params: any) => <Anchor size='sm' href={`/workflows/${params.data.id}`}>{params.value}</Anchor>,
                flex: 1.7
            },
            { headerName: 'Created By', field: 'createdBy' },
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
                                onClick={() => navigate(`/workflows/${params.data.id}`)}
                                leftSection={<IconEye size={15} />}
                            >
                                View
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
                <Title order={3} mb="lg" opacity={0.8}>Workflows</Title>
                <Button size="xs" leftSection={<IconCirclePlus2 size={16} />} onClick={handleNewWorkflow}>Create</Button>
            </Group>
            <Box h="83vh">
                <ThemedAgGridReact
                    overlayNoRowsTemplate='No Workflows Found'
                    loading={isPending && !error}
                    rowData={workflows}
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

export default WorkflowsPage;