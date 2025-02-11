import { Button } from "@mantine/core";
import { IconRestore } from "@tabler/icons-react";
import type { ButtonProps } from "@mantine/core";

export const DiscardButton: React.FC<
  ButtonProps & React.ButtonHTMLAttributes<HTMLButtonElement>
> = (props) => (
  <Button
    leftSection={<IconRestore size={18} />}
    variant="default"
    c="dimmed"
    {...props}
  >
    Discard
  </Button>
);
