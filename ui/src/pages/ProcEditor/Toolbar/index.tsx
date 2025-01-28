import { IconSearch, IconGripVertical, IconCube3dSphere } from '@tabler/icons-react';
import {
    Box,
    Stack,
    Title,
    TextInput,
    Group,
    Text,
    ThemeIcon,
    ScrollArea,
} from '@mantine/core';
import classes from './Toolbar.module.css';
import { useCallback, useMemo, useState } from 'react';
import { useGetAllProcBlocks } from '../../../apis/processors';
import { useParams } from 'react-router-dom';
import { useDragAndDrop } from '../../../hooks/DndCanvas';


export const Toolbar = () => {
    const { workspaceId } = useParams();
    const [dragging, setDragging] = useState<string | null>(null);
    const [searchText, setSearchText] = useState('');
    const { setType, setDataTransfer } = useDragAndDrop();
    const { isPending, error, blocks } = useGetAllProcBlocks(workspaceId || "");

    const handleDragStart = useCallback((event: any, item: any) => {
        event.dataTransfer.effectAllowed = 'move';
        event.dataTransfer.setData('text/plain', 'dragged-item');
        setType("baseNode");
        setDataTransfer(item)
        setDragging(item?.key);
    }, [setType, setDragging, setDataTransfer]);

    const handleDragEnd = useCallback(() => {
        setDragging(null);
    }, [setDragging]);

    const handleSearch = (e: any) => {
        const text = e.target.value;
        setSearchText(text);
    }

    const items = useMemo(() => {
        if (isPending || !blocks || error) return [];

        const filteredBlocks = blocks.filter((c: any) => {
            return c.key.toLowerCase().includes(searchText.toLowerCase()) ||
                c.description.toLowerCase().includes(searchText.toLowerCase());
        });

        return filteredBlocks.map(({ key, label, description }: any) => (
            <div
                key={key}
                draggable
                onDragStart={(e) => handleDragStart(e, { key, label, description })}
                onDragEnd={handleDragEnd}
                className={classes.block}
            >

                <Group pl="md" wrap='nowrap' className={`${dragging === key ? classes.dragging : ''}`}>
                    <IconGripVertical size={25} stroke={1.5} />
                    <ThemeIcon variant='default' size={40} c="dimmed">
                        <IconCube3dSphere size={18} />
                    </ThemeIcon>
                    <Box w="70%">
                        <Text size="sm">{label}</Text>
                        <Text opacity={0.7} size="xs" lineClamp={1}>
                            {description}
                        </Text>
                    </Box>
                </Group>
            </div>
        ))
    }, [dragging, handleDragEnd, handleDragStart, isPending, blocks, error, searchText]);


    return (
        <Stack mt="sm" className={classes.toolbarContainer}>
            <Title order={4} c="gray" px="md">Tools </Title>
            <TextInput
                px="md"
                placeholder="Search blocks"
                size="sm"
                leftSection={<IconSearch size={12} stroke={1.5} />}
                onChange={handleSearch}
                value={searchText}
                styles={{ section: { pointerEvents: 'none' } }}
                mb="sm"
            />
            <ScrollArea scrollbarSize={5} pb="sm">
                <div className={classes.collections}>{items}</div>
            </ScrollArea>
        </Stack>
    );
}