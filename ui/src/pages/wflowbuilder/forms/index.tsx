import { Box, Divider, Group, Text } from "@mantine/core";
import { WebhookTaskForm } from "./Webhook";
import { PropertiesForm } from "./PropertiesForm";
import { IconX } from "@tabler/icons-react";
import { useWorkflowBuilder } from "../../../context/WflowEditorContext";
import { Task } from "../../../types/tasks";
import React from "react";

export type TaskDataFormProps = {
  selectedTask: Task;
};

const getSelectedTaskForm = (
  key: string | undefined
): React.FC<TaskDataFormProps> | null => {
  switch (key) {
    case "Webhook":
      return WebhookTaskForm;
    default:
      return null;
  }
};

export const SelectedTaskForm = () => {
  const { selectedTask, setSelectedTask } = useWorkflowBuilder();
  const DataForm = getSelectedTaskForm(selectedTask?.key);

  if (!selectedTask) return <></>;

  return (
    <Box h="70vh" p="sm">
      <Group justify="space-between" align="center" mb="lg">
        <Text size="lg" fw={500}>
          {selectedTask.name}
        </Text>
        <IconX
          size={18}
          onClick={() => setSelectedTask(null)}
          cursor="pointer"
        />
      </Group>
      <Divider label="Properties" labelPosition="center" />
      <PropertiesForm />
      {DataForm && (
        <Box pb="lg">
          <Divider label="Data" labelPosition="center" mt="md" />
          <DataForm selectedTask={selectedTask} />
        </Box>
      )}
    </Box>
  );
};
