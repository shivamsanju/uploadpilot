import React, { createContext, useContext, useReducer, ReactNode } from "react";
import { EditableProperties, Task } from "../types/tasks";
import {
  validateTasks,
  validateTask,
} from "../pages/wflowbuilder/forms/validations";

type State = {
  tasks: Task[];
  selectedTask: Task | null;
  newTaskId: string | null;
};

type Action =
  | { type: "SET_TASKS"; tasks: Task[] }
  | { type: "REORDER_TASK"; from: number; to: number }
  | { type: "ADD_TASK"; task: Task }
  | { type: "REMOVE_TASK"; id: string }
  | {
      type: "MODIFY_PROPERTIES";
      taskId: string;
      properties: EditableProperties;
    }
  | {
      type: "SET_TASK_ERRORS";
      taskId: string;
      hasErrors: boolean;
      error: string;
    }
  | { type: "MODIFY_TASK_DATA"; id: string; data: Record<string, any> }
  | { type: "SET_SELECTED_TASK"; task: Task | null };

const WorkflowBuilderContext = createContext<{
  state: State;
  dispatch: React.Dispatch<Action>;
} | null>(null);

const initialState: State = {
  tasks: [],
  selectedTask: null,
  newTaskId: null,
};

const workflowReducer = (state: State, action: Action): State => {
  switch (action.type) {
    case "SET_TASKS":
      const validatedTasks = validateTasks(action.tasks);
      return {
        ...state,
        tasks: validatedTasks,
      };
    case "REORDER_TASK":
      const updatedTasks = [...state.tasks];
      const [movedTask] = updatedTasks.splice(action.from, 1);
      updatedTasks.splice(action.to, 0, movedTask);

      const reorderedTasks = updatedTasks.map((task, index) => ({
        ...task,
        position: index + 1,
      }));

      return {
        ...state,
        tasks: reorderedTasks,
      };
    case "ADD_TASK":
      const validatedTask = validateTask(action.task);
      return {
        ...state,
        tasks: [
          ...state.tasks,
          { ...validatedTask, position: state.tasks.length + 1 },
        ],
        newTaskId: action.task.id,
      };
    case "REMOVE_TASK":
      return {
        ...state,
        tasks: state.tasks.filter((task) => task.id !== action.id),
        selectedTask:
          state.selectedTask?.id === action.id ? null : state.selectedTask,
      };
    case "MODIFY_PROPERTIES":
      return {
        ...state,
        tasks: state.tasks.map((task) =>
          task.id === action.taskId ? { ...task, ...action.properties } : task
        ),
      };
    case "SET_TASK_ERRORS":
      return {
        ...state,
        tasks: state.tasks.map((task) =>
          task.id === action.taskId
            ? { ...task, hasErrors: action.hasErrors, error: action.error }
            : task
        ),
      };
    case "MODIFY_TASK_DATA":
      return {
        ...state,
        tasks: state.tasks.map((task) => {
          if (task.id === action.id) {
            return {
              ...task,
              data: { ...task.data, ...action.data },
            };
          }
          return task;
        }),
      };
    case "SET_SELECTED_TASK":
      return {
        ...state,
        selectedTask: action.task,
      };
    default:
      return state;
  }
};

type WorkflowBuilderProps = {
  children: ReactNode;
};

export const WorkflowBuilderProvider: React.FC<WorkflowBuilderProps> = ({
  children,
}) => {
  const [state, dispatch] = useReducer(workflowReducer, initialState);

  return (
    <WorkflowBuilderContext.Provider value={{ state, dispatch }}>
      {children}
    </WorkflowBuilderContext.Provider>
  );
};

export const useWorkflowBuilder = () => {
  const context = useContext(WorkflowBuilderContext);
  if (!context) {
    throw new Error("useWorkflowEditor must be used within a WorkflowBuilder");
  }
  const { state, dispatch } = context;

  const setTasks = (tasks: Task[]) => dispatch({ type: "SET_TASKS", tasks });
  const reorderTasks = (from: number, to: number) =>
    dispatch({ type: "REORDER_TASK", from, to });
  const addTask = (task: Task) => dispatch({ type: "ADD_TASK", task });
  const removeTask = (id: string) => dispatch({ type: "REMOVE_TASK", id });
  const modifyProperties = (taskId: string, properties: EditableProperties) =>
    dispatch({ type: "MODIFY_PROPERTIES", taskId, properties });
  const setTaskErrors = (taskId: string, hasErrors: boolean, error: string) =>
    dispatch({ type: "SET_TASK_ERRORS", taskId, hasErrors, error });
  const modifyTaskData = (id: string, data: Record<string, any>) =>
    dispatch({ type: "MODIFY_TASK_DATA", id, data });
  const setSelectedTask = (task: Task | null) =>
    dispatch({ type: "SET_SELECTED_TASK", task });

  return {
    tasks: state.tasks,
    setTasks,
    reorderTasks,
    addTask,
    removeTask,
    modifyProperties,
    modifyTaskData,
    setTaskErrors,
    selectedTask: state.selectedTask,
    setSelectedTask,
    newTaskId: state.newTaskId,
  };
};
