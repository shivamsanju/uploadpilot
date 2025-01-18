import { Button, ButtonProps } from '@mantine/core';
import { IconBrandGoogle } from '@tabler/icons-react';

export function GoogleButton(props: ButtonProps & React.ComponentPropsWithoutRef<'button'>) {
    return (
        <Button leftSection={<IconBrandGoogle size={16} />} variant="default" {...props} />
    );
}