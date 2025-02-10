import { FormValidateInput } from "@mantine/form";
import { Task } from "../../../../types/tasks";
import { webhookValidations } from "../Webhook";

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
