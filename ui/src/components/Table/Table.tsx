import React, { useState } from "react";
import { Box, Divider, Stack } from "@mantine/core";
import { DataTable, DataTableProps } from "mantine-datatable";
import "mantine-datatable/styles.layer.css";
import { useDebouncedState } from "@mantine/hooks";
import { IconDatabaseOff } from "@tabler/icons-react";

export type TableProps = {
  showSearch?: boolean;
  searchPlaceholder?: string;
  showRefresh?: boolean;
  showExport?: boolean;
  onRefresh?: () => void;
  onSearchFilterChange?: (value: string) => void;
  menuBar?: React.ReactNode;
} & DataTableProps;

export const UploadPilotDataTable: React.FC<TableProps> = (props) => {
  return (
    <Stack gap={2}>
      {props.menuBar && <Box>{props.menuBar}</Box>}
      <Divider p={0} m={0} mt="2" />
      <DataTable
        backgroundColor="transparent"
        selectionCheckboxProps={{
          style: { "*": { cursor: "pointer" } },
        }}
        noRecordsIcon={
          (
            <IconDatabaseOff
              size={50}
              stroke={1}
              style={{ marginBottom: "10px" }}
            />
          ) as any
        }
        {...props}
      />
    </Stack>
  );
};

export const useUploadPilotDataTable = (searchDelay = 1000) => {
  const [searchFilter, onSearchFilterChange] = useDebouncedState<string>(
    "",
    searchDelay
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
