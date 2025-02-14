import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  Menu,
  Stack,
  Box,
  Button,
  Group,
  ActionIcon,
  Text,
  Avatar,
} from "@mantine/core";
import {
  IconDots,
  IconBraces,
  IconChevronsDown,
  IconLogs,
  IconDownload,
  IconCopy,
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
        title: "",
        accessor: "id",
        hidden: width < 768,
        render: (params: any) => (
          <Avatar size={40} radius={40} variant="light">
            {getFileIcon(params?.metadata?.filetype, 20)}
          </Avatar>
        ),
      },
      {
        title: "Name",
        accessor: "metadata.filename",
        elipsis: true,
        render: (params: any) => (
          <>
            <Text fz="sm">{params?.metadata?.filename}</Text>
            <Text fz="xs" c="dimmed">
              Filename
            </Text>
          </>
        ),
      },
      {
        title: "File Type",
        accessor: "metadata.filetype",
        textAlign: "center",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">{params?.metadata?.filetype}</Text>
            <Text fz="xs" c="dimmed">
              Mime Type
            </Text>
          </>
        ),
      },
      {
        title: "Size",
        accessor: "size",
        textAlign: "center",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">{formatBytes(Number(params?.size))}</Text>
            <Text fz="xs" c="dimmed">
              Size
            </Text>
          </>
        ),
      },
      {
        title: "Finished At",
        accessor: "finishedAt",
        textAlign: "center",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">
              {params?.finishedAt &&
                timeAgo.format(new Date(params?.finishedAt))}
            </Text>
            <Text fz="xs" c="dimmed">
              Uploaded
            </Text>
          </>
        ),
      },
      {
        title: "Status",
        accessor: "status",
        textAlign: "center",
        render: (params: any) => (
          <>
            <UploadStatus status={params?.status} />
            <Text fz="xs" c="dimmed">
              {params?.status}
            </Text>
          </>
        ),
      },
      {
        title: "Actions",
        accessor: "actions",
        textAlign: "center",
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
          verticalSpacing="lg"
          horizontalSpacing="lg"
          fetching={isFetchingNextPage}
          noHeader={true}
          showSearch={true}
          searchPlaceholder='Search imports by name or status. For metadata search use {key: "regex"}'
          showExport={true}
          showRefresh={true}
          onRefresh={handleRefresh}
          onSearchFilterChange={onSearchFilterChange}
          columns={colDefs}
          records={uploads}
          selectionCheckboxProps={{ style: { cursor: "pointer" } }}
          onScrollToBottom={fetchNextPage}
          scrollViewportRef={scrollViewportRef}
          noRecordsText="No imports yet"
        />

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
    </Box>
  );
};

export default UploadList;
