import axiosInstance from "../utils/axios";
import { useInfiniteQuery, useQueryClient } from "@tanstack/react-query";

interface UploadParams {
    workspaceId: string;
    batchSize: number;
    search?: string;
}

interface UploadResponse {
    records: any[];
    totalRecords: number;
}

export const useGetUploads = ({
    workspaceId,
    batchSize = 50,
    search = '',
}: UploadParams) => {
    const queryClient = useQueryClient();

    const {
        isPending,
        error,
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        isFetching,
        isFetchNextPageError
    } = useInfiniteQuery({
        queryKey: ['uploads', { workspaceId, search }],
        staleTime: 30 * 1000,
        queryFn: ({ pageParam = 0 }) => {
            if (!workspaceId) {
                throw new Error('workspaceId is required');
            }

            const skipValue = pageParam * batchSize;
            const searchParam = search ? `&search=${encodeURIComponent(search)}` : '';

            return axiosInstance.get<UploadResponse>(
                `/workspaces/${workspaceId}/uploads?skip=${skipValue}&limit=${batchSize}${searchParam}`
            ).then(res => res.data);
        },
        getNextPageParam: (lastPage, allPages) => {
            const totalPages = Math.ceil(lastPage.totalRecords / batchSize);
            const nextPage = allPages.length;
            return nextPage < totalPages ? nextPage : undefined;
        },
        initialPageParam: 0,
    });

    const allUploads = data?.pages.flatMap(page => page?.records || []) ?? [];
    const totalRecords = data?.pages[0]?.totalRecords ?? 0;

    const invalidate = async () => {
        // Remove all existing data for this query
        await queryClient.cancelQueries({
            queryKey: ['uploads', { workspaceId, search }]
        });

        // Reset the query to its initial state
        queryClient.resetQueries({
            queryKey: ['uploads', { workspaceId, search }],
        });

        // Fetch only the first page
        return queryClient.fetchQuery({
            queryKey: ['uploads', { workspaceId, search }],
            queryFn: () => {
                const searchParam = search ? `&search=${encodeURIComponent(search)}` : '';
                return axiosInstance.get<UploadResponse>(
                    `/workspaces/${workspaceId}/uploads?skip=0&limit=${batchSize}${searchParam}`
                ).then(res => res.data);
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
        isFetching
    };
};