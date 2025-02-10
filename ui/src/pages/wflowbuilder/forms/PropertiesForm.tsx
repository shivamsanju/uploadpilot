import {
  Group,
  Stack,
  Switch,
  Text,
  NumberInput,
  TextInput,
} from "@mantine/core";
import { useWorkflowBuilder } from "../../../context/WflowEditorContext";
import { EditableProperties } from "../../../types/tasks";
import { useEffect, useState } from "react";
import classes from "./forms.module.css";

export const PropertiesForm = () => {
  const { selectedTask, modifyProperties } = useWorkflowBuilder();
  const [properties, setProperties] = useState<EditableProperties>({
    name: selectedTask?.name || selectedTask?.label || "",
    retries: 0,
    timeoutMs: 0,
    continueOnError: false,
  });

  useEffect(() => {
    setProperties({
      name: selectedTask?.name || selectedTask?.label || "",
      retries: selectedTask?.retries || 0,
      timeoutMs: selectedTask?.timeoutMs || 0,
      continueOnError: selectedTask?.continueOnError || false,
    });
  }, [selectedTask]);

  useEffect(() => {
    if (!selectedTask) return;
    modifyProperties(selectedTask.id, properties);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [properties]);

  return (
    <Stack gap="lg">
      <TextInput
        label="Name"
        description="Name of the task"
        placeholder="Enter the name"
        onChange={(e) => setProperties({ ...properties, name: e.target.value })}
        value={properties.name}
      />

      <NumberInput
        label="Retry Count"
        description="Number of retries for the task"
        placeholder="Enter the retry count"
        defaultValue={0}
        onChange={(e) =>
          setProperties({
            ...properties,
            retries: isNaN(Number(e)) ? 0 : Number(e),
          })
        }
        value={properties.retries}
      />

      <NumberInput
        label="Timeout MilSec"
        description="Timeout for the task in milliseconds"
        placeholder="Enter the timeout in milliseconds"
        defaultValue={0}
        onChange={(e) =>
          setProperties({
            ...properties,
            timeoutMs: isNaN(Number(e)) ? 0 : Number(e),
          })
        }
        value={properties.timeoutMs}
      />

      <Group justify="space-between" mt="xs">
        <div>
          <Text fw={500}>Continue on error </Text>
          <Text c="dimmed">
            Continue executing next tasks even if the task fails
          </Text>
        </div>
        <Switch
          className={classes.customSwitch}
          onLabel="ON"
          offLabel="OFF"
          defaultChecked={false}
          onChange={(e) =>
            setProperties({ ...properties, continueOnError: e.target.checked })
          }
          checked={properties.continueOnError}
        />
      </Group>
    </Stack>
  );
};
