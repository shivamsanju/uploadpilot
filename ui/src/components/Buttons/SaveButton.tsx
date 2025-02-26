import { Button, ButtonProps, Text, TextProps } from '@mantine/core';
import { Icon, IconDeviceFloppy, IconProps } from '@tabler/icons-react';

export const SaveButton: React.FC<
  ButtonProps &
    React.ButtonHTMLAttributes<HTMLButtonElement> & {
      iconProps?: React.ForwardRefExoticComponent<
        IconProps & React.RefAttributes<Icon>
      >;
      labelProps?: TextProps;
    }
> = props => (
  <Button
    leftSection={<IconDeviceFloppy size={18} {...props.iconProps} />}
    {...props}
  >
    <Text {...props.labelProps}>Save</Text>
  </Button>
);
