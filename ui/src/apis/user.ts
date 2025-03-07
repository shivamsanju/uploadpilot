import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import Session from 'supertokens-web-js/recipe/session';
import { axiosBaseInstance } from '../utils/axios';

export const useGetUserDetails = () => {
  const {
    isPending,
    error,
    data: user,
  } = useQuery({
    queryKey: ['user'],
    refetchInterval: 1200000,
    queryFn: async () => {
      const resp = await axiosBaseInstance.get(`/user`);
      return resp.data;
    },
  });

  return { isPending, error, user };
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['user'],
    mutationFn: (data: Record<string, any>) => {
      return axiosBaseInstance.put(`/user`, data).then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] });
    },
  });
};

export const useGetSession = () => {
  const {
    isPending,
    error,
    data: session,
  } = useQuery({
    queryKey: ['session'],
    queryFn: async () => {
      return await Session.doesSessionExist();
    },
  });

  return { isPending, error, session };
};
