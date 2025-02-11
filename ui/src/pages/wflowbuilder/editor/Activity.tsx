import {
  IconAlertTriangle,
  IconGripVertical,
  IconX,
} from "@tabler/icons-react";
import { motion } from "framer-motion";
import cx from "clsx";
import { Box, Group, Text } from "@mantine/core";
import classes from "./editor.module.css";
import { getBlockIcon } from "../../../utils/blockicon";
import { showConfirmationPopup } from "../../../components/Popups/ConfirmPopup";
import { useWorkflowBuilderV2 } from "../../../context/WflowEditorContextV2";
import { ActivityInvocation } from "../../../types/workflow";

type ActivityItemProps = {
  item: ActivityInvocation;
  provided: any;
  snapshot: any;
};

export const ActivityItem: React.FC<ActivityItemProps> = ({
  item,
  provided,
  snapshot,
}) => {
  const { selectedActivity, removeActivity, setSelectedActivity } =
    useWorkflowBuilderV2();

  const handleRemoveActivity = () => {
    showConfirmationPopup({
      message: "Are you sure you want to remove this task?",
      onOk: () => {
        removeActivity(item.id);
      },
    });
  };

  const handleSelectActivity = () => {
    if (selectedActivity?.id !== item.id) {
      setSelectedActivity(item);
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
          [classes.selected]: item.id === selectedActivity?.id,
        })}
        {...provided.draggableProps}
        ref={provided.innerRef}
      >
        <Box flex={1} onClick={handleSelectActivity} p="sm">
          <Group gap="md" align="center">
            <Box c="dimmed">{getBlockIcon(item.key, 25)}</Box>
            <Box>
              <Text flex={1}>{item.label}</Text>
              {item.hasErrors ? (
                <Text c="red" style={{ fontSize: "11px" }} mt={2} fz="xs">
                  <Group align="center" gap={4}>
                    <IconAlertTriangle size={12} stroke={1.5} />
                    Some settings need attention
                  </Group>
                </Text>
              ) : (
                <Text fz="xs" c="dimmed" mt={2}>
                  {item.key}
                </Text>
              )}
            </Box>
          </Group>
        </Box>
        <Group gap="xs" align="center" p="sm">
          <IconX
            className={classes.deleteBtn}
            onClick={handleRemoveActivity}
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
