import { FormValidateInput } from "@mantine/form";
import { Task } from "../../../../types/tasks";
import { webhookValidations } from "../Webhook";
import { ActivityInvocation } from "../../../../types/workflow";

export type IValidations<T> = FormValidateInput<T>;

const getTaskValidations = (key: string): IValidations<any> => {
  const validations: IValidations<any> = {};

  switch (key) {
    case "Webhook":
      return webhookValidations;
    default:
      return validations;
  }
};

export const validateTask = (task: Task): Task => {
  const validations = getTaskValidations(task.key);
  for (const [key, value] of Object.entries(validations)) {
    const err = value ? value(task.data[key]) : null;
    if (err) {
      task.hasErrors = true;
      task.error = err;
    }
  }

  return task;
};

export const validateTasks = (tasks: Task[]): Task[] => {
  for (let i = 0; i < tasks.length; i++) {
    tasks[i] = validateTask(tasks[i]);
  }
  return tasks;
};

export const validateActivity = (
  activity: ActivityInvocation,
  variables: Record<string, string>
): ActivityInvocation => {
  const validations = getTaskValidations(activity.key);
  for (const [key, value] of Object.entries(validations)) {
    const err = value ? value(variables[activity.id + key]) : null;
    if (err) {
      activity.hasErrors = true;
      activity.error = err;
    }
  }

  return activity;
};

export const validateActivities = (
  activities: ActivityInvocation[],
  variables: Record<string, string>
) => {
  for (let i = 0; i < activities.length; i++) {
    activities[i] = validateActivity(activities[i], variables);
  }
  return activities;
};
