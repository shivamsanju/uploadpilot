import axiosInstance from "../utils/axios";
import { useQuery } from "@tanstack/react-query";

export const useGetSession = () => {

    const { isPending, error, data: session } = useQuery({
        queryKey: ['session'],
        refetchInterval: 60000,
        staleTime: 60000,
        queryFn: () => {
            return axiosInstance
                .get(`/session`)
                .then((res) => res.data)
        }
    })


    return { isPending, error, session }
}