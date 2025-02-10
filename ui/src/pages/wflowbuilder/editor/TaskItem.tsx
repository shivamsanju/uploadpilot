import {
  IconAlertTriangle,
  IconGripVertical,
  IconX,
} from "@tabler/icons-react";
import { motion } from "framer-motion";
import cx from "clsx";
import { Box, Group, Text } from "@mantine/core";
import classes from "./editor.module.css";
import { Task } from "../../../types/tasks";
import { useWorkflowBuilder } from "../../../context/WflowEditorContext";
import { getBlockIcon } from "../../../utils/blockicon";
import { showConfirmationPopup } from "../../../components/Popups/ConfirmPopup";

type TaskItemProps = {
  item: Task;
  provided: any;
  snapshot: any;
};

export const TaskItem: React.FC<TaskItemProps> = ({
  item,
  provided,
  snapshot,
}) => {
  const { selectedTask, removeTask, setSelectedTask } = useWorkflowBuilder();

  const handleRemoveTask = () => {
    showConfirmationPopup({
      message: "Are you sure you want to remove this task?",
      onOk: () => {
        removeTask(item.id);
      },
    });
  };

  const handleSelectTask = () => {
    if (selectedTask?.id !== item.id) {
      setSelectedTask(item);
    }
  };

  return (
    <motion.div
      initial={{ x: "100%", opacity: 0 }}
      animate={{ x: 0, opacity: 1 }}
      exit={{ x: "100%", opacity: 0 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
    >
      <Group
        justify="space-between"
        className={cx(classes.item, {
          [classes.itemDragging]: snapshot.isDragging,
          [classes.selected]: item.id === selectedTask?.id,
        })}
        {...provided.draggableProps}
        ref={provided.innerRef}
      >
        <Box flex={1} onClick={handleSelectTask} p="sm">
          <Group gap="md" align="center">
            <Box c="dimmed">{getBlockIcon(item.key, 25)}</Box>
            <Box>
              <Text flex={1}>{item.name}</Text>
              {item.hasErrors ? (
                <Text c="red" style={{ fontSize: "11px" }} mt={2} fz="xs">
                  <Group align="center" gap={4}>
                    <IconAlertTriangle size={12} stroke={1.5} />
                    Some settings are missing
                  </Group>
                </Text>
              ) : (
                <Text fz="xs" c="dimmed" mt={2}>
                  {item.label}
                </Text>
              )}
            </Box>
          </Group>
        </Box>
        <Group gap="xs" align="center" p="sm">
          <IconX
            className={classes.deleteBtn}
            onClick={handleRemoveTask}
            color="red"
            size={25}
            stroke={1.5}
          />
          <Box
            {...provided.dragHandleProps}
            className={classes.dragHandle}
            c="dimmed"
          >
            <IconGripVertical size={30} stroke={1.5} />
          </Box>
        </Group>
      </Group>
    </motion.div>
  );
};
