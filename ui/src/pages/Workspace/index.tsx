import {
  ActionIcon,
  Box,
  Button,
  Group,
  Modal,
  Paper,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import {
  useCreateWorkspaceMutation,
  useGetWorkspaces,
} from "../../apis/workspace";
import { useNavigate } from "react-router-dom";
import classes from "./Workspace.module.css";
import {
  IconCategory,
  IconCategoryPlus,
  IconChevronsRight,
  IconPlus,
} from "@tabler/icons-react";
import { ErrorLoadingWrapper } from "../../components/ErrorLoadingWrapper";
import { useState } from "react";

const WorkspaceLandingPage = () => {
  const [opened, toggle] = useState(false);
  const { isPending, error, workspaces } = useGetWorkspaces();
  const { mutateAsync, isPending: isCreating } = useCreateWorkspaceMutation();
  const navigate = useNavigate();

  const form = useForm({
    initialValues: {
      name: "",
      tags: [],
    },
    validate: {
      name: (value) => {
        if (!value.trim()) {
          return "Workspace name is required";
        }
        if (value.trim().length > 20 || value.trim().length < 2) {
          return "Workspace name must be between 2 and 20 characters";
        }
        return null;
      },
    },
    validateInputOnChange: true,
  });

  const handleCreateWorkspace = async (values: any) => {
    try {
      const id = await mutateAsync(values.name);
      navigate(`/workspace/${id}`);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <ErrorLoadingWrapper error={error} isPending={isPending || isCreating}>
      <Group justify="center" mb="50">
        <Box mt={30} w="90%">
          <Paper
            withBorder
            p="md"
            radius="md"
            className={classes.wsItemAdd}
            onClick={() => toggle(true)}
          >
            <Group justify="center" h="100%">
              <IconPlus size={30} stroke={2} color="gray" />
            </Group>
          </Paper>
          {workspaces?.length > 0 &&
            workspaces.map((workspace: any) => (
              <Paper
                withBorder
                p="md"
                radius="md"
                key={workspace.id}
                className={classes.wsItem}
                w="100%"
              >
                <Group justify="space-between" key={workspace.id} h="100%">
                  <Group gap="sm">
                    <IconCategory size={30} stroke={2} color="gray" />
                    <Text size="sm" fw="bold" opacity={0.7}>
                      {workspace.name}
                    </Text>
                  </Group>
                  <ActionIcon
                    size="lg"
                    onClick={() => navigate(`/workspace/${workspace.id}`)}
                  >
                    <IconChevronsRight />
                  </ActionIcon>
                </Group>
              </Paper>
            ))}
        </Box>
        <Modal
          title="Create new workspace"
          size="lg"
          padding="xl"
          transitionProps={{ transition: "pop" }}
          opened={opened}
          onClose={() => toggle(false)}
        >
          <form
            onSubmit={form.onSubmit((values) => handleCreateWorkspace(values))}
          >
            <TextInput
              mb="xl"
              label="Workspace name"
              description="Name of the workspace"
              leftSection={<IconCategoryPlus stroke={2} color="gray" />}
              placeholder="Enter a workspace name"
              {...form.getInputProps("name")}
            />
            <Group justify="flex-end">
              <Button type="submit">Create</Button>
            </Group>
          </form>
        </Modal>
      </Group>
    </ErrorLoadingWrapper>
  );
};

export default WorkspaceLandingPage;
