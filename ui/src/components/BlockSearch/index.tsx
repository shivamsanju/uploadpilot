import { IconSearch, IconCube3dSphere } from "@tabler/icons-react";
import {
  Box,
  Group,
  Text,
  ThemeIcon,
  ScrollArea,
  Modal,
  TextInput,
} from "@mantine/core";
import { useState, useRef, useEffect } from "react";
import { useGetAllProcBlocks } from "../../apis/processors";
import { useParams } from "react-router-dom";
import classes from "./BlockSearch.module.css";
import { useCanvas } from "../../context/ProcEditorContext";

export const BlockSearch = ({
  opened,
  close,
}: {
  opened: boolean;
  close: () => void;
}) => {
  const { workspaceId } = useParams();
  const { isPending, error, blocks } = useGetAllProcBlocks(workspaceId || "");
  const [focusedIndex, setFocusedIndex] = useState(-1);
  const itemRefs = useRef<(HTMLDivElement | null)[]>([]);
  const [filtered, setFiltered] = useState<any[]>(blocks || []);
  const { onSelectNewNode } = useCanvas();

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const searchText = e.target.value;
    if (isPending || !blocks || error) return [];

    const filteredBlocks = blocks.filter(
      (c: any) =>
        c.key.toLowerCase().includes(searchText.toLowerCase()) ||
        c.description.toLowerCase().includes(searchText.toLowerCase())
    );

    setFiltered(filteredBlocks);
    setFocusedIndex(0);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (!filtered.length) return;
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setFocusedIndex((prev) => (prev + 1) % filtered.length);
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setFocusedIndex((prev) => (prev - 1 + filtered.length) % filtered.length);
    } else if (e.key === "Enter" && filtered[focusedIndex]) {
      handleSelect(filtered[focusedIndex]);
    }
  };

  const handleSelect = (item: any) => {
    onSelectNewNode(item, "baseNode");
    close();
  };

  useEffect(() => {
    if (itemRefs.current[focusedIndex]) {
      itemRefs.current[focusedIndex]?.scrollIntoView({
        behavior: "smooth",
        block: "nearest",
      });
    }
  }, [focusedIndex]);

  useEffect(() => {
    if (blocks) {
      setFiltered(blocks);
    }
  }, [blocks]);

  return (
    <Modal opened={opened} onClose={close} withCloseButton={false} size="xl">
      <TextInput
        px="md"
        placeholder="Search blocks"
        size="lg"
        variant="subtle"
        leftSection={<IconSearch size={20} stroke={1.5} />}
        onChange={handleSearch}
        onKeyDown={handleKeyDown}
        autoFocus
        styles={{ section: { pointerEvents: "none" } }}
        mb="sm"
      />
      <ScrollArea mah="80vh" scrollbarSize={5}>
        {filtered.length > 0 &&
          filtered.map((item: any, index: number) => (
            <Group
              key={item?.key}
              ref={(el) => (itemRefs.current[index] = el)}
              wrap="nowrap"
              p="sm"
              className={`${classes.block} ${
                index === focusedIndex ? classes.focused : ""
              }`}
              tabIndex={0}
              onClick={() => handleSelect(item)}
            >
              <ThemeIcon variant="default" size={40} c="dimmed">
                <IconCube3dSphere size={18} />
              </ThemeIcon>
              <Box w="70%">
                <Text size="sm">{item?.label}</Text>
                <Text opacity={0.7} size="xs" lineClamp={1}>
                  {item?.description}
                </Text>
              </Box>
            </Group>
          ))}
      </ScrollArea>
    </Modal>
  );
};
