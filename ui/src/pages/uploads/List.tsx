import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  Menu,
  Stack,
  Box,
  Button,
  Group,
  ActionIcon,
  Text,
  TextInput,
} from "@mantine/core";
import {
  IconDots,
  IconBraces,
  IconChevronsDown,
  IconLogs,
  IconDownload,
  IconCopy,
  IconSearch,
  IconBolt,
} from "@tabler/icons-react";
import { useParams } from "react-router-dom";
import { useGetUploads } from "../../apis/upload";
import { timeAgo } from "../../utils/datetime";
import { formatBytes } from "../../utils/utility";
import {
  UploadPilotDataTable,
  useUploadPilotDataTable,
} from "../../components/Table/Table";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { DataTableColumn } from "mantine-datatable";
import { getFileIcon } from "../../utils/fileicons";
import { useViewportSize } from "@mantine/hooks";
import { UploadStatus } from "./Status";
import { LogsModal } from "./Logs";
import { MetadataModal } from "./Metadata";
import { ContainerOverlay } from "../../components/Overlay";
import { RefreshButton } from "../../components/Buttons/RefreshButton/RefreshButton";

const batchSize = 20;

const handleDownload = (url: string) => {
  window.open(url, "_blank");
};

const handleCopyToClipboard = (url: string) => {
  navigator.clipboard.writeText(url);
};

