import React, { createContext, useContext, useReducer, ReactNode } from "react";
import {
  ActivityInvocation,
  EditableProperties,
  Statement,
} from "../types/workflow";

type State = {
  variables: Record<string, string>;
  activities: ActivityInvocation[];
  selectedActivity: ActivityInvocation | null;
  newActivityId: string | null;
};

type Action =
  | {
      type: "SET_ACTIVITIES_AND_VARIABLES";
      statement: Statement | null;
      variables: Record<string, string>;
    }
  | { type: "REORDER_ACTIVITY"; from: number; to: number }
  | { type: "ADD_ACTIVITY"; activity: ActivityInvocation }
  | { type: "REMOVE_ACTIVITY"; id: string }
  | {
      type: "EDIT_ACTIVITY";
      activityId: string;
      properties: EditableProperties;
    }
  | {
      type: "SET_ACTIVITY_ERRORS";
      activityId: string;
      hasErrors: boolean;
      error: string;
    }
  | {
      type: "EDIT_ACTIVITY_VARIABLES";
      id: string;
      variables: Record<string, string>;
    }
  | { type: "SET_SELECTED_ACTIVITY"; activity: ActivityInvocation | null };

const WorkflowBuilderContextV2 = createContext<{
  state: State;
  dispatch: React.Dispatch<Action>;
} | null>(null);

const initialState: State = {
  variables: {},
  activities: [],
  selectedActivity: null,
  newActivityId: null,
};

const workflowReducer = (state: State, action: Action): State => {
  switch (action.type) {
    case "SET_ACTIVITIES_AND_VARIABLES":
      const activities =
        action?.statement?.sequence?.elements
          ?.map((s) => s.activity)
          ?.filter((e) => e !== undefined) || [];

      return {
        ...state,
        variables: action.variables,
        activities: activities.map((a) => a as ActivityInvocation),
      };
    case "REORDER_ACTIVITY":
      const updatedActivities = [...state.activities];
      const [movedActivity] = updatedActivities.splice(action.from, 1);
      updatedActivities.splice(action.to, 0, movedActivity);

      return {
        ...state,
        activities: updatedActivities,
      };
    case "ADD_ACTIVITY":
      return {
        ...state,
        activities: [...state.activities, action.activity],
        newActivityId: action.activity.id,
      };
    case "REMOVE_ACTIVITY":
      return {
        ...state,
        activities: state.activities.filter(
          (activity) => activity.id !== action.id
        ),
      };
    case "EDIT_ACTIVITY":
      return {
        ...state,
        variables: state.variables || {},
        activities: state.activities.map((activity) => {
          if (activity.id === action.activityId) {
            return {
              ...activity,
              ...action.properties,
            };
          }
          return activity;
        }),
      };
    case "SET_ACTIVITY_ERRORS":
      return {
        ...state,
        activities: state.activities.map((activity) => {
          if (activity.id === action.activityId) {
            return {
              ...activity,
              hasErrors: action.hasErrors,
              error: action.error,
            };
          }
          return activity;
        }),
      };
    case "EDIT_ACTIVITY_VARIABLES":
      return {
        ...state,
        activities: state.activities.map((activity) => {
          if (activity.id === action.id) {
            return {
              ...activity,
              arguments: Object.keys(action.variables),
            };
          }
          return activity;
        }),
        variables: { ...state.variables, ...action.variables },
      };
    case "SET_SELECTED_ACTIVITY":
      return {
        ...state,
        selectedActivity: action.activity,
      };
    default:
      return state;
  }
};

type WorkflowBuilderProps = {
  children: ReactNode;
};

export const WorkflowBuilderProviderV2: React.FC<WorkflowBuilderProps> = ({
  children,
}) => {
  const [state, dispatch] = useReducer(workflowReducer, initialState);

  return (
    <WorkflowBuilderContextV2.Provider value={{ state, dispatch }}>
      {children}
    </WorkflowBuilderContextV2.Provider>
  );
};

export const useWorkflowBuilderV2 = () => {
  const context = useContext(WorkflowBuilderContextV2);
  if (!context) {
    throw new Error("useWorkflowEditor must be used within a WorkflowBuilder");
  }
  const { state, dispatch } = context;

  const setActivitiesAndVariables = (
    statement: Statement | null,
    variables: Record<string, string>
  ) => dispatch({ type: "SET_ACTIVITIES_AND_VARIABLES", statement, variables });
  const reorderActivity = (from: number, to: number) =>
    dispatch({ type: "REORDER_ACTIVITY", from, to });
  const addActivity = (activity: ActivityInvocation) =>
    dispatch({ type: "ADD_ACTIVITY", activity });
  const removeActivity = (id: string) =>
    dispatch({ type: "REMOVE_ACTIVITY", id });
  const editActivity = (activityId: string, properties: EditableProperties) =>
    dispatch({ type: "EDIT_ACTIVITY", activityId, properties });
  const setActivityErrors = (
    activityId: string,
    hasErrors: boolean,
    error: string
  ) => dispatch({ type: "SET_ACTIVITY_ERRORS", activityId, hasErrors, error });
  const editActivityVariables = (
    id: string,
    variables: Record<string, string>
  ) => dispatch({ type: "EDIT_ACTIVITY_VARIABLES", id, variables });
  const setSelectedActivity = (activity: ActivityInvocation | null) =>
    dispatch({ type: "SET_SELECTED_ACTIVITY", activity });

  return {
    variables: state.variables,
    activities: state.activities,
    selectedActivity: state.selectedActivity,
    newActivityId: state.newActivityId,
    setActivitiesAndVariables,
    reorderActivity,
    addActivity,
    removeActivity,
    editActivity,
    setActivityErrors,
    editActivityVariables,
    setSelectedActivity,
  };
};
