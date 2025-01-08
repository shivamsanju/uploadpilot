import { useEffect, useMemo, useState } from 'react';
import { TextInput, Stack, Button, Group, Title, Table, ActionIcon } from '@mantine/core';
import { IconEye, IconSearch, IconSquarePlus } from '@tabler/icons-react';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../config';
import { useNavigate } from 'react-router-dom';

const headers = ["Name", "Tags", "Created By", "Created At", "Updated At", "Actions"];
const WorkflowsPage = () => {
    const [search, setSearch] = useState('');
    const [workflows, setWorkflows] = useState<any[]>([]);
    const navigate = useNavigate();

    const handleNewWorkflow = () => {
        navigate('/workflows/new');
    }


    const getWorkflows = async () => {
        try {
            const resp = await axios.get(getApiDomain() + "/workflows");
            return resp.data;
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to fetch workflows',
                color: 'red',
            })
        }
    }

    const filteredworkflows: any = useMemo(() => {
        if (!workflows || workflows.length === 0) {
            return {
                head: headers,
                body: [],
            };
        }
        const fw = workflows.filter((workflow) =>
            workflow.name.toLowerCase().includes(search.toLowerCase()) ||
            workflow.description.toLowerCase().includes(search.toLowerCase()) ||
            workflow.tags.some((tag: any) => tag.toLowerCase().includes(search.toLowerCase()))
        )
        return {
            head: headers,
            body: fw.map((workflow) => ([
                workflow.name,
                workflow.tags.join(', '),
                workflow.createdBy,
                new Date(workflow.createdAt).toLocaleString('en-CA'),
                new Date(workflow.updatedAt).toLocaleString('en-CA'),
                <ActionIcon variant='filled' onClick={() => navigate(`/workflows/${workflow.id}`)}>
                    <IconEye size={16} />
                </ActionIcon>
            ])),
        }
    }, [search, workflows, navigate]);

    useEffect(() => {
        getWorkflows().then((data) => {
            setWorkflows(data)
        });
    }, [])

    return (
        <>
            <Title order={3} mb="lg" opacity={0.8}>Workflows</Title>
            <Group justify='space-between' align='flex-start'>
                <TextInput
                    placeholder="Search workflows..."
                    leftSection={<IconSearch size={16} />}
                    value={search}
                    onChange={(event) => setSearch(event.currentTarget.value)}
                    mb="md"
                    w="90%"
                />
                <Button size="xs" leftSection={<IconSquarePlus size={16} />} onClick={handleNewWorkflow}>Create</Button>
            </Group>
            <Stack gap="md">
                <Table data={filteredworkflows} withTableBorder stickyHeader p="md" />
            </Stack>
        </>
    );
}

export default WorkflowsPage;