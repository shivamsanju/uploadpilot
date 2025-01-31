import { Button, ButtonProps, } from '@mantine/core';
import { IconRefresh } from '@tabler/icons-react';
import classes from "./RefreshButton.module.css";
import { useDisclosure, useTimeout } from '@mantine/hooks';


export const RefreshButton: React.FC<React.HTMLAttributes<HTMLButtonElement> & ButtonProps> = (props) => {
    const [rotate, handlers] = useDisclosure(false);
    const { start } = useTimeout(() => {
        handlers.close();
    }, 1000);



    const handleRefresh = (e: React.MouseEvent<HTMLButtonElement>) => {
        if (rotate) return;
        handlers.open();
        start();
        if (props.onClick) props.onClick(e);
    }


    return (
        <Button
            className={classes.refreshBtn}
            variant="subtle"
            leftSection={
                <IconRefresh
                    size={15}
                    className={`${rotate ? classes.rotate : ''}`}
                />
            }
            {...props}
            onClick={handleRefresh}
        >
            Refresh
        </Button>
    );
}

