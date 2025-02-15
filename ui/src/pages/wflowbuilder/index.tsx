import { Box, Group, Paper, ScrollArea, Title } from "@mantine/core";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useGetProcessor } from "../../apis/processors";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { WorkflowYamlEditor } from "./editor";
import { BlockSearch } from "./blocksearch";
import { ContainerOverlay } from "../../components/Overlay";
import { useState } from "react";

const WorkflowBuilderPage = () => {
  const { workspaceId, processorId } = useParams();
  const [editor, setEditor] = useState<any>(null);

  const { isPending, error, processor } = useGetProcessor(
    workspaceId as string,
    processorId as string
  );

  if (!workspaceId || !processorId) {
    return <AppLoader h="70vh" />;
  }

  if (error) {
    return <ErrorCard title={error.name} message={error.message} h="70vh" />;
  }

  return (
    <Box mb={50}>
      <ContainerOverlay visible={isPending} />
      <Group justify="space-between" mb="xl">
        <Title order={3} opacity={0.7}>
          Workflow builder for processor {processor?.name}
        </Title>
      </Group>
      <Paper withBorder>
        <Group justify="center" align="flex-start" gap={0}>
          <Box w="60%">
            {processor && (
              <WorkflowYamlEditor
                processor={processor}
                workspaceId={workspaceId}
                setEditor={setEditor}
                editor={editor}
              />
            )}
          </Box>
          <ScrollArea h="75vh" w="40%" scrollbarSize={6}>
            <Box m={0} px="md">
              <BlockSearch processorId={processorId} editor={editor} />
            </Box>
          </ScrollArea>
        </Group>
      </Paper>
    </Box>
  );
};

export default WorkflowBuilderPage;
