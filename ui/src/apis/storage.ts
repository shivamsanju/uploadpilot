import { notifications } from "@mantine/notifications";
import { Datastore } from "../types/connector";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const useGetStorageConnectors = () => {
    const { isPending, error, data: connectors } = useQuery({
        queryKey: ['storageConnectors'],
        queryFn: () =>
            axiosInstance
                .get("/storage/connectors")
                .then((res) => res.data)
    })

    return { isPending, error, connectors }
}



export const useCreateStorageConnectorMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['storageConnectors'],
        mutationFn: (sc: any) => axiosInstance.post("/storage/connectors", sc).then((res) => res.data),
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

export const useGetDataStores = () => {
    const { isPending, error, data: connectors } = useQuery({
        queryKey: ['dataStores'],
        queryFn: () =>
            axiosInstance
                .get("/storage/datastores")
                .then((res) => res.data)
    })

    return { isPending, error, connectors }
}


export const useCreateDataStoreMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['dataStores'],
        mutationFn: (sc: Datastore) => axiosInstance.post("/storage/datastores", sc).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['dataStores'] });
            notifications.show({
                title: "Success",
                message: "Datastore created successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to create datastore",
                color: "red",
            });
        },
    })
};