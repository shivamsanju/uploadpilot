import { IconDots, IconEdit, IconTrash } from '@tabler/icons-react';
import { ActionIcon, Avatar, Box, Group, Menu, Text } from '@mantine/core';
import { useGetUsersInWorkspace } from '../../../apis/workspace';
import { useParams } from 'react-router-dom';
import { useMemo } from 'react';
import { DataTableColumn } from 'mantine-datatable';
import { UploadPilotDataTable } from '../../../components/Table/Table';
import { ErrorCard } from '../../../components/ErrorCard/ErrorCard';
import { AppLoader } from '../../../components/Loader/AppLoader';

const getRandomAvatar = () => {
    const randomIndex = Math.floor(Math.random() * 5) + 1;
    return `https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-${randomIndex}.png`
}

const WorkspaceUsersList = () => {
    const { workspaceId } = useParams();
    const { isPending, error, users } = useGetUsersInWorkspace(workspaceId || '');

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
            accessor: 'rate',
            title: 'Rate',
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
                            <Menu.Item leftSection={<IconEdit size={16} stroke={1.5} />}>
                                Edit role
                            </Menu.Item>
                            <Menu.Item leftSection={<IconTrash size={16} stroke={1.5} />} color="red">
                                Remove from workspace
                            </Menu.Item>
                        </Menu.Dropdown>
                    </Menu>
                </Group>
            )
        },

    ], []);

    if (error) {
        return <ErrorCard title="Error" message={error.message} h="70vh" />
    }

    return (
        <Box mr="md">
            <UploadPilotDataTable
                fetching={isPending}
                showSearch={false}
                columns={columns}
                records={users}
                verticalSpacing="md"
                horizontalSpacing="md"
                noHeader={true}
                customLoader={<AppLoader h="50vh" />}
                noRecordsText="No users found"
            />
        </Box>
    );
}

export default WorkspaceUsersList