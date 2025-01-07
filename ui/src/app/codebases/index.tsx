import { useEffect, useMemo, useState } from 'react';
import { Container, TextInput, Stack, Button, Modal, Group } from '@mantine/core';
import { IconSearch, IconPlus } from '@tabler/icons-react';
import CodeBaseCard from '../../components/CodebaseCard';
import { Codebase } from '../../types/codebase';
import CreateCodebaseForm from '../../components/CodebaseForm';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../config';

const CodeBases = () => {
    const [search, setSearch] = useState('');
    const [opened, setOpened] = useState(false);
    const [codebases, setCodebases] = useState<Codebase[]>([]);


    const handleSubmit = async (c: Codebase) => {
        try {
            const resp = await axios.post(getApiDomain() + "/codebase", c);
            if (resp.status === 200) {
                notifications.show({
                    title: 'Success',
                    message: 'Codebase added successfully',
                    color: 'green',
                })
                setCodebases((prev) => [resp.data, ...prev]);
                setOpened(false);
            }
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to add codebase',
                color: 'red',
            })
        }
    }

    const getCodebases = async () => {
        try {
            const resp = await axios.get(getApiDomain() + "/codebase");
            return resp.data;
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to fetch codebases',
                color: 'red',
            })
        }
    }

    const filteredCodebases = useMemo(() => {
        return codebases.filter((codebase) =>
            codebase.name.toLowerCase().includes(search.toLowerCase()) ||
            codebase.description.toLowerCase().includes(search.toLowerCase()) ||
            codebase.tags.some((tag) => tag.toLowerCase().includes(search.toLowerCase()))
        )
    }, [search, codebases]);

    useEffect(() => {
        getCodebases().then((data) => {
            setCodebases(data)
        });
    }, [])

    return (
        <Container size="lg">
            <Group justify='space-between' align='flex-start'>
                <TextInput
                    placeholder="Search codebases..."
                    leftSection={<IconSearch size={16} />}
                    value={search}
                    onChange={(event) => setSearch(event.currentTarget.value)}
                    mb="md"
                    w="40%"
                />
                <Button size="xs" leftSection={<IconPlus size={16} />} onClick={() => setOpened(true)}>Add New Codebase</Button>
            </Group>
            <Stack gap="md">
                {filteredCodebases.map((codebase) => (
                    <CodeBaseCard codebase={codebase} />
                ))}
                <Modal
                    opened={opened}
                    onClose={() => setOpened(false)}
                    title="Add New Codebase"
                >
                    <CreateCodebaseForm onSubmit={handleSubmit} />
                </Modal>
            </Stack>
        </Container>
    );
}

export default CodeBases;