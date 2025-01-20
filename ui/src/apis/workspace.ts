import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const useGetWorkspaces = () => {
    const queryClient = useQueryClient();

    const { isPending, error, data } = useQuery({
        queryKey: ['workspaces'],
        queryFn: () =>
            axiosInstance
                .get(`/workspaces`)
                .then((res) => res.data)
    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    return { isPending, error, workspaces: data, invalidate }
}

export const useCreateWorkspaceMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['workspaces'],
        mutationFn: (name: string) => axiosInstance.post("/workspaces", { name }).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workspaces'] });
            notifications.show({
                title: "Success",
                message: "Workspace created successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to create workspace",
                color: "red",
            });
        },
    })
};


export const useGetAllAllowedSources = (workspaceId: string) => {
    const { isPending, error, data: allowedSources } = useQuery({
        queryKey: ['allowedSources'],
        queryFn: () => {
            if (!workspaceId) {
                return Promise.reject(new Error('workspaceId is required'));
            }
            return axiosInstance
                .get(`workspaces/${workspaceId}/allowedSources`)
                .then((res) => res.data)
        }

    })

    return { isPending, error, allowedSources }
}

export const useGetUsersInWorkspace = (workspaceId: string) => {
    const queryClient = useQueryClient();

    const { isPending, error, data } = useQuery({
        queryKey: ['workspaceUsers', workspaceId],
        queryFn: () => {
            if (!workspaceId) {
                return Promise.reject(new Error('workspaceId is required'));
            }
            return axiosInstance
                .get(`workspaces/${workspaceId}/users`)
                .then((res) => res.data)
        }
    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['workspaceUsers', workspaceId] });
    return { isPending, error, users: data, invalidate }
}