import { Box, Group, LoadingOverlay, Title } from "@mantine/core";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useGetProcessor } from "../../apis/processors";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { WorkflowBuilderProvider } from "../../context/WflowEditorContext";
import { Builder } from "./builder";

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
    <WorkflowBuilderProvider>
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
        <Builder workspaceId={workspaceId} processorId={processorId} />
      </Box>
    </WorkflowBuilderProvider>
  );
};

export default WorkflowBuilderPage;
