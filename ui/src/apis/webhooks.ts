import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Webhook } from "../types/webhook";


export const useGetWebhooks = (workspaceId: string) => {
    const queryClient = useQueryClient();
    const { isPending, error, data: webhooks } = useQuery<Webhook[]>({
        queryKey: ['webhooks', workspaceId],
        queryFn: () => {
            if (!workspaceId) {
                return Promise.reject(new Error('workspaceId is required'));
            }
            return axiosInstance
                .get(`/workspaces/${workspaceId}/webhooks`)
                .then((res) => res.data)
        }

    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['webhooks', workspaceId] });

    return { isPending, error, webhooks, invalidate }
}

export const useGetWebhook = (workspaceId: string, webhookId: string) => {
    const queryClient = useQueryClient();
    const { isPending, error, data: webhook } = useQuery<Webhook>({
        queryKey: ['webhook', workspaceId, webhookId],
        queryFn: () => {
            if (!workspaceId || !webhookId) {
                return Promise.reject(new Error('workspaceId and webhookId are required'));
            }
            return axiosInstance
                .get(`/workspaces/${workspaceId}/webhooks`)
                .then((res) => res.data)
        }

    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['webhook', workspaceId, webhookId] });

    return { isPending, error, webhook, invalidate }
}

export const useCreateWebhookMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['webhooks'],
        mutationFn: ({ workspaceId, webhook }: { workspaceId: string, webhook: Webhook }) => {
            return axiosInstance.post(`/workspaces/${workspaceId}/webhooks`, webhook).then((res) => res.data)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['webhooks'] });
            notifications.show({
                title: "Success",
                message: "Webhook created successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to create webhook",
                color: "red",
            });
        },
    })
};

export const useDeleteWebhookMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['webhooks'],
        mutationFn: ({ workspaceId, webhookId }: { workspaceId: string, webhookId: string }) => {
            return axiosInstance.delete(`/workspaces/${workspaceId}/webhooks/${webhookId}`).then((res) => res.data)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['webhooks'] });
            notifications.show({
                title: "Success",
                message: "Webhook deleted successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to delete webhook",
                color: "red",
            });
        },
    })
};

export const useUpdateWebhookMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['webhooks'],
        mutationFn: ({ workspaceId, webhookId, webhook }: { workspaceId: string, webhookId: string, webhook: Webhook }) => {
            return axiosInstance.put(`/workspaces/${workspaceId}/webhooks/${webhookId}`, webhook).then((res) => res.data)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['webhooks'] });
            notifications.show({
                title: "Success",
                message: "Webhook updated successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to update webhook",
                color: "red",
            });
        },
    })
}

export const useEnableDisableWebhookMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['webhooks'],
        mutationFn: ({ workspaceId, webhookId, enabled }: { workspaceId: string, webhookId: string, enabled: boolean }) => {
            return axiosInstance.patch(`/workspaces/${workspaceId}/webhooks/${webhookId}`, { enabled }).then((res) => res.data)
        },
        onSuccess: (_, { enabled }) => {
            queryClient.invalidateQueries({ queryKey: ['webhooks'] });
            notifications.show({
                title: "Success",
                message: `Webhook ${enabled ? 'enabled' : 'disabled'} successfully`,
                color: "green",
            });
        },
        onError: (error: any, { enabled }) => {
            notifications.show({
                title: "Error",
                message: `Failed to ${enabled ? 'enable' : 'disable'} webhook. Reason: ${error?.response?.data?.message || error.message}`,
                color: "red",
            });
        },
    })
}