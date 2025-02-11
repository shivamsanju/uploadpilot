import {
  Box,
  Group,
  LoadingOverlay,
  Paper,
  ScrollArea,
  Title,
} from "@mantine/core";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useGetProcessor } from "../../apis/processors";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { WorkflowBuilderProviderV2 } from "../../context/WflowEditorContextV2";
import { WorkflowYamlEditor } from "./editor";
import { BlockSearch } from "./blocksearch";

const WorkflowBuilderPage = () => {
  const { workspaceId, processorId } = useParams();

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
    <WorkflowBuilderProviderV2>
      <Box mb={50}>
        <LoadingOverlay
          visible={isPending}
          overlayProps={{ backgroundOpacity: 0 }}
          zIndex={1000}
        />
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
                />
              )}
            </Box>
            <ScrollArea h="75vh" w="40%" scrollbarSize={6}>
              <Box m={0} px="md">
                <BlockSearch processorId={processorId} />
              </Box>
            </ScrollArea>
          </Group>
        </Paper>
      </Box>
    </WorkflowBuilderProviderV2>
  );
};

export default WorkflowBuilderPage;
