import { IconCircleCheck, IconCircleOff, IconDots, IconEdit, IconTrash, IconWebhook, IconWebhookOff } from '@tabler/icons-react';
import { ActionIcon, Avatar, Badge, Box, Group, Menu, Modal, Paper, Stack, Text, Title } from '@mantine/core';
import { useParams } from 'react-router-dom';
import { useCallback, useMemo, useState } from 'react';
import { DataTableColumn } from 'mantine-datatable';
import { UploadPilotDataTable } from '../../../components/Table/Table';
import { ErrorCard } from '../../../components/ErrorCard/ErrorCard';
import { showNotification } from '@mantine/notifications';
import { useDeleteWebhookMutation, useEnableDisableWebhookMutation, useGetWebhooks } from '../../../apis/webhooks';
import AddWebhookForm from './Add';
import { timeAgo } from '../../../utils/datetime';


const WebhooksList = ({ opened, setOpened }: { opened: boolean, setOpened: any }) => {
    const [mode, setMode] = useState<'add' | 'edit'>('add');
    const [initialValues, setInitialValues] = useState(null);
    const { workspaceId } = useParams();
    const { isPending, error, webhooks } = useGetWebhooks(workspaceId || '');
    const { mutateAsync } = useDeleteWebhookMutation();
    const { mutateAsync: enableDisableWebhook } = useEnableDisableWebhookMutation();


    const handleRemoveWebhook = useCallback(async (webhookId: string) => {
        if (!workspaceId || workspaceId === '' || !webhookId || webhookId === '') {
            showNotification({
                color: 'red',
                title: 'Error',
                message: 'Workspace ID or Webhook ID is not available'
            })
            return
        };

        try {
            await mutateAsync({ workspaceId, webhookId });
        } catch (error) {
            console.error(error);
        }

    }, [workspaceId, mutateAsync]);

    const handleEnableDisableWebhook = useCallback(async (webhookId: string, enabled: boolean) => {
        if (!workspaceId || workspaceId === '' || !webhookId || webhookId === '') {
            showNotification({
                color: 'red',
                title: 'Error',
                message: 'Workspace ID or Webhook ID is not available'
            })
            return
        };

        try {
            await enableDisableWebhook({ workspaceId, webhookId, enabled });
        } catch (error) {
            console.error(error);
        }

    }, [workspaceId, enableDisableWebhook]);

    const handleEdit = useCallback(async (values: any) => {
        setInitialValues(values);
        setMode('edit');
        setOpened(true);
    }, [setOpened, setMode, setInitialValues]);


    const columns: DataTableColumn[] = useMemo(() => [
        {
            accessor: 'event',
            title: 'Event',
            render: (item: any) => (
                <Group gap="sm">
                    <Avatar size={40} radius={40} variant='light'>
                        {item?.enabled ? <IconWebhook color="green" /> : <IconWebhookOff />}
                    </Avatar>
                    <div>
                        <Text fz="sm" fw={500}>
                            {item.event}
                        </Text>
                        <Text c="dimmed" fz="xs">
                            Event
                        </Text>
                    </div>
                </Group>
            ),
        },
        {
            accessor: 'url',
            title: 'URL',
            render: (item: any) => (
                <>
                    <Text fz="sm">{item.url}</Text>
                    <Text fz="xs" c="dimmed">
                        URL
                    </Text>
                </>
            ),
        },
        {
            title: 'updated At',
            accessor: 'updatedAt',
            render: (params: any) => (
                <>
                    <Text fz="sm">{params?.updatedAt && timeAgo.format(new Date(params?.updatedAt))}</Text>
                    <Text fz="xs" c="dimmed">
                        Last Updated
                    </Text>
                </>
            )
        },
        {
            accessor: 'enabled',
            title: 'Status',
            textAlign: 'center',
            render: (item: any) => (
                <>
                    <Badge color={item?.enabled ? 'green' : 'red'} size="sm">{item?.enabled ? 'Enabled' : 'Disabled'}</Badge>
                </>
            ),
        },
        {
            accessor: 'actions',
            title: 'Actions',
            render: (item: any) => (
                <Group gap={0} justify="flex-end">
                    <Menu
                        transitionProps={{ transition: 'pop' }}
                        withArrow
                        position="bottom-end"
                        withinPortal
                    >
                        <Menu.Target>
                            <ActionIcon variant="subtle" color="dimmed">
                                <IconDots size={16} stroke={1.5} />
                            </ActionIcon>
                        </Menu.Target>
                        <Menu.Dropdown>
                            <Menu.Item
                                leftSection={<IconEdit size={16} stroke={1.5} />}
                                disabled={!item?.enabled}
                                onClick={() => handleEdit(item)}
                            >
                                Edit
                            </Menu.Item>
                            <Menu.Item
                                leftSection={item?.enabled ? <IconCircleOff size={16} stroke={1.5} /> : <IconCircleCheck size={16} stroke={1.5} />}
                                color={item?.enabled ? 'red' : 'green'}
                                onClick={() => handleEnableDisableWebhook(item.id, !item?.enabled)}
                            >
                                {item?.enabled ? 'Disable' : 'Enable'}
                            </Menu.Item>
                            <Menu.Item
                                leftSection={<IconTrash size={16} stroke={1.5} />}
                                color="red"
                                onClick={() => handleRemoveWebhook(item.id)}
                            >
                                Delete
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>
                </Group>
            )
        },

    ], [handleRemoveWebhook, handleEnableDisableWebhook, handleEdit]);

    if (error) {
        return <ErrorCard title="Error" message={error.message} h="70vh" />
    }

    return (
        <Box mr="md">
            {!isPending && (!webhooks || webhooks.length === 0) ? (
                <Stack justify="center" align="center">
                    <Title order={3} opacity={0.7}>Create your first webhook</Title>
                    <Paper p={{ base: "md", md: "xl" }} mt="md" miw="300" maw="1000" w="50vw">
                        <AddWebhookForm mode={mode} setOpened={setOpened} workspaceId={workspaceId || ""} initialValues={initialValues} setInitialValues={setInitialValues} setMode={setMode} />
                    </Paper>
                </Stack>
            ) : (
                <UploadPilotDataTable
                    minHeight={700}
                    fetching={isPending}
                    showSearch={false}
                    columns={columns}
                    records={webhooks}
                    verticalSpacing="md"
                    horizontalSpacing="md"
                    noHeader={true}
                    noRecordsText="No webhooks found"
                />
            )}

            <Modal
                padding="xl"
                transitionProps={{ transition: 'pop' }}
                opened={opened}
                onClose={() => {
                    setOpened(false);
                    setInitialValues(null);
                    setMode('add');
                }}
                title={<Title order={3} opacity={0.7}>Add Webhook</Title>}
                closeOnClickOutside={false}
                size="xl"
            >
                <AddWebhookForm mode={mode} setOpened={setOpened} workspaceId={workspaceId || ""} initialValues={initialValues} setInitialValues={setInitialValues} setMode={setMode} />
            </Modal>
        </Box>
    );
}

export default WebhooksList