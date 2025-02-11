import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const useGetProcessors = (workspaceId: string) => {
  const queryClient = useQueryClient();
  const {
    isPending,
    error,
    isFetching,
    data: processors,
  } = useQuery({
    queryKey: ["processors", workspaceId],
    queryFn: () => {
      if (!workspaceId) {
        return Promise.reject(new Error("workspaceId is required"));
      }
      return axiosInstance
        .get(`/workspaces/${workspaceId}/processors`)
        .then((res) => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ["processors", workspaceId] });

  return { isPending, error, processors, invalidate, isFetching };
};

export const useGetProcessor = (workspaceId: string, processorId: string) => {
  const queryClient = useQueryClient();
  const {
    isPending,
    error,
    data: processor,
  } = useQuery({
    queryKey: ["processorDetails", workspaceId, processorId],
    queryFn: () => {
      if (!workspaceId || !processorId) {
        return Promise.reject(
          new Error("workspaceId and processorId are required")
        );
      }
      return axiosInstance
        .get(`/workspaces/${workspaceId}/processors/${processorId}`)
        .then((res) => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({
      queryKey: ["processorDetails", workspaceId, processorId],
    });

  return { isPending, error, processor, invalidate };
};

export const useCreateProcessorMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["processors"],
    mutationFn: ({
      workspaceId,
      processor,
    }: {
      workspaceId: string;
      processor: any;
    }) => {
      return axiosInstance
        .post(`/workspaces/${workspaceId}/processors`, processor)
        .then((res) => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["processors"] });
      notifications.show({
        title: "Success",
        message: "Processor created successfully",
        color: "green",
      });
    },
    onError: () => {
      notifications.show({
        title: "Error",
        message: "Failed to create Processor",
        color: "red",
      });
    },
  });
};

export const useUpdateProcessorMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["processors"],
    mutationFn: ({
      workspaceId,
      processorId,
      processor,
    }: {
      workspaceId: string;
      processorId: string;
      processor: any;
    }) => {
      return axiosInstance
        .put(`/workspaces/${workspaceId}/processors/${processorId}`, processor)
        .then((res) => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["processors"] });
      notifications.show({
        title: "Success",
        message: "Processor created successfully",
        color: "green",
      });
    },
    onError: () => {
      notifications.show({
        title: "Error",
        message: "Failed to create Processor",
        color: "red",
      });
    },
  });
};

export const useDeleteProcessorMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["processors"],
    mutationFn: ({
      workspaceId,
      processorId,
    }: {
      workspaceId: string;
      processorId: string;
    }) => {
      return axiosInstance
        .delete(`/workspaces/${workspaceId}/processors/${processorId}`)
        .then((res) => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["processors"] });
      notifications.show({
        title: "Success",
        message: "Processor deleted successfully",
        color: "green",
      });
    },
    onError: () => {
      notifications.show({
        title: "Error",
        message: "Failed to delete Processor",
        color: "red",
      });
    },
  });
};

export const useEnableDisableProcessorMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["processors"],
    mutationFn: ({
      workspaceId,
      processorId,
      enabled,
    }: {
      workspaceId: string;
      processorId: string;
      enabled: boolean;
    }) => {
      return axiosInstance
        .put(
          `/workspaces/${workspaceId}/processors/${processorId}/${
            enabled ? "enable" : "disable"
          }`
        )
        .then((res) => res.data);
    },
    onSuccess: (_, { enabled }) => {
      queryClient.invalidateQueries({ queryKey: ["processors"] });
      notifications.show({
        title: "Success",
        message: `Processors ${enabled ? "enabled" : "disabled"} successfully`,
        color: "green",
      });
    },
    onError: (error: any, { enabled }) => {
      notifications.show({
        title: "Error",
        message: `Failed to ${
          enabled ? "enable" : "disable"
        } Processors. Reason: ${
          error?.response?.data?.message || error.message
        }`,
        color: "red",
      });
    },
  });
};

// Editor related

export const useGetAllProcBlocks = (workspaceId: string) => {
  const {
    isPending,
    error,
    data: blocks,
  } = useQuery({
    queryKey: ["procblocks"],
    queryFn: () => {
      if (!workspaceId) {
        return Promise.reject(new Error("workspaceId is required"));
      }
      return axiosInstance
        .get(`/workspaces/${workspaceId}/processors/blocks`)
        .then((res) => res.data);
    },
  });

  return { isPending, error, blocks };
};

export const useUpdateProcessorWorkflowMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["processorDetails"],
    mutationFn: ({
      workspaceId,
      processorId,
      workflow,
    }: {
      workspaceId: string;
      processorId: string;
      workflow: string;
    }) => {
      return axiosInstance
        .put(`/workspaces/${workspaceId}/processors/${processorId}/workflow`, {
          workflow,
        })
        .then((res) => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["processorDetails"] });
      notifications.show({
        title: "Success",
        message: "Workflow updated successfully",
        color: "green",
      });
    },
    onError: (error: any) => {
      notifications.show({
        title: "Error",
        message: `Failed to update workflow Reason: ${
          error?.response?.data?.message || error.message
        }`,
        color: "red",
      });
    },
  });
};
