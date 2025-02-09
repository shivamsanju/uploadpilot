import { useForm } from "@mantine/form";
import { Button, Group, Stack, ScrollArea } from "@mantine/core";
import { IconDeviceFloppy } from "@tabler/icons-react";
import { NodeFormProps } from ".";
import { CommonForm } from "./Common";

const NoDataForm: React.FC<NodeFormProps> = ({ nodeData, saveNodeData }) => {
  const form = useForm({
    initialValues: {
      retry: nodeData?.retry || 0,
      timeoutMilSec: nodeData?.timeoutMilSec || 0,
      continueOnError: nodeData?.continueOnError || false,
    },
  });

  return (
    <form onSubmit={form.onSubmit(saveNodeData)}>
      <Stack justify="space-between" h="82vh">
        <ScrollArea scrollbarSize={6}>
          <Stack gap="lg" pr="lg">
            <CommonForm form={form} />
          </Stack>
        </ScrollArea>
      </Stack>
      <Group justify="center">
        <Button type="submit" leftSection={<IconDeviceFloppy size={16} />}>
          Save
        </Button>
      </Group>
    </form>
  );
};

export default NoDataForm;
