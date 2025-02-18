import {
  IconDots,
  IconRouteOff,
  IconEdit,
  IconClockBolt,
  IconLogs,
} from "@tabler/icons-react";
import {
  ActionIcon,
  Box,
  Group,
  Menu,
  Modal,
  Text,
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

const ProcessorList = ({
  opened,
  setOpened,
}: {
  opened: boolean;
  setOpened: any;
}) => {
  const [workflowId, setWorkflowId] = useState<string>("");
  const [runId, setRunId] = useState<string>("");

  const { width } = useViewportSize();
  const { workspaceId, processorId } = useParams();
  const { isPending, error, runs } = useGetProcessorRuns(
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
        accessor: "runId",
        title: "Run ID",
        render: (item: any) => (
          <Group gap="xs" align="center">
            <IconClockBolt size={16} stroke={1.5} />
            {item?.runId}
          </Group>
        ),
      },
      {
        accessor: "status",
        title: "Status",
        textAlign: "center",
      },
      {
        title: "Started At",
        accessor: "startTime",
        textAlign: "center",
        hidden: width < 768,
        render: (item: any) =>
          new Date(item?.startTime).toLocaleString("en-US"),
      },
      {
        title: "Ended At",
        accessor: "endTime",
        textAlign: "center",
        hidden: width < 768,
        render: (item: any) => new Date(item?.endTime).toLocaleString("en-US"),
      },
      {
        title: "Execution Time",
        accessor: "durationSeconds",
        textAlign: "center",
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
        minHeight={700}
        columns={columns}
        records={runs}
        verticalSpacing="sm"
        horizontalSpacing="md"
        noHeader={false}
        noRecordsText="No runs yet"
        noRecordsIcon={<IconRouteOff size={100} />}
        highlightOnHover
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

export default ProcessorList;
