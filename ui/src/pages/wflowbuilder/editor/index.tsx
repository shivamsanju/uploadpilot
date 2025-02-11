import { DragDropContext, Draggable, Droppable } from "@hello-pangea/dnd";
import { AnimatePresence } from "framer-motion";
import { ActivityItem } from "./Activity";
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
import { showConfirmationPopup } from "../../../components/Popups/ConfirmPopup";
import { useWorkflowBuilderV2 } from "../../../context/WflowEditorContextV2";
import { Statement } from "../../../types/workflow";
import { useUpdateProcessorWorkflowMutation } from "../../../apis/processors";

type Props = {
  workspaceId: string;
  processorId: string;
  statement: Statement | null;
  variables: Record<string, string>;
};

export const WorkflowEditor: React.FC<Props> = ({
  workspaceId,
  processorId,
  statement,
  variables: initialVariables,
}) => {
  const {
    activities,
    reorderActivity,
    newActivityId,
    variables,
    setActivitiesAndVariables,
  } = useWorkflowBuilderV2();
  const { mutateAsync, isPending: isSaving } =
    useUpdateProcessorWorkflowMutation();

  const scrollBotomRef = useRef<HTMLDivElement>(null);

  const saveActivities = async () => {
    if (!workspaceId || !processorId) {
      return;
    }
    try {
      await mutateAsync({
        workspaceId,
        processorId,
        workflow: {
          root: {
            sequence: {
              elements: activities.map((a) => ({
                activity: a,
              })),
            },
          },
          variables,
        },
      });
    } catch (error) {
      console.error(error);
    }
  };

  const discardTasks = () => {
    showConfirmationPopup({
      message:
        "Are you sure you want to discard changes? This action cannot be undone.",
      onOk: () => {
        setActivitiesAndVariables(statement, initialVariables);
      },
    });
  };

  useEffect(() => {
    if (newActivityId && scrollBotomRef.current) {
      scrollBotomRef.current.scrollIntoView({
        behavior: "smooth",
        block: "nearest",
      });
    }
  }, [newActivityId]);

  useEffect(() => {
    setActivitiesAndVariables(statement, initialVariables);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [statement, initialVariables]);

  return (
    <Box>
      <LoadingOverlay
        visible={isSaving}
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
          <SaveButton onClick={saveActivities} />
        </Group>
      </Group>
      <ScrollArea
        h="70vh"
        scrollbarSize={6}
        className={classes.canvasContainer}
      >
        <DragDropContext
          onDragEnd={({ destination, source }) =>
            reorderActivity(source.index, destination?.index || 0)
          }
        >
          <Droppable droppableId="dnd-list" direction="vertical">
            {(provided) => (
              <div {...provided.droppableProps} ref={provided.innerRef}>
                <AnimatePresence>
                  {activities.map((item, index) => (
                    <Draggable
                      key={item.id}
                      index={index}
                      draggableId={item.id}
                    >
                      {(provided, snapshot) => (
                        <ActivityItem
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
