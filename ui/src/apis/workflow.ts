import { notifications } from "@mantine/notifications";
import axiosInstance from "../utils/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Workflow } from "../types/workflow";

export const useGetWorkflows = () => {
    const { isPending, error, data: workflows } = useQuery({
        queryKey: ['workflows'],
        queryFn: () =>
            axiosInstance
                .get("/workflows")
                .then((res) => res.data)
    })

    return { isPending, error, workflows }
}

export const useCreateWorkflowMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['workflows'],
        mutationFn: (wf: Workflow) => axiosInstance.post("/workflows", wf).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
            notifications.show({
                title: "Success",
                message: "Workflow created successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to create workflow",
                color: "red",
            });
        },
    })
};
