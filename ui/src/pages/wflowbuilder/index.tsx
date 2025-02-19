import {
  ActionIcon,
  Anchor,
  Box,
  Breadcrumbs,
  Group,
  Paper,
  ScrollArea,
  Text,
  Title,
} from "@mantine/core";
import { useNavigate, useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useGetProcessor } from "../../apis/processors";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { WorkflowYamlEditor } from "./editor";
import { BlockSearch } from "./blocksearch";
import { ContainerOverlay } from "../../components/Overlay";
import { useState } from "react";
import { IconChevronLeft } from "@tabler/icons-react";

const WorkflowBuilderPage = () => {
  const { workspaceId, processorId } = useParams();
  const [editor, setEditor] = useState<any>(null);
  const navigate = useNavigate();

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
      <Breadcrumbs separator=">">
        <Anchor href={`/`}>Workspaces</Anchor>
        <Anchor href={`/workspace/${workspaceId}/processors`}>
          Processors
        </Anchor>
        <Text>{processor?.name}</Text>
      </Breadcrumbs>
      <Group mt="xs" mb="xl">
        <ActionIcon
          variant="default"
          radius="xl"
          size="sm"
          onClick={() => navigate(`/workspace/${workspaceId}/processors`)}
        >
          <IconChevronLeft size={16} />
        </ActionIcon>
        <Title order={3}>Workflow builder</Title>
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
