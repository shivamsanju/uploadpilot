import { IconSearch } from "@tabler/icons-react";
import {
  Box,
  Group,
  Text,
  ThemeIcon,
  ScrollArea,
  TextInput,
  Title,
  Button,
} from "@mantine/core";
import { useState, useEffect } from "react";
import { useGetAllProcBlocks } from "../../../apis/processors";
import { useParams } from "react-router-dom";
import classes from "./BlockSearch.module.css";
import { getBlockIcon } from "../../../utils/blockicon";
import { useWorkflowBuilderV2 } from "../../../context/WflowEditorContextV2";

export const BlockSearch = ({ processorId }: { processorId: string }) => {
  const { workspaceId } = useParams();
  const { isPending, error, blocks } = useGetAllProcBlocks(workspaceId || "");
  const [filtered, setFiltered] = useState<any[]>(blocks || []);
  const { addActivity } = useWorkflowBuilderV2();

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const searchText = e.target.value;
    if (isPending || !blocks || error) return [];

    const filteredBlocks = blocks.filter(
      (c: any) =>
        c.key.toLowerCase().includes(searchText.toLowerCase()) ||
        c.description.toLowerCase().includes(searchText.toLowerCase())
    );

    setFiltered(filteredBlocks);
  };

  const handleAddTask = (item: any) => {
    addActivity({
      id: crypto.randomUUID(),
      key: item?.key,
      label: item?.label,
      retries: 0,
      timeout: 0,
      continueOnError: false,
      arguments: [],
      result: "",
    });
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
            Add Tasks
          </Title>
          <Text fz="xs" c="dimmed">
            Search and select blocks
          </Text>
        </Box>
        <TextInput
          placeholder="Search blocks"
          leftSection={<IconSearch size={18} stroke={1.5} />}
          onChange={handleSearch}
          autoFocus
          w={{
            base: "100%",
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
                  {getBlockIcon(item?.key, 50)}
                </ThemeIcon>
                <Box w="70%">
                  <Text size="sm" mb={5}>
                    {item?.label}
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
