import { useEffect, useMemo, useState } from 'react';
import { TextInput, Stack, Button, Group, Title, Table, ActionIcon } from '@mantine/core';
import { IconBrandAws, IconBrandAzure, IconBrandGoogle, IconEye, IconFolder, IconSearch, IconSquarePlus, IconTrash } from '@tabler/icons-react';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../../config';
import { useNavigate } from 'react-router-dom';

const headers = ["Name", "Type", "Tags", "Created By", "Created At", "Actions"];
const ConnectorsPage = () => {
    const [search, setSearch] = useState('');
    const [connectors, setConnectors] = useState<any[]>([]);
    const navigate = useNavigate();

    const handleNewConnector = () => {
        navigate('/storage/connectors/new');
    }

    const getConnectors = async () => {
        try {
            const resp = await axios.get(getApiDomain() + "/storage/connectors");
            return resp.data;
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to fetch connectors',
                color: 'red',
            })
        }
    }

    const handleDeleteConnector = async (id: string) => {
        if (!window.confirm("Are you sure you want to delete this connector?")) {
            return
        }
        try {
            const resp = await axios.delete(getApiDomain() + "/connectors/" + id);
            if (resp.status === 200) {
                notifications.show({
                    title: 'Success',
                    message: 'Connector deleted successfully',
                    color: 'green',
                })
                getConnectors().then((data) => {
                    setConnectors(data)
                });
            }
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to delete connector',
                color: 'red',
            })
        }
    }

    const filteredConnectors: any = useMemo(() => {
        if (!connectors || connectors.length === 0) {
            return {
                head: headers,
                body: [],
            };
        }
        const fc = connectors.filter((connector) =>
            connector.name.toLowerCase().includes(search.toLowerCase())
        )
        return {
            head: headers,
            body: fc.map((connector) => ([
                connector.name,
                <>
                    {connector.type === "s3" && <Group align='center'><IconBrandAws size={16} /> S3</Group>}
                    {connector.type === "gcs" && <Group align='center'><IconBrandGoogle size={16} /> GCS</Group>}
                    {connector.type === "azure" && <Group align='center'><IconBrandAzure size={16} /> Azure</Group>}
                    {connector.type === "local" && <Group align='center'><IconFolder size={16} /> Local</Group>}
                </>,
                connector.tags?.join(', '),
                connector.createdBy,
                new Date(connector.createdAt).toLocaleString('en-CA'),
                <ActionIcon variant='filled' onClick={() => handleDeleteConnector(connector.id)}>
                    <IconTrash size={16} />
                </ActionIcon>
            ])),
        }
    }, [search, connectors]);

    useEffect(() => {
        getConnectors().then((data) => {
            setConnectors(data)
        });
    }, [])

    return (
        <>
            <Title order={3} mb="lg" opacity={0.8}>Connectors</Title>
            <Group justify='space-between' align='flex-start'>
                <TextInput
                    placeholder="Search connectors..."
                    leftSection={<IconSearch size={16} />}
                    value={search}
                    onChange={(event) => setSearch(event.currentTarget.value)}
                    mb="md"
                    w="90%"
                />
                <Button size="xs" leftSection={<IconSquarePlus size={16} />} onClick={handleNewConnector}>Create</Button>
            </Group>
            <Stack gap="md">
                <Table data={filteredConnectors} withTableBorder stickyHeader p="md" />
            </Stack>
        </>
    );
}

export default ConnectorsPage;
