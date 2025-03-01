import { CodeHighlight } from '@mantine/code-highlight';
import {
  Accordion,
  Box,
  Group,
  ScrollArea,
  Text,
  TextInput,
} from '@mantine/core';
import { IconIndentIncrease, IconSearch } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useGetAllProcessingActivities } from '../../../apis/processors';
import classes from './BlockSearch.module.css';

export const BlockSearch = () => {
  const { workspaceId } = useParams();
  const { isPending, error, blocks } = useGetAllProcessingActivities(
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

  useEffect(() => {
    if (blocks) {
      setFiltered(blocks);
    }
  }, [blocks]);

  return (
    <Box>
      <TextInput
        m="sm"
        placeholder="Search activities"
        leftSection={<IconSearch size={18} stroke={1.5} />}
        onChange={handleSearch}
        autoFocus
      />
      <ScrollArea h="67vh" scrollbarSize={5}>
        <Accordion chevronPosition="right" maw="100%">
          {filtered?.length > 0 &&
            filtered.map((item: any, index: number) => (
              <Accordion.Item
                value={item.name}
                key={index}
                className={classes.blockItem}
              >
                <Accordion.Control icon={<IconIndentIncrease size={18} />}>
                  <Block {...item} />
                </Accordion.Control>
                <Accordion.Panel maw="100%">
                  <BlockDescription {...item} />
                </Accordion.Panel>
              </Accordion.Item>
            ))}
        </Accordion>
      </ScrollArea>
    </Box>
  );
};

interface BlockProps {
  displayName: string;
  image?: string;
}

export const Block = ({ displayName, image }: BlockProps) => {
  if (!image) {
    image = 'https://img.icons8.com/clouds/256/000000/homer-simpson.png';
  }
  return (
    <Group wrap="nowrap">
      <Text size="sm">{displayName}</Text>
    </Group>
  );
};

export const BlockDescription = ({ description, workflow }: any) => {
  return (
    <Box>
      <Text fz="xs" c="dimmed" mb="sm">
        {description}
      </Text>
      <CodeHighlight
        maw="100%"
        code={workflow}
        language="yaml"
        copyLabel="Copy"
        copiedLabel="Copied"
      />
    </Box>
  );
};
