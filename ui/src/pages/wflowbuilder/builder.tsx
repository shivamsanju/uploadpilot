import { Box, Group, Paper, ScrollArea } from "@mantine/core";
import { WorkflowEditor } from "./editor";
import { SelectedTaskForm } from "./forms";
import { useWorkflowBuilder } from "../../context/WflowEditorContext";
import { BlockSearch } from "./blocksearch";

type Props = {
  workspaceId: string;
  processorId: string;
};
export const Builder: React.FC<Props> = ({ workspaceId, processorId }) => {
  const { selectedTask } = useWorkflowBuilder();

  return (
    <Paper withBorder>
      <Group justify="center" align="flex-start" gap={0}>
        <Box w="40%">
          <WorkflowEditor workspaceId={workspaceId} processorId={processorId} />
        </Box>
        <ScrollArea h="75vh" w="60%" scrollbarSize={6}>
          <Box m={0} px="md">
            {selectedTask ? (
              <SelectedTaskForm />
            ) : (
              <BlockSearch processorId={processorId} />
            )}
          </Box>
        </ScrollArea>
      </Group>
    </Paper>
  );
};
