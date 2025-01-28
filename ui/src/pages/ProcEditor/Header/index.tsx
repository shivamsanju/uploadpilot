import { Group, Title, Text, Button, ActionIcon, Box } from "@mantine/core";
import { IconChevronLeft } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import UserButton from "../../../components/UserMenu";
import { useCanvas } from "../../../hooks/DndCanvas";

export const ProcEditorHeader = () => {
    const navigate = useNavigate();
    const { isUpdating, isPending, handleSave, handleDiscard, workspaceId } = useCanvas();

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
            <Group>
                <Button variant="default" size="sm" c="dimmed" loading={isUpdating || isPending} onClick={handleDiscard}>Discard</Button>
                <Button size="sm" loading={isUpdating || isPending} onClick={handleSave}>Save</Button>
            </Group>
            <UserButton />
        </Group>
    )
};