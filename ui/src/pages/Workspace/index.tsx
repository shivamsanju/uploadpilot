import {
  ActionIcon,
  Box,
  Divider,
  Group,
  LoadingOverlay,
  Paper,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import {
  useCreateWorkspaceMutation,
  useGetWorkspaces,
} from "../../apis/workspace";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useNavigate } from "react-router-dom";
import classes from "./Workspace.module.css";
import { IconCategory, IconChevronsRight } from "@tabler/icons-react";

const WorkspaceLandingPage = () => {
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
  });

  const handleCreateWorkspace = async (values: any) => {
    try {
      const id = await mutateAsync(values.name);
      navigate(`/workspaces/${id}`);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Group justify="center" mb="50">
      {error ? (
        <ErrorCard title="Error" message={error.message} h="55vh" />
      ) : isPending ? (
        <AppLoader h="55vh" />
      ) : (
        <Box mt={30}>
          <LoadingOverlay
            visible={isCreating}
            overlayProps={{ radius: "sm", blur: 1 }}
          />
          <form
            onSubmit={form.onSubmit((values) => handleCreateWorkspace(values))}
          >
            <Paper withBorder p="md" radius="md">
              <Group className={classes.wsForm} justify="space-between">
                <Text size="md" w="100%">
                  Create a new workspace
                </Text>
                <TextInput
                  size="sm"
                  leftSection={<IconCategory stroke={2} color="gray" />}
                  w="80%"
                  placeholder="Enter a workspace name"
                  {...form.getInputProps("name")}
                />
                <ActionIcon type="submit" size="xl" radius={"50%"}>
                  <IconChevronsRight />
                </ActionIcon>
              </Group>
            </Paper>
          </form>
          {workspaces && workspaces.length > 0 && (
            <>
              <Divider mt={50} mb={50} />
              <Stack>
                {workspaces &&
                  workspaces.length > 0 &&
                  workspaces.map((workspace: any) => (
                    <Paper withBorder p="md" radius="md" key={workspace.id}>
                      <Group
                        justify="space-between"
                        key={workspace.id}
                        className={classes.wsItem}
                      >
                        <Group gap="sm">
                          <IconCategory size={30} stroke={2} color="gray" />
                          <Text size="sm" fw="bold" opacity={0.7}>
                            {workspace.name}
                          </Text>
                        </Group>
                        <ActionIcon
                          size="xl"
                          radius={"50%"}
                          onClick={() =>
                            navigate(`/workspaces/${workspace.id}`)
                          }
                        >
                          <IconChevronsRight />
                        </ActionIcon>
                      </Group>
                    </Paper>
                  ))}
              </Stack>
            </>
          )}
        </Box>
      )}
    </Group>
  );
};

export default WorkspaceLandingPage;
