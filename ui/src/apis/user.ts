import axiosInstance from "../utils/axios";
import { useQuery } from "@tanstack/react-query";

export const useGetCurrentUserDetails = () => {
    const { isPending, error, data: me } = useQuery({
        queryKey: ['me'],
        queryFn: () => {
            return axiosInstance
                .get(`/users/me`)
                .then((res) => {
                    return {
                        ...res.data,
                        image: "https://avatar.iran.liara.run/public/33"
                    }
                })
        }
    })

    return { isPending, error, me }
}