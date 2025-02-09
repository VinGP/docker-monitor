"use client";

import {useMutation, useQuery} from "@tanstack/react-query";
import {clearContainers, getContainers} from "../lib/api";

export const useContainers = () => {
    const {
        data: containers,
        isLoading,
        error,
        refetch,
        isFetching
    } = useQuery({
        queryKey: ['containers'],
        queryFn: getContainers,
        refetchInterval: 2000,
    });

    const {mutate: clear, status: clearStatus} = useMutation<void, Error, void>({
        onSuccess: () => refetch(),

        mutationFn: () => clearContainers(),

    });

    const isClearing = clearStatus === 'pending';


    return {
        containers,
        isLoading,
        isError: !!error,
        refresh: refetch,
        isFetching,
        isClearing,
        clear
    };
};