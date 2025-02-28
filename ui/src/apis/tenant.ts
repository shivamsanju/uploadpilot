import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { TENANT_ID_KEY } from '../constants/tenancy';
import type { TenantOnboardingRequest } from '../types/tenant';
import axiosInstance from '../utils/axios';

export const useOnboardTenant = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['tenants'],
    mutationFn: (data: TenantOnboardingRequest) => {
      return axiosInstance.post(`/tenants`, data).then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] });
      notifications.show({
        title: 'Success',
        message: 'Tenant onboarding successfull',
        color: 'green',
      });
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to onboard tenant',
        color: 'red',
      });
    },
  });
};

export const useGetActiveTenant = () => {
  const {
    isPending,
    error,
    data: session,
  } = useQuery({
    queryKey: ['activeTenant'],
    refetchInterval: 5000, // 5 seconds
    queryFn: async () => {
      return localStorage.getItem(TENANT_ID_KEY);
    },
  });

  return { isPending, error, session };
};

export const useSetActiveTenant = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ['activeTenant'],
    mutationFn: (tenantId: string) => {
      return axiosInstance
        .put(`/tenants/active`, { tenantId })
        .then(res => res.data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['activeTenant'] });
    },
  });
};
