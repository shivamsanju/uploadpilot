import { IconDots, IconEdit, IconEye, IconTrash } from '@tabler/icons-react';
import { ActionIcon, Avatar, Box, Group, LoadingOverlay, Menu, Modal, Text } from '@mantine/core';
import { useGetUsersInWorkspace, useRemoveUserFromWorkspaceMutation } from '../../apis/workspace';
import { useParams } from 'react-router-dom';
import { useCallback, useMemo, useState } from 'react';
import { DataTableColumn } from 'mantine-datatable';
import { UploadPilotDataTable } from '../../components/Table/Table';
import { ErrorCard } from '../../components/ErrorCard/ErrorCard';
import { showNotification } from '@mantine/notifications';
import AddUserForm from './Add';
import { useViewportSize } from '@mantine/hooks';

const getRandomAvatar = () => {
    const randomIndex = Math.floor(Math.random() * 5) + 1;
    return `https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-${randomIndex}.png`
}

const WorkspaceUsersList = ({ opened, setOpened }: { opened: boolean, setOpened: any }) => {
    const [mode, setMode] = useState<'add' | 'edit' | 'view'>('add');
    const { width } = useViewportSize();

    const [initialValues, setInitialValues] = useState(null);
    const { workspaceId } = useParams();
    const { isPending, error, users } = useGetUsersInWorkspace(workspaceId || '');
    const { mutateAsync, isPending: removePending } = useRemoveUserFromWorkspaceMutation();


    const handleRemoveUser = useCallback(async (userId: string) => {
        if (!workspaceId || workspaceId === '' || !userId || userId === '') {
            showNotification({
                color: 'red',
                title: 'Error',
                message: 'Workspace ID or User ID is not available'
            })
            return
        };

        try {
            await mutateAsync({ workspaceId, userId });
        } catch (error) {
            console.error(error);
        }

    }, [workspaceId, mutateAsync]);


    const handleViewEdit = useCallback((item: any, mode: "edit" | "view") => {
        setMode(mode);
        setInitialValues(item);
        setOpened(true);
    }, [setMode, setInitialValues, setOpened]);


    const columns: DataTableColumn[] = useMemo(() => [
        {
            accessor: 'name',
            title: 'Name',
            render: (item: any) => (
                <Group gap="sm">
                    <Avatar size={40} src={getRandomAvatar()} radius={40} />
                    <div>
                        <Text fz="sm" fw={500}>
                            {item.name}
                        </Text>
                        <Text c="dimmed" fz="xs">
                            Name
                        </Text>
                    </div>
                </Group>
            ),
        },
        {
            accessor: 'email',
            title: 'Email',
            hidden: width < 768,
            render: (item: any) => (
                <>
                    <Text fz="sm">{item.email}</Text>
                    <Text fz="xs" c="dimmed">
                        Email
                    </Text>
                </>
            ),
        },
        {
            accessor: 'role',
            title: 'Role',
            render: (item: any) => (
                <>
                    <Text fz="sm">{item.role}</Text>
                    <Text fz="xs" c="dimmed">
                        Role
                    </Text>
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
                                leftSection={<IconEye size={16} stroke={1.5} />}
                                onClick={() => handleViewEdit(item, "view")}>
                                <Text>View</Text>
                            </Menu.Item>
                            <Menu.Item
                                leftSection={<IconEdit size={16} stroke={1.5} />}
                                onClick={() => handleViewEdit(item, "edit")}>
                                <Text>Edit role</Text>
                            </Menu.Item>
                            <Menu.Item
                                leftSection={<IconTrash size={16} stroke={1.5} />}
                                color="red"
                                onClick={() => handleRemoveUser(item.userId)}
                            >
                                <Text>Remove from workspace</Text>
                            </Menu.Item>

                        </Menu.Dropdown>
                    </Menu>
                </Group>
            )
        },

    ], [handleRemoveUser, handleViewEdit, width]);

    if (error) {
        return <ErrorCard title="Error" message={error.message} h="70vh" />
    }

    return (
        <Box mr="md">
            <LoadingOverlay visible={removePending} overlayProps={{ radius: "sm", blur: 1 }} />
            <UploadPilotDataTable
                minHeight={500}
                fetching={isPending}
                showSearch={false}
                columns={columns}
                records={users}
                verticalSpacing="md"
                horizontalSpacing="md"
                noHeader={true}
                noRecordsText="No users found"
            />
            <Modal
                padding="xl"
                transitionProps={{ transition: 'pop' }}
                opened={opened}
                onClose={() => {
                    setOpened(false);
                    setMode('add');
                    setInitialValues(null);
                }}
                title={mode === 'edit' ? 'Edit User' : mode === 'view' ? 'User Details' : 'Add User'}
                closeOnClickOutside={false}
                size="lg"
            >
                <AddUserForm mode={mode} setOpened={setOpened} workspaceId={workspaceId || ""} initialValues={initialValues} setInitialValues={setInitialValues} setMode={setMode} />
            </Modal>
        </Box>
    );
}

export default WorkspaceUsersList