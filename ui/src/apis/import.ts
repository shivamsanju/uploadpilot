import axiosInstance from "../utils/axios";
import { useInfiniteQuery, useQueryClient } from "@tanstack/react-query";

interface ImportParams {
    workspaceId: string;
    batchSize: number;
    search?: string;
}

interface ImportResponse {
    records: any[];
    totalRecords: number;
}

export const useGetImports = ({
    workspaceId,
    batchSize = 50,
    search = '',
}: ImportParams) => {
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
        queryKey: ['imports', { workspaceId, search }],
        staleTime: 30 * 1000,
        queryFn: ({ pageParam = 0 }) => {
            if (!workspaceId) {
                throw new Error('workspaceId is required');
            }

            const skipValue = pageParam * batchSize;
            const searchParam = search ? `&search=${encodeURIComponent(search)}` : '';

            return axiosInstance.get<ImportResponse>(
                `/workspaces/${workspaceId}/imports?skip=${skipValue}&limit=${batchSize}${searchParam}`
            ).then(res => res.data);
        },
        getNextPageParam: (lastPage, allPages) => {
            const totalPages = Math.ceil(lastPage.totalRecords / batchSize);
            const nextPage = allPages.length;
            return nextPage < totalPages ? nextPage : undefined;
        },
        initialPageParam: 0,
    });

    const allImports = data?.pages.flatMap(page => page?.records || []) ?? [];
    const totalRecords = data?.pages[0]?.totalRecords ?? 0;

    const invalidate = async () => {
        // Remove all existing data for this query
        await queryClient.cancelQueries({
            queryKey: ['imports', { workspaceId, search }]
        });

        // Reset the query to its initial state
        queryClient.resetQueries({
            queryKey: ['imports', { workspaceId, search }],
        });

        // Fetch only the first page
        return queryClient.fetchQuery({
            queryKey: ['imports', { workspaceId, search }],
            queryFn: () => {
                const searchParam = search ? `&search=${encodeURIComponent(search)}` : '';
                return axiosInstance.get<ImportResponse>(
                    `/workspaces/${workspaceId}/imports?skip=0&limit=${batchSize}${searchParam}`
                ).then(res => res.data);
            },
        });
    };

    return {
        isPending,
        error,
        isFetchNextPageError,
        imports: allImports,
        totalRecords,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        invalidate,
        isFetching
    };
};