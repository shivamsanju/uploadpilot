import { useForm } from "@mantine/form";
import {
  useAddUserToWorkspaceMutation,
  useEditUserInWorkspaceMutation,
} from "../../../apis/workspace";
import { Button, Group, Select, Stack, TextInput } from "@mantine/core";
import { IconAt, IconLockAccess } from "@tabler/icons-react";
import { ErrorLoadingWrapper } from "../../../components/ErrorLoadingWrapper";

const ROLES = ["Owner", "Contributor", "Viewer"];

type Props = {
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
  workspaceId: string;
  mode?: "edit" | "add" | "view";
  setMode: React.Dispatch<React.SetStateAction<"edit" | "add" | "view">>;
  initialValues?: any;
  setInitialValues: React.Dispatch<React.SetStateAction<any>>;
};

const AddUserForm: React.FC<Props> = ({
  setOpened,
  setInitialValues,
  setMode,
  workspaceId,
  mode = "add",
  initialValues,
}) => {
  const form = useForm({
    initialValues:
      (mode === "edit" || mode === "view") && initialValues
        ? initialValues
        : {
            email: "",
            role: "",
          },
    validate: {
      email: (value) =>
        /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(value)
          ? null
          : "Invalid email",
      role: (value) =>
        value && ROLES.includes(value)
          ? null
          : "Role is required and must be one of the following: Owner, Contributor, Viewer",
    },
  });

  const { mutateAsync, isPending } = useAddUserToWorkspaceMutation();
  const { mutateAsync: editUser } = useEditUserInWorkspaceMutation();

  const handleAdd = async (values: any) => {
    try {
      await mutateAsync({
        workspaceId,
        email: values.email,
        role: values.role,
      });
      setOpened(false);
      setInitialValues(null);
      setMode("add");
    } catch (error) {
      console.error(error);
    }
  };

  const handleEdit = async (values: any) => {
    try {
      await editUser({
        workspaceId,
        userId: initialValues.userId,
        role: values.role,
      });
      setOpened(false);
      setInitialValues(null);
      setMode("add");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <ErrorLoadingWrapper error={null} isPending={isPending}>
      <form onSubmit={form.onSubmit(mode === "edit" ? handleEdit : handleAdd)}>
        <Stack gap="xl">
          <TextInput
            leftSection={<IconAt size={16} />}
            withAsterisk
            label="Email"
            description="Email address of the user"
            placeholder="Email"
            disabled={mode === "edit" || mode === "view"}
            {...form.getInputProps("email")}
          />
          <Select
            withAsterisk
            leftSection={<IconLockAccess size={16} />}
            label="Role"
            description="Role to be assigned to the user"
            placeholder="Role"
            data={ROLES}
            disabled={mode === "view"}
            {...form.getInputProps("role")}
          />
        </Stack>
        {mode !== "view" && (
          <Group justify="flex-end" mt={50}>
            <Button type="submit">
              {mode === "edit" && "Save"}
              {mode === "add" && "Add"}
            </Button>
          </Group>
        )}
      </form>
    </ErrorLoadingWrapper>
  );
};

export default AddUserForm;
