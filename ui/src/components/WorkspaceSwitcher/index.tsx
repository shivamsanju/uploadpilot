import { Stack, Title, Button } from '@mantine/core';
import { IconSwitch3 } from '@tabler/icons-react';
import { useNavigate, useParams } from "react-router-dom";
import { useGetWorkspaces } from '../../apis/workspace';


const WorkspaceSwitcher = () => {
    const navigate = useNavigate();
    const { workspaceId } = useParams();
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
            <Title order={4} opacity={0.7} lineClamp={1}>{workspaces.find((w: any) => w?.id === workspaceId)?.name}</Title>
            <Button
                variant='subtle'
                leftSection={<IconSwitch3 stroke={1.5} />}
                onClick={handleWorkspaceSwitch}
            >
                Switch Workspace
            </Button>
        </Stack>
    ) : <></>;
};

export default WorkspaceSwitcher;