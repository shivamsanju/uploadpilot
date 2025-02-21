import { Box, Group, Title } from "@mantine/core";
import WorkflowsList from "./List";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";

const WorkflowsPage = () => {
  const { workspaceId } = useParams();

  if (!workspaceId) {
    return <AppLoader h="70vh" />;
  }

  return (
    <Box mb={50}>
      <Group justify="space-between" mb="xl">
        <Title order={3} opacity={0.7}>
          Processors
        </Title>
      </Group>
      <WorkflowsList />
    </Box>
  );
};

export default WorkflowsPage;
