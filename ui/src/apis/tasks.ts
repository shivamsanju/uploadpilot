import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Task } from '../types/tasks';
import axiosInstance from '../utils/axios';

export const useGetTasks = (workspaceId: string, processorId: string) => {
  const queryClient = useQueryClient();
  const {
    isPending,
    error,
    isFetching,
    data: tasks,
  } = useQuery<Task[]>({
    queryKey: ['tasks', workspaceId, processorId],
    queryFn: () => {
      if (!workspaceId) {
        return Promise.reject(new Error('workspaceId is required'));
      }
      return axiosInstance
        .get(`/workspaces/${workspaceId}/processors/${processorId}/tasks`)
        .then(res => res.data);
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({
      queryKey: ['tasks', workspaceId, processorId],
    });

  return { isPending, error, tasks, invalidate, isFetching };
};

export const useSaveTasksMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['tasks'],
    mutationFn: ({
      workspaceId,
      processorId,
      tasks,
    }: {
      workspaceId: string;
      processorId: string;
      tasks: Task[];
    }) => {
      return axiosInstance
        .put(
          `/workspaces/${workspaceId}/processors/${processorId}/tasks`,
          tasks,
        )
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      notifications.show({
        title: 'Success',
        message: 'Tasks saved successfully',
        color: 'green',
      });
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to save tasks',
        color: 'red',
      });
    },
  });
};
