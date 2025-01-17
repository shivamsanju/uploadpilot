import React, { useState } from 'react';
import { Group, Stack, TextInput } from '@mantine/core';
import { DataTable, DataTableProps } from 'mantine-datatable';
import { IconDownload, IconRefresh, IconSearch } from '@tabler/icons-react';
import classes from "./Table.module.css";
import 'mantine-datatable/styles.layer.css';
import { useDebouncedState, useDisclosure, useTimeout } from '@mantine/hooks';

export type TableProps = {
    showSearch?: boolean
    showRefresh?: boolean
    showExport?: boolean
    onRefresh?: () => void
    onSearchFilterChange?: (value: string) => void
} & DataTableProps


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
    }


    return (
        <Stack>
            <Group justify='flex-end' gap="md">
                {props.showExport && <IconDownload size={18} className={classes.icon} />}
                {props.showRefresh &&
                    <IconRefresh
                        size={18}
                        className={`${classes.icon} ${rotate ? classes.rotate : ''}`}
                        onClick={handleRefresh}
                    />}
                {props.showSearch &&
                    <TextInput
                        leftSection={<IconSearch size={18} />}
                        variant="subtle"
                        placeholder="Search"
                        size="xs"
                        className={classes.search}
                        onChange={(e) => handleSearchChange(e.target.value)}
                    />
                }
            </Group>
            <DataTable
                backgroundColor="transparent"
                {...props}
            />

        </Stack>
    );
}




export const useUploadPilotDataTable = (searchDelay = 500) => {
    const [searchFilter, onSearchFilterChange] = useDebouncedState<string>('', searchDelay);

    const [page, onPageChange] = useState<number>(1);
    const [recordsPerPage, onRecordsPerPageChange] = useState<number>(10);

    return { searchFilter, onSearchFilterChange, page, onPageChange, recordsPerPage, onRecordsPerPageChange }
}