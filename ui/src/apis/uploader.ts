import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { WorkspaceConfig } from '../types/uploader';
import { axiosTenantInstance } from '../utils/axios';

export const useGetUploaderConfig = (workspaceId: string) => {
  const queryClient = useQueryClient();
  const {
    isPending,
    error,
    data: config,
  } = useQuery({
    queryKey: ['workspace.uploaderConfig', workspaceId],
    queryFn: () => {
      if (!workspaceId) {
        return Promise.reject(new Error('workspaceId is required'));
      }
      return axiosTenantInstance
        .get(`/workspaces/${workspaceId}/config`)
        .then(res => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({
      queryKey: ['workspace.uploaderConfig', workspaceId],
    });

  return { isPending, error, config, invalidate };
};

export const useUpdateUploaderConfigMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['workspace.uploaderConfig'],
    mutationFn: ({
      workspaceId,
      config,
    }: {
      workspaceId: string;
      config: WorkspaceConfig;
    }) => {
      return axiosTenantInstance
        .put(`/workspaces/${workspaceId}/config`, config)
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspace.uploaderConfig'] });
      notifications.show({
        title: 'Success',
        message: 'Uploader configuration updated successfully',
        color: 'green',
      });
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to update uploader configuration',
        color: 'red',
      });
    },
  });
};
