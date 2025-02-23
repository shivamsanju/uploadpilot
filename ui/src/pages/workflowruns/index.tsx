import {
  ActionIcon,
  Anchor,
  Box,
  Breadcrumbs,
  Group,
  Text,
  Title,
} from "@mantine/core";
import WorkflowRunsList from "./List";
import { useNavigate, useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { IconChevronLeft } from "@tabler/icons-react";

const WorkflowRunsPage = () => {
  const { workspaceId, processorId } = useParams();
  const navigate = useNavigate();

  if (!workspaceId) {
    return <AppLoader h="70vh" />;
  }

  return (
    <Box>
      <Breadcrumbs separator=">">
        <Anchor href={`/`}>Workspaces</Anchor>
        <Anchor href={`/workspace/${workspaceId}/processors`}>
          Processors
        </Anchor>
        <Text>{processorId}</Text>
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
        <Title order={3}>Workflow runs</Title>
      </Group>
      <WorkflowRunsList />
    </Box>
  );
};

export default WorkflowRunsPage;
