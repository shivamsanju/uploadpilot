import {
  IconDots,
  IconEdit,
  IconLogs,
  IconBolt,
  IconSearch,
} from "@tabler/icons-react";
import {
  ActionIcon,
  Badge,
  Box,
  Button,
  Group,
  Menu,
  Modal,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useParams } from "react-router-dom";
import { useCallback, useMemo, useState } from "react";
import { DataTableColumn } from "mantine-datatable";
import { UploadPilotDataTable } from "../../components/Table/Table";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { useGetProcessorRuns } from "../../apis/processors";
import { useViewportSize } from "@mantine/hooks";
import { ContainerOverlay } from "../../components/Overlay";
import { LogsModal } from "./Logs";
import { RefreshButton } from "../../components/Buttons/RefreshButton/RefreshButton";
import { statusConfig } from "./status";

const ProcessorRunsList = () => {
  const [workflowId, setWorkflowId] = useState<string>("");
  const [runId, setRunId] = useState<string>("");
  const [opened, setOpened] = useState(false);
  const [selectedRecords, setSelectedRecords] = useState<any[]>([]);

  const { width } = useViewportSize();
  const { workspaceId, processorId } = useParams();
  const { isPending, error, runs, invalidate } = useGetProcessorRuns(
    workspaceId || "",
    processorId || ""
  );

  const handleViewLogs = useCallback((runId: string, workflowId: string) => {
    setRunId(runId);
    setWorkflowId(workflowId);
    setOpened(true);
  }, []);

  const closeLogs = useCallback(() => {
    setRunId("");
    setWorkflowId("");
    setOpened(false);
  }, []);

  const columns: DataTableColumn[] = useMemo(
    () => [
      {
        accessor: "id",
        title: "",
        width: 20,
        render: (item: any) => (
          <Stack justify="center">
            <IconBolt size={16} stroke={1.5} />
          </Stack>
        ),
      },
      {
        accessor: "runId",
        title: "Run ID",
      },
      {
        accessor: "status",
        title: "Status",
        render: (item: any) => (
          <Badge
            variant="light"
            color={statusConfig[item?.status?.toLowerCase() || ""]}
            size="sm"
          >
            {item?.status}
          </Badge>
        ),
      },
      {
        title: "Started At",
        accessor: "startTime",
        hidden: width < 768,
        render: (item: any) =>
          new Date(item?.startTime).toLocaleString("en-US"),
      },
      {
        title: "Ended At",
        accessor: "endTime",
        hidden: width < 768,
        render: (item: any) => {
          if (item?.endTime !== "0001-01-01T00:00:00Z") {
            return new Date(item?.endTime).toLocaleString("en-US");
          }
          return "-";
        },
      },
      {
        title: "Execution Time",
        accessor: "durationSeconds",
        hidden: width < 768,
        render: (item: any) => `${item?.durationSeconds} s`,
      },
      {
        accessor: "actions",
        title: "Actions",
        textAlign: "right",
        render: (item: any) => (
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
                <Menu.Item
                  leftSection={<IconLogs size={16} stroke={1.5} />}
                  onClick={() => handleViewLogs(item?.runId, item?.workflowId)}
                >
                  <Text>View Logs</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconEdit size={16} stroke={1.5} />}
                  disabled={true}
                >
                  <Text>Trigger Again</Text>
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ],
    [width, handleViewLogs]
  );

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  return (
    <Box mr="md">
      <ContainerOverlay visible={isPending} />
      <UploadPilotDataTable
        minHeight={500}
        columns={columns}
        records={runs}
        verticalSpacing="xs"
        horizontalSpacing="md"
        noHeader={false}
        noRecordsText="No runs yet"
        highlightOnHover
        menuBar={
          <Group gap="sm" align="center" justify="space-between">
            <Group gap="sm">
              <Button
                variant="subtle"
                onClick={() => setOpened(true)}
                leftSection={<IconBolt size={18} />}
              >
                Re Trigger
              </Button>

              <RefreshButton onClick={invalidate} />
            </Group>
            <TextInput
              value={""}
              // onChange={(e) => onSearchFilterChange(e.target.value)}
              placeholder="Search by name or status"
              leftSection={<IconSearch size={18} />}
              variant="subtle"
            />
          </Group>
        }
        selectedRecords={selectedRecords}
        onSelectedRecordsChange={setSelectedRecords}
      />
      <Modal
        padding="xl"
        transitionProps={{ transition: "pop" }}
        opened={opened}
        onClose={() => {
          setOpened(false);
        }}
        title={
          <Title order={5} opacity={0.7}>
            Logs
          </Title>
        }
        closeOnClickOutside={false}
        size="lg"
      >
        {workspaceId && processorId && runId && workflowId && (
          <div>
            <LogsModal
              open={opened}
              onClose={closeLogs}
              workspaceId={workspaceId}
              processorId={processorId}
              workflowId={workflowId}
              runId={runId}
            />
          </div>
        )}
      </Modal>
    </Box>
  );
};

export default ProcessorRunsList;