const UploadList = ({ setTotalRecords }: any) => {
  const scrollViewportRef = useRef<HTMLDivElement>(null);
  const { width } = useViewportSize();
  const [openModal, setOpenModal] = useState(false);
  const [modalVariant, setModalVariant] = useState<"logs" | "metadata">("logs");
  const [uploadId, setUploadId] = useState<string | null>(null);
  const [metadata, setMetadata] = useState({});
  const [selectedRecords, setSelectedRecords] = useState<any[]>([]);

  const { workspaceId } = useParams();
  const { searchFilter, onSearchFilterChange } = useUploadPilotDataTable();

  const {
    isPending,
    error,
    isFetchNextPageError,
    uploads,
    fetchNextPage,
    totalRecords,
    isFetchingNextPage,
    invalidate,
    hasNextPage,
  } = useGetUploads({
    workspaceId: workspaceId || "",
    batchSize,
    search: searchFilter,
  });

  const handleViewLogs = useCallback(
    (importId: string) => {
      const uploadItem = uploads?.find((item: any) => item.id === importId);
      if (uploadItem) {
        setOpenModal(true);
        setUploadId(uploadItem.id);
        setModalVariant("logs");
      }
    },
    [uploads]
  );

  const handleViewMetadata = useCallback(
    (importId: string) => {
      const uploadItem = uploads?.find((item: any) => item.id === importId);
      if (uploadItem) {
        setOpenModal(true);
        setMetadata(uploadItem.metadata || {});
        setModalVariant("metadata");
      }
    },
    [uploads]
  );

  const handleRefresh = () => {
    invalidate();
    scrollViewportRef.current?.scrollTo(0, 0);
  };

  const colDefs: DataTableColumn[] = useMemo(() => {
    return [
      {
        title: "Name",
        accessor: "metadata.filename",
        elipsis: true,
        render: (params: any) => (
          <Group align="center" gap="sm">
            {getFileIcon(params?.fileType, 20)}
            <Text fz="sm">{params?.fileName}</Text>
            {/* <Text fz="xs" c="dimmed">
              Filename
            </Text> */}
          </Group>
        ),
      },
      {
        title: "File Type",
        accessor: "metadata.filetype",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">{params?.fileType}</Text>
            {/* <Text fz="xs" c="dimmed">
              Mime Type
            </Text> */}
          </>
        ),
      },
      {
        title: "Size",
        accessor: "size",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">{formatBytes(Number(params?.size))}</Text>
            {/* <Text fz="xs" c="dimmed">
              Size
            </Text> */}
          </>
        ),
      },
      {
        title: "Finished At",
        accessor: "finishedAt",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">
              {params?.finishedAt &&
                timeAgo.format(new Date(params?.finishedAt))}
            </Text>
            {/* <Text fz="xs" c="dimmed">
              Uploaded
            </Text> */}
          </>
        ),
      },
      {
        title: "Status",
        accessor: "status",
        render: (params: any) => (
          <Group align="center" gap="sm">
            <UploadStatus status={params?.status} />
            {params?.status}
          </Group>
        ),
      },
      {
        title: "Actions",
        accessor: "actions",
        textAlign: "right",
        width: 100,
        render: (params: any) => (
          <Group gap={0} justify="flex-end">
            <Menu
              transitionProps={{ transition: "pop" }}
              withArrow
              position="bottom-end"
              withinPortal
            >
              <Menu.Target>
                <ActionIcon variant="subtle" color="dimmed">
                  <IconDots size={16} stroke={1.5} />
                </ActionIcon>
              </Menu.Target>
              <Menu.Dropdown>
                {params?.url && (
                  <Menu.Item
                    onClick={() => handleDownload(params?.url)}
                    leftSection={<IconDownload size={18} />}
                  >
                    <Text>Download</Text>
                  </Menu.Item>
                )}
                {params?.url && (
                  <Menu.Item
                    onClick={() => handleCopyToClipboard(params?.url)}
                    leftSection={<IconCopy size={18} />}
                  >
                    <Text>Copy URL</Text>
                  </Menu.Item>
                )}
                <Menu.Item
                  onClick={() => handleViewLogs(params?.id)}
                  leftSection={<IconLogs size={18} />}
                >
                  <Text>View Logs</Text>
                </Menu.Item>
                <Menu.Item
                  onClick={() => handleViewMetadata(params?.id)}
                  leftSection={<IconBraces size={18} />}
                >
                  <Text>View Metadata</Text>
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ];
  }, [handleViewLogs, handleViewMetadata, width]);

  useEffect(() => {
    setTotalRecords(totalRecords);
  }, [setTotalRecords, totalRecords]);

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  if (isFetchNextPageError) {
    return (
      <ErrorCard title="Error" message={"Failed to fetch next page"} h="70vh" />
    );
  }

  return (
    <Box mt="lg">
      {/* Loading overlay only while pending, not on refetch*/}
      <ContainerOverlay visible={isPending} />{" "}
      <LogsModal
        open={openModal && modalVariant === "logs"}
        onClose={() => setOpenModal(false)}
        uploadId={uploadId || ""}
        workspaceId={workspaceId || ""}
      />
      <MetadataModal
        open={openModal && modalVariant === "metadata"}
        onClose={() => setOpenModal(false)}
        metadata={metadata || {}}
      />
      <Box mr="md">
        <UploadPilotDataTable
          minHeight={500}
          verticalSpacing="xs"
          horizontalSpacing="lg"
          noHeader={false}
          showSearch={true}
          searchPlaceholder='Search imports by name or status. For metadata search use {key: "regex"}'
          showExport={true}
          showRefresh={true}
          onRefresh={handleRefresh}
          onSearchFilterChange={onSearchFilterChange}
          onScrollToBottom={fetchNextPage}
          columns={colDefs}
          records={uploads}
          scrollViewportRef={scrollViewportRef}
          noRecordsText="No imports yet"
          selectedRecords={selectedRecords}
          onSelectedRecordsChange={setSelectedRecords}
          menuBar={
            <Group gap="sm" align="center" justify="space-between">
              <Group gap="sm">
                <Button variant="subtle" leftSection={<IconBolt size={18} />}>
                  Process
                </Button>
                <RefreshButton variant="subtle" onClick={handleRefresh} />
              </Group>
              <TextInput
                value={searchFilter}
                onChange={(e) => onSearchFilterChange(e.target.value)}
                placeholder="Search by name or status"
                leftSection={<IconSearch size={18} />}
                variant="subtle"
              />
            </Group>
          }
        />
      </Box>
      <Stack align="center" justify="center" p="md" key="selectedRecords">
        <Button
          display={hasNextPage ? "block" : "none"}
          leftSection={<IconChevronsDown size={16} />}
          variant="subtle"
          disabled={isPending || isFetchingNextPage || !hasNextPage}
          loading={isFetchingNextPage}
          onClick={() => fetchNextPage({})}
        >
          Load More
        </Button>
      </Stack>
    </Box>
  );
};

export default UploadList;
