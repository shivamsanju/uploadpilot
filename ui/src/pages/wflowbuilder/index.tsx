import { Box, Group, LoadingOverlay, Title } from "@mantine/core";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useGetProcessor } from "../../apis/processors";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { Builder } from "./builder";
import { WorkflowBuilderProviderV2 } from "../../context/WflowEditorContextV2";

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
        {processor?.id && (
          <Builder workspaceId={workspaceId} processorId={processorId} />
        )}
      </Box>
    </WorkflowBuilderProviderV2>
  );
};

export default WorkflowBuilderPage;
