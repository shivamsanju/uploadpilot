import { DragDropContext, Draggable, Droppable } from "@hello-pangea/dnd";
import { AnimatePresence } from "framer-motion";
import { useWorkflowBuilder } from "../../../context/WflowEditorContext";
import { TaskItem } from "./TaskItem"; // Import the new TaskItem component
import { useEffect, useRef } from "react";
import {
  Box,
  Group,
  LoadingOverlay,
  ScrollArea,
  Text,
  Title,
} from "@mantine/core";
import { DiscardButton } from "../../../components/Buttons/DiscardButton";
import { SaveButton } from "../../../components/Buttons/SaveButton";
import classes from "./editor.module.css";
import { useGetTasks, useSaveTasksMutation } from "../../../apis/tasks";
import { showConfirmationPopup } from "../../../components/Popups/ConfirmPopup";

type Props = {
  workspaceId: string;
  processorId: string;
};

export const WorkflowEditor: React.FC<Props> = ({
  workspaceId,
  processorId,
}) => {
  const { tasks, reorderTasks, newTaskId, setTasks } = useWorkflowBuilder();
  const { tasks: existingTasks, isPending } = useGetTasks(
    workspaceId,
    processorId
  );
  const { mutateAsync, isPending: isSaving } = useSaveTasksMutation();

  const scrollBotomRef = useRef<HTMLDivElement>(null);

  const saveTasks = async () => {
    if (!workspaceId || !processorId || !tasks || tasks.length === 0) {
      return;
    }
    try {
      await mutateAsync({
        workspaceId,
        processorId,
        tasks,
      });
    } catch (error) {
      console.error(error);
    }
  };

  const discardTasks = () => {
    if (existingTasks) {
      showConfirmationPopup({
        message:
          "Are you sure you want to discard changes? This action cannot be undone.",
        onOk: () => {
          setTasks(existingTasks);
        },
      });
    }
  };

  useEffect(() => {
    if (newTaskId && scrollBotomRef.current) {
      scrollBotomRef.current.scrollIntoView({
        behavior: "smooth",
        block: "nearest",
      });
    }
  }, [newTaskId]);

  useEffect(() => {
    if (existingTasks) {
      setTasks(existingTasks);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [existingTasks]);

  return (
    <Box>
      <LoadingOverlay
        visible={isPending || isSaving}
        overlayProps={{ backgroundOpacity: 0 }}
        zIndex={1000}
      />
      <Group justify="space-between" align="center" p="xs">
        <Box>
          <Title order={4} opacity={0.8}>
            Steps
          </Title>
          <Text c="dimmed">Drag and drop tasks to build your workflow</Text>
        </Box>
        <Group gap="md">
          <DiscardButton onClick={discardTasks} />
          <SaveButton onClick={saveTasks} />
        </Group>
      </Group>
      <ScrollArea
        h="70vh"
        scrollbarSize={6}
        className={classes.canvasContainer}
      >
        <DragDropContext
          onDragEnd={({ destination, source }) =>
            reorderTasks(source.index, destination?.index || 0)
          }
        >
          <Droppable droppableId="dnd-list" direction="vertical">
            {(provided) => (
              <div {...provided.droppableProps} ref={provided.innerRef}>
                <AnimatePresence>
                  {tasks.map((item, index) => (
                    <Draggable
                      key={item.id}
                      index={index}
                      draggableId={item.id}
                    >
                      {(provided, snapshot) => (
                        <TaskItem
                          item={item}
                          provided={provided}
                          snapshot={snapshot}
                        />
                      )}
                    </Draggable>
                  ))}
                </AnimatePresence>
                {provided.placeholder}
              </div>
            )}
          </Droppable>
        </DragDropContext>
        <div ref={scrollBotomRef} />
      </ScrollArea>
    </Box>
  );
};
