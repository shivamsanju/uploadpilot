import type { ButtonProps } from '@mantine/core';
import { Button, Text, TextProps } from '@mantine/core';
import { Icon, IconProps, IconRestore } from '@tabler/icons-react';

export const DiscardButton: React.FC<
  ButtonProps &
    React.ButtonHTMLAttributes<HTMLButtonElement> & {
      iconProps?: React.ForwardRefExoticComponent<
        IconProps & React.RefAttributes<Icon>
      >;
      labelProps?: TextProps;
    }
> = props => (
  <Button
    leftSection={<IconRestore size={18} {...props.iconProps} />}
    variant="default"
    c="dimmed"
    {...props}
  >
    <Text {...props.labelProps}>Discard</Text>
  </Button>
);
