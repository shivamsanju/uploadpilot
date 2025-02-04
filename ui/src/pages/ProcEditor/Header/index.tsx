import { Group, Title, Text, ActionIcon, Box, Tooltip } from "@mantine/core";
import { IconChevronLeft } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import UserButton from "../../../components/UserMenu";
import { useCanvas } from "../../../context/ProcEditorContext";

export const ProcEditorHeader = () => {
  const navigate = useNavigate();
  const { workspaceId, processor } = useCanvas();

  const handleBack = () => navigate(`/workspaces/${workspaceId}/processors`);

  return (
    <Group justify="space-between" align="center" px="md" h="100%">
      <Group gap="xl" align="center">
        <ActionIcon
          radius="50%"
          variant="default"
          size="lg"
          onClick={handleBack}
        >
          <IconChevronLeft color="gray" />
        </ActionIcon>
        <Box maw="70vw">
          <Title order={4} opacity={0.7}>
            {processor?.name}
          </Title>
          <Tooltip
            label={`Triggers: ${processor?.triggers?.join(", ")}`}
            position="bottom"
            w="300px"
            multiline
            withArrow
          >
            <Text c="dimmed" lineClamp={1}>
              Triggers: {processor?.triggers?.join(", ")}
            </Text>
          </Tooltip>
        </Box>
      </Group>
      <UserButton />
    </Group>
  );
};
