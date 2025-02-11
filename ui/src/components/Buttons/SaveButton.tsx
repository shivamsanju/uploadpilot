import { Button, ButtonProps } from "@mantine/core";
import { IconDeviceFloppy } from "@tabler/icons-react";

export const SaveButton: React.FC<
  ButtonProps & React.ButtonHTMLAttributes<HTMLButtonElement>
> = (props) => (
  <Button leftSection={<IconDeviceFloppy size={18} />} {...props}>
    Save
  </Button>
);
