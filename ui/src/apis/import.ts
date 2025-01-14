import axiosInstance from "../utils/axios";
import { useQuery } from "@tanstack/react-query";

export const useGetImports = (uploaderId: string) => {
    const { isPending, error, data: imports } = useQuery({
        queryKey: ['imports', uploaderId],
        queryFn: () => {
            if (!uploaderId) {
                return Promise.reject(new Error('uploaderId is required'));
            }
            return axiosInstance
                .get(`/uploaders/${uploaderId}/imports`)
                .then((res) => res.data)
        }
    })

    return { isPending, error, imports }
}