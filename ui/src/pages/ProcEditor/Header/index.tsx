import { Group, Title, Text, ActionIcon, Box } from "@mantine/core";
import { IconChevronLeft } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import UserButton from "../../../components/UserMenu";
import { useCanvas } from "../../../context/EditorCtx";

export const ProcEditorHeader = () => {
    const navigate = useNavigate();
    const { workspaceId } = useCanvas();

    const handleBack = () => navigate(`/workspaces/${workspaceId}/processors`);

    return (
        <Group justify="space-between" align="center" px="md" h="100%">
            <Group gap="xl" align="center">
                <ActionIcon radius="50%" variant="default" size="lg" onClick={handleBack} >
                    <IconChevronLeft color="gray" />
                </ActionIcon>
                <Box>
                    <Title order={4} opacity={0.7}>Processor XYZ</Title>
                    <Text c="dimmed">Manage users and roles in this workspace</Text>
                </Box>
            </Group>
            <UserButton />
        </Group>
    )
};