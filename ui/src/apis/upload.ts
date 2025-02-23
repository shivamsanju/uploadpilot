import axiosInstance from "../utils/axios";
import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import {
  areBracketsBalanced,
  getFilterParams,
  getSearchParam,
} from "../utils/utility";
import { notifications } from "@mantine/notifications";

interface UploadParams {
  workspaceId: string;
  batchSize: number;
  search?: string;
  filter?: Record<string, string[]>;
}

interface UploadResponse {
  records: any[];
  totalRecords: number;
}

export const useGetUploads = ({
  workspaceId,
  batchSize = 50,
  search = "",
  filter = {},
}: UploadParams) => {
  const queryClient = useQueryClient();

  let formattedSearch =
    search?.length > 2 && areBracketsBalanced(search) ? search : "";

  const {
    isPending,
    error,
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isFetching,
    isFetchNextPageError,
  } = useInfiniteQuery({
    queryKey: ["uploads", { workspaceId, formattedSearch, filter }],
    staleTime: 1 * 1000,
    queryFn: ({ pageParam = 0 }) => {
      if (!workspaceId) {
        throw new Error("workspaceId is required");
      }
      const offset = pageParam * batchSize;
      const searchParam = getSearchParam(formattedSearch);

      const filterParam = getFilterParams(filter);

      return axiosInstance
        .get<UploadResponse>(
          `/workspaces/${workspaceId}/uploads?offset=${offset}&limit=${batchSize}${searchParam}${filterParam}`
        )
        .then((res) => res.data);
    },
    getNextPageParam: (lastPage, allPages) => {
      const totalPages = Math.ceil(lastPage.totalRecords / batchSize);
      const nextPage = allPages.length;
      return nextPage < totalPages ? nextPage : undefined;
    },
    initialPageParam: 0,
  });

  const allUploads = data?.pages.flatMap((page) => page?.records || []) ?? [];
  const totalRecords = data?.pages[0]?.totalRecords ?? 0;

  const invalidate = async () => {
    // Remove all existing data for this query
    await queryClient.cancelQueries({
      queryKey: ["uploads", { workspaceId, search }],
    });

    // Reset the query to its initial state
    queryClient.resetQueries({
      queryKey: ["uploads", { workspaceId, search }],
    });

    // Fetch only the first page
    return queryClient.fetchQuery({
      queryKey: ["uploads", { workspaceId, search }],
      queryFn: () => {
        const searchParam = search
          ? `&search=${encodeURIComponent(search)}`
          : "";
        return axiosInstance
          .get<UploadResponse>(
            `/workspaces/${workspaceId}/uploads?skip=0&limit=${batchSize}${searchParam}`
          )
          .then((res) => res.data);
      },
    });
  };

  return {
    isPending,
    error,
    isFetchNextPageError,
    uploads: allUploads,
    totalRecords,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    invalidate,
    isFetching,
  };
};

export const useDownloadUploadedFile = (workspaceId: string) => {
  return useMutation({
    mutationKey: ["downloadFile", workspaceId],
    mutationFn: async ({ uploadId }: { uploadId: string }) => {
      if (!workspaceId || !uploadId) {
        throw new Error("workspaceId and uploadId are required");
      }
      const response = await axiosInstance.get(
        `/workspaces/${workspaceId}/uploads/${uploadId}/download`
      );
      return response.data;
    },
  });
};

export const useTriggerProcessUpload = (workspaceId: string) => {
  return useMutation({
    mutationKey: ["processUpload", workspaceId],
    mutationFn: async ({ uploadId }: { uploadId: string }) => {
      if (!workspaceId || !uploadId) {
        throw new Error("workspaceId and uploadId are required");
      }
      const response = await axiosInstance.post(
        `/workspaces/${workspaceId}/uploads/${uploadId}/process`
      );
      return response.data;
    },
    onSuccess: (_, { uploadId }) => {
      notifications.show({
        title: "Success",
        message: `Processing started for upload ${uploadId}`,
        color: "green",
      });
    },
    onError: (_, { uploadId }) => {
      notifications.show({
        title: "Error",
        message: `Failed to start processing for upload ${uploadId}`,
        color: "red",
      });
    },
  });
};
