import axiosInstance from "../utils/axios";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export const useGetImports = ({ uploaderId, skip, limit, search }: {
    uploaderId: string,
    skip: number,
    limit: number,
    search: string
}) => {
    const queryClient = useQueryClient();

    const { isPending, error, data } = useQuery({
        queryKey: ['imports', uploaderId, skip, limit, search],
        queryFn: () => {
            if (!uploaderId) {
                return Promise.reject(new Error('uploaderId is required'));
            }
            return axiosInstance
                .get(`/uploaders/${uploaderId}/imports?skip=${skip}&limit=${limit}&search=${search}`)
                .then((res) => res.data)
        }
    })
    const invalidate = () => queryClient.invalidateQueries({ queryKey: ['imports', uploaderId, skip, limit, search] });

    return { isPending, error, imports: data?.records || [], totalRecords: data?.totalRecords || 0, invalidate }
}