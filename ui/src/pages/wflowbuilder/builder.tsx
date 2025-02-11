import { Box, Group, Paper, ScrollArea } from "@mantine/core";
import { BlockSearch } from "./blocksearch";
import YanlEditor from "./edits/yaml";

type Props = {
  workspaceId: string;
  processorId: string;
};
export const Builder: React.FC<Props> = ({ workspaceId, processorId }) => {
  return (
    <Paper withBorder>
      <Group justify="center" align="flex-start" gap={0}>
        <Box w="60%">
          {/* <WorkflowEditor
            workspaceId={workspaceId}
            processorId={processorId}
            statement={statement}
            variables={variables}
          /> */}
          <YanlEditor />
        </Box>
        <ScrollArea h="75vh" w="40%" scrollbarSize={6}>
          <Box m={0} px="md">
            {/* {selectedTask ? (
              <SelectedTaskForm />
            ) : (
              <BlockSearch processorId={processorId} />
            )} */}
            <BlockSearch processorId={processorId} />
          </Box>
        </ScrollArea>
      </Group>
    </Paper>
  );
};
