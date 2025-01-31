import { Handle, NodeToolbar, Position, useNodeConnections } from '@xyflow/react';
import { Box, Group, Paper, ThemeIcon, Text, Grid } from '@mantine/core';
import { IconCircleDashedCheck, IconCirclePlusFilled, IconCube3dSphere } from '@tabler/icons-react';
import classes from './BaseNode.module.css';
import { useCanvas } from '../../context/EditorCtx';

export const BaseNode = (node: any) => {
    const conn = useNodeConnections({
        id: node.id,
        handleId: 'target',
    });

    const { setOpenedNodeId, nodes, onConnectEnd, openBlocksModal } = useCanvas();

    const handleNewNode = (id: string) => {
        onConnectEnd(id)
        openBlocksModal()
    }

    return (

        <Paper p="sm" withBorder className={classes.node} w={300} radius="sm" onDoubleClick={() => setOpenedNodeId(node.id)}>
            <Grid>
                <Grid.Col span={11}>
                    <Group wrap='nowrap'>
                        <ThemeIcon variant='default' size={40} c="dimmed">
                            <IconCube3dSphere size={18} />
                        </ThemeIcon>
                        <Box>
                            <Text size="sm">{node?.data?.label}</Text>
                            <Text opacity={0.7} size="xs" lineClamp={1}>
                                {node?.data?.description}
                            </Text>
                        </Box>
                    </Group>
                </Grid.Col>
                <Grid.Col span={1}>
                    <Box mt="-10">
                        <IconCircleDashedCheck size={18} color={conn.length > 0 ? 'green' : 'dimmed'} />
                    </Box>
                </Grid.Col>
            </Grid>
            <Handle
                type="target"
                position={Position.Top}
                className={classes.targetHandle}
                id="target"
                isConnectable={conn.length === 0 && nodes.length > 1}
            />
            <Handle
                type="source"
                position={Position.Bottom}
                id="source"
                color='red'
                className={classes.sourceHandle}
            />
            <NodeToolbar
                className={classes.nodeToolbar}
                position={Position.Bottom}
                isVisible
                onClick={() => handleNewNode(node.id)}
            >
                <IconCirclePlusFilled />
            </NodeToolbar>
        </Paper>
    );
}