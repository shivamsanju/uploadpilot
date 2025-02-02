import React, { useState } from "react";
import { Box, Button, Divider, Group, Stack, TextInput } from "@mantine/core";
import { DataTable, DataTableProps } from "mantine-datatable";
import { IconDownload, IconRefresh, IconSearch } from "@tabler/icons-react";
import classes from "./Table.module.css";
import "mantine-datatable/styles.layer.css";
import { useDebouncedState, useDisclosure, useTimeout } from "@mantine/hooks";

export type TableProps = {
  showSearch?: boolean;
  searchPlaceholder?: string;
  showRefresh?: boolean;
  showExport?: boolean;
  onRefresh?: () => void;
  onSearchFilterChange?: (value: string) => void;
} & DataTableProps;

export const UploadPilotDataTable: React.FC<TableProps> = (props) => {
  const [rotate, handlers] = useDisclosure(false);
  const { start } = useTimeout(() => {
    handlers.close();
  }, 1000);

  const handleSearchChange = (value: string) => {
    props.onSearchFilterChange && props.onSearchFilterChange(value);
  };

  const handleRefresh = () => {
    if (rotate) return;
    handlers.open();
    props.onRefresh && props.onRefresh();
    start();
  };

  return (
    <Stack gap={2}>
      <Group justify="space-between" p={0} m={0}>
        <Box w="70%">
          {props.showSearch && (
            <TextInput
              leftSection={<IconSearch size={18} />}
              variant="subtle"
              placeholder={props.searchPlaceholder || "Search"}
              className={classes.search}
              onChange={(e) => handleSearchChange(e.target.value)}
            />
          )}
        </Box>
        <Group gap="md" justify="flex-end">
          {props.showExport && (
            <Button
              variant="subtle"
              className={classes.tableExtraBtn}
              leftSection={<IconDownload size={18} />}
            >
              Export
            </Button>
          )}
          {props.showRefresh && (
            <Button
              className={classes.tableExtraBtn}
              onClick={handleRefresh}
              variant="subtle"
              leftSection={
                <IconRefresh
                  size={15}
                  className={`${rotate ? classes.rotate : ""}`}
                />
              }
            >
              Refresh
            </Button>
          )}
        </Group>
      </Group>
      {(props.showSearch || props.showRefresh || props.showExport) && (
        <Divider p={0} m={0} />
      )}
      <DataTable backgroundColor="transparent" {...props} />
    </Stack>
  );
};

export const useUploadPilotDataTable = (searchDelay = 1000) => {
  const [searchFilter, onSearchFilterChange] = useDebouncedState<string>(
    "",
    searchDelay,
  );

  const [page, onPageChange] = useState<number>(1);
  const [recordsPerPage, onRecordsPerPageChange] = useState<number>(10);

  return {
    searchFilter,
    onSearchFilterChange,
    page,
    onPageChange,
    recordsPerPage,
    onRecordsPerPageChange,
  };
};
