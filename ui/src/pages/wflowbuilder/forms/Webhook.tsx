import { Stack, TextInput } from "@mantine/core";
import { IconHttpPost, IconLock } from "@tabler/icons-react";
import { TaskDataFormProps } from ".";
import { IValidations } from "./validations";
import { useTaskForm } from "./hooks/useForm";

export type WebhookFormData = { url: string };

export const webhookValidations: IValidations<WebhookFormData> = {
  url: (value: string) => (/^http(s)?:\/\//.test(value) ? null : "Invalid URL"),
};

export const WebhookTaskForm: React.FC<TaskDataFormProps> = ({
  selectedTask,
}) => {
  const form = useTaskForm<WebhookFormData>(selectedTask, webhookValidations);

  return (
    <form>
      <Stack gap="lg">
        <TextInput
          leftSection={<IconHttpPost size={16} color="#E0A526" />}
          withAsterisk
          label="Target URL"
          type="url"
          description="A post request will be sent to this URL with zip file url in body. {url: 'https://example.com/file.zip'}"
          placeholder="Enter the webhook url"
          {...form.getInputProps("url")}
        />

        <TextInput
          leftSection={<IconLock size={16} />}
          label="Signing Secret"
          description="Signing secret for the webhook"
          placeholder="Enter the signing secret"
          {...form.getInputProps("secret")}
        />
      </Stack>
    </form>
  );
};
