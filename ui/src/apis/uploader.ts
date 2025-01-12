import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Uploader } from "../types/uploader";

export const useGetUploaders = () => {
    const { isPending, error, data: uploaders } = useQuery({
        queryKey: ['uploaders'],
        queryFn: () =>
            axiosInstance
                .get("/uploaders")
                .then((res) => res.data)
    })

    return { isPending, error, uploaders }
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

export const useGetAllAllowedSources = () => {
    const { isPending, error, data: allowedSources } = useQuery({
        queryKey: ['allowedSources'],
        queryFn: () =>
            axiosInstance
                .get("/uploaders/allowedSources")
                .then((res) => res.data)
    })

    return { isPending, error, allowedSources }
}