import { ActionIcon, Box, Group, Title } from "@mantine/core";
import WorkflowRunsList from "./List";
import { IconPlus } from "@tabler/icons-react";
import { useState } from "react";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";

const WorkflowRunsPage = () => {
  const { workspaceId } = useParams();
  const [workflowAddOpen, setWorkflowAddOpen] = useState(false);

  if (!workspaceId) {
    return <AppLoader h="70vh" />;
  }

  return (
    <Box mb={50}>
      <Group justify="space-between" mb="xl">
        <Title order={3} opacity={0.7}>
          Workflow Runs
        </Title>
        <Box
          style={{
            position: "relative",
            display: "inline-block",
            marginRight: "1rem",
          }}
        >
          <ActionIcon
            size="lg"
            variant="outline"
            onClick={() => setWorkflowAddOpen(true)}
          >
            <IconPlus />
          </ActionIcon>
        </Box>
      </Group>
      <WorkflowRunsList
        opened={workflowAddOpen}
        setOpened={setWorkflowAddOpen}
      />
    </Box>
  );
};

export default WorkflowRunsPage;
