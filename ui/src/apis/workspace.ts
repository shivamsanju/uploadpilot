import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateWorkspaceData } from '../types/workspace';
import { axiosTenantInstance } from '../utils/axios';

export const useGetWorkspaces = () => {
  const queryClient = useQueryClient();

  const { isPending, error, data } = useQuery({
    queryKey: ['workspaces'],
    queryFn: () => axiosTenantInstance.get(`/workspaces`).then(res => res.data),
  });

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ['workspaces'] });
  return { isPending, error, workspaces: data, invalidate };
};

export const useCreateWorkspaceMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['workspaces'],
    mutationFn: (data: CreateWorkspaceData) =>
      axiosTenantInstance.post('/workspaces', data).then(res => res.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
      notifications.show({
        title: 'Success',
        message: 'Workspace created successfully',
        color: 'green',
      });
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to create workspace',
        color: 'red',
      });
    },
  });
};

export const useGetUsersInWorkspace = (workspaceId: string) => {
  const queryClient = useQueryClient();

  const { isPending, error, data, isFetching } = useQuery({
    queryKey: ['workspaceUsers', workspaceId],
    queryFn: () => {
      if (!workspaceId) {
        return Promise.reject(new Error('workspaceId is required'));
      }
      return axiosTenantInstance
        .get(`workspaces/${workspaceId}/users`)
        .then(res => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({
      queryKey: ['workspaceUsers', workspaceId],
    });
  return { isPending, error, users: data, invalidate, isFetching };
};

export const useAddUserToWorkspaceMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['workspaceUsers'],
    mutationFn: ({
      workspaceId,
      email,
      role,
    }: {
      workspaceId: string;
      email: string;
      role: string;
    }) => {
      return axiosTenantInstance
        .post(`/workspaces/${workspaceId}/users`, { email, role })
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaceUsers'] });
      notifications.show({
        title: 'Success',
        message: 'User added to workspace successfully',
        color: 'green',
      });
    },
    onError: (error: any) => {
      notifications.show({
        title: 'Error',
        message: `Failed to add user to workspace. Reason: ${
          error?.response?.data?.message || error.message
        }`,
        color: 'red',
      });
    },
  });
};

export const useRemoveUserFromWorkspaceMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['workspaceUsers'],
    mutationFn: ({
      workspaceId,
      userId,
    }: {
      workspaceId: string;
      userId: string;
    }) => {
      return axiosTenantInstance
        .delete(`/workspaces/${workspaceId}/users/${userId}`)
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaceUsers'] });
      notifications.show({
        title: 'Success',
        message: 'User removed from workspace successfully',
        color: 'green',
      });
    },
    onError: (error: any) => {
      notifications.show({
        title: 'Error',
        message:
          'Failed to remove user from workspace. Reason: ' +
            error?.response?.data?.message || error.message,
        color: 'red',
      });
    },
  });
};

export const useEditUserInWorkspaceMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['workspaceUsers'],
    mutationFn: ({
      workspaceId,
      userId,
      role,
    }: {
      workspaceId: string;
      userId: string;
      role: string;
    }) => {
      return axiosTenantInstance
        .put(`/workspaces/${workspaceId}/users/${userId}`, { role })
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaceUsers'] });
      notifications.show({
        title: 'Success',
        message: 'User role modified successfully',
        color: 'green',
      });
    },
    onError: (error: any) => {
      console.log(error);
      notifications.show({
        title: 'Error',
        message: `Failed to modify user role. Reason: ${
          error?.response?.data?.message || error.message
        }`,
        color: 'red',
      });
    },
  });
};
