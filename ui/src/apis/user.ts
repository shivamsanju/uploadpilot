import { useNavigate } from "react-router-dom";
import axiosInstance from "../utils/axios";
import { useQuery } from "@tanstack/react-query";

export const useGetSession = () => {
  const navigate = useNavigate();

  const {
    isPending,
    error,
    data: session,
  } = useQuery({
    queryKey: ["session"],
    refetchInterval: 60000,
    staleTime: 60000,
    queryFn: () => {
      return axiosInstance
        .get(`/session`)
        .then((res) => {
          if (res.status !== 200) {
            localStorage.removeItem("uploadpilottoken");
            navigate("/auth");
          }
          return res.data;
        })
        .catch(() => {
          localStorage.removeItem("uploadpilottoken");
          navigate("/auth");
        });
    },
  });

  return { isPending, error, session };
};
