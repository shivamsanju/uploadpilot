import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { APIKey, type CreateApiKeyData } from '../types/apikey';
import axiosInstance from '../utils/axios';

export const useGetApiKeys = () => {
  const queryClient = useQueryClient();
  const {
    isPending,
    error,
    data: apikeys,
  } = useQuery<APIKey[]>({
    queryKey: ['apikeys'],
    queryFn: () => {
      return axiosInstance.get(`/api-keys`).then(res => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ['apikeys'] });

  return { isPending, error, apikeys, invalidate };
};

export const useCreateApiKeyMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['ApiKeys'],
    mutationFn: ({ data }: { data: CreateApiKeyData }) => {
      return axiosInstance.post(`/api-keys`, data).then(res => res.data);
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
    mutationFn: ({ id }: { id: string }) => {
      return axiosInstance.post(`/api-keys/${id}/revoke`).then(res => res.data);
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
