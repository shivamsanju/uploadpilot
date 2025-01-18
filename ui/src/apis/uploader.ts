import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Uploader } from "../types/uploader";

export const useGetUploaders = ({ skip, limit, search }: {
    skip: number,
    limit: number,
    search: string
}) => {
    const queryClient = useQueryClient();

    const { isPending, error, data } = useQuery({
        queryKey: ['uploaders', skip, limit, search],
        queryFn: () =>
            axiosInstance
                .get(`/uploaders?skip=${skip}&limit=${limit}&search=${search}`)
                .then((res) => res.data)
    })

    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['uploaders', skip, limit, search] });
    return { isPending, error, uploaders: data?.records || [], totalRecords: data?.totalRecords || 0, invalidate }
}


export const useCreateUploaderMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['uploaders'],
        mutationFn: (wf: Uploader) => axiosInstance.post("/uploaders", wf).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['uploaders'] });
            notifications.show({
                title: "Success",
                message: "Uploader created successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to create uploader",
                color: "red",
            });
        },
    })
};



export const useGetUploaderDetailsById = (uploaderId: string) => {
    const { isPending, error, data: uploader } = useQuery({
        queryKey: ['uploaderDetails', uploaderId],
        queryFn: () => {
            if (!uploaderId) {
                return Promise.reject(new Error('uploaderId is required'));
            }
            return axiosInstance
                .get("/uploaders/" + uploaderId)
                .then((res) => res.data)
        }

    })

    return { isPending, error, uploader }
}

export const useUpdateUploaderConfigMutation = (uploaderId: string) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['uploaderDetails', uploaderId],
        mutationFn: (wf: Uploader) => axiosInstance.put(`/uploaders/${uploaderId}/config`, wf).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['uploaderDetails', uploaderId] });
            notifications.show({
                title: "Success",
                message: "Uploader updated successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to update uploader",
                color: "red",
            });
        },
    })
};

export const useGetAllAllowedSources = () => {
    const { isPending, error, data: allowedSources } = useQuery({
        queryKey: ['allowedSources'],
        queryFn: () =>
            axiosInstance
                .get("/uploaders/sources/allowed")
                .then((res) => res.data)
    })

    return { isPending, error, allowedSources }
}