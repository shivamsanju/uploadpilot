import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const useGetStorageConnectors = ({ skip, limit, search }: {
    skip: number,
    limit: number,
    search: string
}) => {
    const queryClient = useQueryClient();

    const { isPending, error, data } = useQuery({
        queryKey: ['storageConnectors', skip, limit, search],
        queryFn: () =>
            axiosInstance
                .get(`/storageConnectors?skip=${skip}&limit=${limit}&search=${search}`)
                .then((res) => res.data)
    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['storageConnectors', skip, limit, search] });

    return { isPending, error, connectors: data?.records || [], totalRecords: data?.totalRecords || 0, invalidate }
}



export const useCreateStorageConnectorMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['storageConnectors'],
        mutationFn: (sc: any) => axiosInstance.post("/storageConnectors", sc).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['storageConnectors'] });
            notifications.show({
                title: "Success",
                message: "Connector created successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to create connector",
                color: "red",
            });
        },
    })
};