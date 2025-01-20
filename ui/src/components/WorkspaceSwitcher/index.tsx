import { Stack, Title, Button } from '@mantine/core';
import { IconSwitchHorizontal } from '@tabler/icons-react';
import { useNavigate } from "react-router-dom";
import { useGetWorkspaces } from '../../apis/workspace';


const WorkspaceSwitcher = () => {
    const navigate = useNavigate();
    const { isPending, error, workspaces } = useGetWorkspaces();


    const handleWorkspaceSwitch = () => {
        navigate(`/`);
    }

    if (isPending) {
        return <></>;
    }

    if (error) {
        return <>Error selecting workspace</>;
    }

    return (workspaces && workspaces.length > 0) ? (
        <Stack align='center' p="md">
            <Title order={4} opacity={0.7} lineClamp={1}>{workspaces[0].name}</Title>
            <Button
                variant='outline'
                leftSection={<IconSwitchHorizontal size={20} stroke={1.5} />}
                onClick={handleWorkspaceSwitch}
            >
                Switch Workspace
            </Button>
        </Stack>
    ) : <></>;
};

export default WorkspaceSwitcher;