import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { UploaderConfig } from "../types/uploader";


export const useGetUploaderConfig = (workspaceId: string) => {
    const queryClient = useQueryClient();
    const { isPending, error, data: config } = useQuery({
        queryKey: ['workspace.uploaderConfig', workspaceId],
        queryFn: () => {
            if (!workspaceId) {
                return Promise.reject(new Error('workspaceId is required'));
            }
            return axiosInstance
                .get(`/workspaces/${workspaceId}/config`)
                .then((res) => res.data)
        }

    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['workspace.uploaderConfig', workspaceId] });

    return { isPending, error, config, invalidate }
}

export const useUpdateUploaderConfigMutation = (workspaceId: string) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['workspace.uploaderConfig', workspaceId],
        mutationFn: (config: UploaderConfig) => axiosInstance.put(`/workspaces/${workspaceId}/config`, config).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workspace.uploaderConfig', workspaceId] });
            notifications.show({
                title: "Success",
                message: "Uploader configuration updated successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to update uploader configuration",
                color: "red",
            });
        },
    })
};
