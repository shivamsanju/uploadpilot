import {
  IconCircleCheck,
  IconCircleOff,
  IconDots,
  IconEye,
  IconTrash,
  IconRoute,
  IconRouteOff,
  IconEdit,
  IconChevronRightPipe,
} from "@tabler/icons-react";
import {
  ActionIcon,
  Anchor,
  Avatar,
  Badge,
  Box,
  Group,
  LoadingOverlay,
  Menu,
  Modal,
  Pill,
  Text,
  Title,
} from "@mantine/core";
import { useNavigate, useParams } from "react-router-dom";
import { useCallback, useMemo, useState } from "react";
import { DataTableColumn } from "mantine-datatable";
import { UploadPilotDataTable } from "../../components/Table/Table";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";
import { showNotification } from "@mantine/notifications";
import {
  useDeleteProcessorMutation,
  useEnableDisableProcessorMutation,
  useGetProcessors,
} from "../../apis/processors";
import AddProcessorForm from "./Add";
import { timeAgo } from "../../utils/datetime";
import { useViewportSize } from "@mantine/hooks";

const ProcessorList = ({
  opened,
  setOpened,
}: {
  opened: boolean;
  setOpened: any;
}) => {
  const [mode, setMode] = useState<"add" | "edit" | "view">("add");
  const { width } = useViewportSize();
  const [initialValues, setInitialValues] = useState(null);
  const { workspaceId } = useParams();
  const { isPending, error, processors, isFetching } = useGetProcessors(
    workspaceId || ""
  );
  const { mutateAsync, isPending: isDeleting } = useDeleteProcessorMutation();
  const { mutateAsync: enableDisableProcessor, isPending: isEnabling } =
    useEnableDisableProcessorMutation();
  const navigate = useNavigate();

  const handleRemoveProcessor = useCallback(
    async (processorId: string) => {
      if (
        !workspaceId ||
        workspaceId === "" ||
        !processorId ||
        processorId === ""
      ) {
        showNotification({
          color: "red",
          title: "Error",
          message: "Workspace ID or Processor ID is not available",
        });
        return;
      }

      try {
        await mutateAsync({ workspaceId, processorId });
      } catch (error) {
        console.error(error);
      }
    },
    [workspaceId, mutateAsync]
  );

  const handleEnableDisableProcessor = useCallback(
    async (processorId: string, enabled: boolean) => {
      if (
        !workspaceId ||
        workspaceId === "" ||
        !processorId ||
        processorId === ""
      ) {
        showNotification({
          color: "red",
          title: "Error",
          message: "Workspace ID or Processor ID is not available",
        });
        return;
      }

      try {
        await enableDisableProcessor({ workspaceId, processorId, enabled });
      } catch (error) {
        console.error(error);
      }
    },
    [workspaceId, enableDisableProcessor]
  );

  const handleViewEdit = useCallback(
    async (values: any, mode: "view" | "edit") => {
      setInitialValues(values);
      setMode(mode);
      setOpened(true);
    },
    [setOpened, setMode, setInitialValues]
  );

  const columns: DataTableColumn[] = useMemo(
    () => [
      {
        accessor: "name",
        title: "Name",
        render: (item: any) => (
          <Group gap="sm">
            <Avatar
              size={40}
              radius={40}
              variant="light"
              color={item?.enabled ? "appcolor" : "gray"}
            >
              {item?.enabled ? <IconRoute /> : <IconRouteOff />}
            </Avatar>
            <div>
              <Anchor
                onClick={() =>
                  navigate(`/workspaces/${workspaceId}/processors/${item?.id}`)
                }
                fz="sm"
                fw={500}
              >
                {item.name}
              </Anchor>
              <Text c="dimmed" fz="xs">
                Name
              </Text>
            </div>
          </Group>
        ),
      },
      {
        accessor: "triggers",
        title: "Triggers",
        hidden: width < 768,
        render: (item: any) => (
          <div>
            <Text fz="sm" fw={500}>
              {item.triggers && item.triggers.length > 0 ? (
                <>
                  {item.triggers
                    .slice(0, 5)
                    .map((trigger: any, index: number) => (
                      <Pill mr="xs" size="xs" key={index} variant="light">
                        {trigger}
                      </Pill>
                    ))}
                  {item.triggers.length > 5 && (
                    <Pill size="xs" color="blue">
                      +{item.triggers.length - 5}
                    </Pill>
                  )}
                </>
              ) : (
                "No Triggers"
              )}
            </Text>
            <Text c="dimmed" fz="xs">
              Triggers
            </Text>
          </div>
        ),
      },
      {
        title: "updated At",
        accessor: "updatedAt",
        hidden: width < 768,
        render: (params: any) => (
          <>
            <Text fz="sm">
              {params?.updatedAt && timeAgo.format(new Date(params?.updatedAt))}
            </Text>
            <Text fz="xs" c="dimmed">
              Last Updated
            </Text>
          </>
        ),
      },
      {
        accessor: "enabled",
        title: "Status",
        hidden: width < 768,
        render: (item: any) => (
          <>
            <Badge color={item?.enabled ? "green" : "red"} size="sm">
              {item?.enabled ? "Enabled" : "Disabled"}
            </Badge>
          </>
        ),
      },
      // {
      //     accessor: 'goto',
      //     title: 'goto',
      //     textAlign: 'right',
      //     hidden: width < 768,
      //     render: (item: any) => (
      //         <ActionIcon
      //             variant="light"
      //             size="lg"
      //             onClick={() => navigate(`/workspaces/${workspaceId}/processors/${item?.id}`)}
      //         >
      //             <IconChevronRight />
      //         </ActionIcon>
      //     ),
      // },
      {
        accessor: "actions",
        title: "Actions",
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
                  leftSection={<IconEye size={16} stroke={1.5} />}
                  disabled={!item?.enabled}
                  onClick={() => handleViewEdit(item, "view")}
                >
                  <Text>View</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconEdit size={16} stroke={1.5} />}
                  disabled={!item?.enabled}
                  onClick={() => handleViewEdit(item, "edit")}
                >
                  <Text>Edit</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconChevronRightPipe size={16} stroke={1.5} />}
                  disabled={!item?.enabled}
                  onClick={() =>
                    navigate(
                      `/workspaces/${workspaceId}/processors/${item?.id}`
                    )
                  }
                >
                  <Text>Canvas</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={
                    item?.enabled ? (
                      <IconCircleOff size={16} stroke={1.5} />
                    ) : (
                      <IconCircleCheck size={16} stroke={1.5} />
                    )
                  }
                  color={item?.enabled ? "red" : "green"}
                  onClick={() =>
                    handleEnableDisableProcessor(item.id, !item?.enabled)
                  }
                >
                  <Text>{item?.enabled ? "Disable" : "Enable"}</Text>
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconTrash size={16} stroke={1.5} />}
                  color="red"
                  onClick={() => handleRemoveProcessor(item.id)}
                >
                  <Text>Delete</Text>
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        ),
      },
    ],
    [
      handleRemoveProcessor,
      handleEnableDisableProcessor,
      handleViewEdit,
      width,
      navigate,
      workspaceId,
    ]
  );

  if (error) {
    return <ErrorCard title="Error" message={error.message} h="70vh" />;
  }

  return (
    <Box mr="md">
      <LoadingOverlay
        visible={isDeleting || isEnabling || isPending || isFetching}
        overlayProps={{ radius: "sm", blur: 2 }}
      />
      <UploadPilotDataTable
        minHeight={700}
        showSearch={false}
        columns={columns}
        records={processors}
        verticalSpacing="md"
        horizontalSpacing="md"
        noHeader={true}
        noRecordsText="No processors. Create a processor by clicking the plus icon on top right to get started."
        noRecordsIcon={<IconRouteOff size={100} />}
        highlightOnHover
      />
      <Modal
        padding="xl"
        transitionProps={{ transition: "pop" }}
        opened={opened}
        onClose={() => {
          setOpened(false);
          setInitialValues(null);
          setMode("add");
        }}
        title={
          <Title order={5} opacity={0.7}>
            {mode === "edit"
              ? "Edit Processor"
              : mode === "view"
              ? "View Details"
              : "Create Processor"}
          </Title>
        }
        closeOnClickOutside={false}
        size="lg"
      >
        <AddProcessorForm
          mode={mode}
          setOpened={setOpened}
          workspaceId={workspaceId || ""}
          initialValues={initialValues}
          setInitialValues={setInitialValues}
          setMode={setMode}
        />
      </Modal>
    </Box>
  );
};

export default ProcessorList;
