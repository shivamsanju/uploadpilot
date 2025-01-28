import { Handle, Position, useNodeConnections } from '@xyflow/react';
import { Box, Group, Paper, ThemeIcon, Text, Indicator } from '@mantine/core';
import { IconAlertTriangle, IconFileTypePdf } from '@tabler/icons-react';
import classes from './BaseNode.module.css';

export const BaseNode = (node: any) => {
    const conn = useNodeConnections({
        id: node.id,
        handleId: 'target',
    });

    return (
        <Indicator label={<IconAlertTriangle size={12} />} size={20} offset={7} disabled>
            <Paper p="sm" withBorder className={classes.node} w={300} radius="sm">
                <Handle type="target" position={Position.Top} className={classes.handle} id="target" isConnectable={conn.length === 0} />
                <Group wrap='nowrap'>
                    <ThemeIcon variant='default' size={40} c="dimmed">
                        <IconFileTypePdf size={18} />
                    </ThemeIcon>
                    <Box>
                        <Text size="sm">{node?.data?.label}</Text>
                        <Text opacity={0.7} size="xs" lineClamp={1}>
                            {node?.data?.description}
                        </Text>
                    </Box>
                </Group>
                <Handle type="source" position={Position.Bottom} className={classes.handle} id="source" />
            </Paper>
        </Indicator>
    );
}