import {
  Box,
  Button,
  Group,
  ScrollArea,
  Text,
  TextInput,
  ThemeIcon,
  Title,
} from '@mantine/core';
import { IconSearch } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useGetAllProcessingTasks } from '../../../apis/processors';
import { getBlockIcon } from '../../../utils/blockicon';
import classes from './BlockSearch.module.css';

export const BlockSearch = ({
  processorId,
  editor,
}: {
  processorId: string;
  editor: any;
}) => {
  const { workspaceId } = useParams();
  const { isPending, error, blocks } = useGetAllProcessingTasks(
    workspaceId || '',
  );
  const [filtered, setFiltered] = useState<any[]>(blocks || []);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const searchText = e.target.value;
    if (isPending || !blocks || error) return [];

    const filteredBlocks = blocks.filter(
      (c: any) =>
        c.name?.toLowerCase()?.includes(searchText.toLowerCase()) ||
        c.description?.toLowerCase()?.includes(searchText.toLowerCase()),
    );

    setFiltered(filteredBlocks);
  };

  const handleAddTask = (item: any) => {
    if (!item || !editor) return;
    // Get base indentation from the current line

    const model = editor.getModel();
    const selection = editor.getSelection();
    const position = selection.getStartPosition();
    const lineContent = model.getLineContent(position.lineNumber);
    const baseIndentation = lineContent.match(/^\s*/)[0];

    // Multiline text to insert
    const text = item?.workflow
      ?.split('\n')
      ?.map((line: string, index: number) =>
        index === 0 ? line : baseIndentation + line,
      ) // preserve indentation for all lines except the first one
      ?.join('\n');

    const id = { major: 1, minor: 1 };
    const op = {
      identifier: id,
      range: {
        startLineNumber: selection?.selectionStartLineNumber || 1,
        startColumn: selection?.selectionStartColumn || 1,
        endLineNumber: selection?.endLineNumber || 1,
        endColumn: selection?.endColumn || 1,
      },
      text,
      forceMoveMarkers: true,
    };
    editor.executeEdits('my-source', [op]);
    editor.focus();
  };

  useEffect(() => {
    if (blocks) {
      setFiltered(blocks);
    }
  }, [blocks]);

  return (
    <Box>
      <Group justify="space-between" gap="sm" align="center" mb="md">
        <Box p="xs">
          <Title order={4} opacity={0.8}>
            Add activities
          </Title>
          <Text fz="xs" c="dimmed">
            Search and select activities
          </Text>
        </Box>
        <TextInput
          placeholder="Search activities"
          leftSection={<IconSearch size={18} stroke={1.5} />}
          onChange={handleSearch}
          autoFocus
          w={{
            base: '100%',
            md: 400,
          }}
          flex={{
            grow: 1,
          }}
        />
      </Group>

      <ScrollArea h="60vh" scrollbarSize={5}>
        {filtered?.length > 0 &&
          filtered.map((item: any, index: number) => (
            <Group
              className={classes.block}
              key={index}
              wrap="nowrap"
              p="lg"
              justify="space-between"
            >
              <Group flex={1}>
                <ThemeIcon variant="default" size={80} c="dimmed">
                  {getBlockIcon(item?.name, 50)}
                </ThemeIcon>
                <Box w="70%">
                  <Text size="sm" mb={5}>
                    {item?.displayName}
                  </Text>
                  <Text opacity={0.7} size="xs" lineClamp={3}>
                    {item?.description}
                  </Text>
                </Box>
              </Group>
              <Button
                variant="filled"
                onClick={() => handleAddTask(item)}
                className={classes.addBtn}
              >
                Add
              </Button>
            </Group>
          ))}
      </ScrollArea>
    </Box>
  );
};
