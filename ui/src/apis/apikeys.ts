import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { APIKey, type CreateApiKeyData } from '../types/apikey';
import axiosInstance from '../utils/axios';

export const useGetApiKeysInWorkspace = (workspaceId: string) => {
  const queryClient = useQueryClient();
  const {
    isPending,
    error,
    data: apikeys,
  } = useQuery<APIKey[]>({
    queryKey: ['apikeys', workspaceId],
    queryFn: () => {
      if (!workspaceId) {
        return Promise.reject(new Error('workspaceId is required'));
      }
      return axiosInstance
        .get(`/workspaces/${workspaceId}/apikeys`)
        .then(res => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ['apikeys', workspaceId] });

  return { isPending, error, apikeys, invalidate };
};

export const useCreateApiKeyMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['ApiKeys'],
    mutationFn: ({
      workspaceId,
      data,
    }: {
      workspaceId: string;
      data: CreateApiKeyData;
    }) => {
      return axiosInstance
        .post(`/workspaces/${workspaceId}/apikeys`, data)
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apikeys'] });
      notifications.show({
        title: 'Success',
        message: 'APIKey created successfully',
        color: 'green',
      });
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to create APIKey',
        color: 'red',
      });
    },
  });
};

export const useRevokeApiKeyMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['apikeys'],
    mutationFn: ({ workspaceId, id }: { workspaceId: string; id: string }) => {
      return axiosInstance
        .post(`/workspaces/${workspaceId}/apikeys/${id}/revoke`)
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apikeys'] });
      notifications.show({
        title: 'Success',
        message: 'APIKey revoked successfully',
        color: 'green',
      });
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to revoke APIKey',
        color: 'red',
      });
    },
  });
};
