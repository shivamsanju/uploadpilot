import { ImportPolicy } from "../types/importpolicy";
import axiosInstance from "../utils/axios";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useMutation } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";

export const useGetImportPolicies = () => {
    const { isPending, error, data: importPolicies } = useQuery({
        queryKey: ['importPolicies'],
        queryFn: () =>
            axiosInstance
                .get("/importPolicies")
                .then((res) => res.data)
    })

    return { isPending, error, importPolicies }
}


export const useGetImportPolicyDetails = (importpolicyId: string) => {
    const { isPending, error, data: importPolicy } = useQuery<any, Error, ImportPolicy>({
        queryKey: ['importPolicyDetails'],
        queryFn: () => {
            if (!importpolicyId) {
                return Promise.reject(new Error('importpolicyId is required'));
            }
            return axiosInstance
                .get("/importPolicies/" + importpolicyId)
                .then((res) => {
                    console.log(res)
                    return res.data
                })
        },
        enabled: !!importpolicyId,
    })

    return { isPending, error, importPolicy }
}


export const useCreateImportPolicyMutation = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationKey: ['importPolicies'],
        mutationFn: (newPolicy: ImportPolicy) => axiosInstance.post("/importPolicies", newPolicy).then((res) => res.data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['importPolicies'] });
            notifications.show({
                title: "Success",
                message: "Import policy added successfully",
                color: "green",
            });
        },
        onError: () => {
            notifications.show({
                title: "Error",
                message: "Failed to add import policy",
                color: "red",
            });
        },
    })
};

