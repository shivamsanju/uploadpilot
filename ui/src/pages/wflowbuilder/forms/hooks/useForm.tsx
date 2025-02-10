import { useEffect } from "react";
import { useWorkflowBuilder } from "../../../../context/WflowEditorContext";
import { useForm } from "@mantine/form";
import { IValidations } from "../validations";
import { Task } from "../../../../types/tasks";

export const useTaskForm = <T extends Record<string, any>>(
  selectedTask: Task,
  validations: IValidations<T>
) => {
  const { modifyTaskData, setTaskErrors } = useWorkflowBuilder();

  // Ensure useForm uses the generic type T correctly
  const form = useForm<T>({
    initialValues: selectedTask.data as T, // Ensure correct type assertion
    validateInputOnBlur: true,
    validateInputOnChange: true,
    validate: validations,
    onValuesChange: (values) => {
      setTaskErrors(selectedTask.id, !form.isValid(), "");
      modifyTaskData(selectedTask.id, values);
    },
  });

  useEffect(() => {
    if (!selectedTask) return;
    form.setValues(selectedTask.data as T);
    form.validate();
    setTaskErrors(selectedTask.id, !form.isValid(), "");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedTask.id]);

  return form;
};
