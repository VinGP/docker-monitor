"use client";

import {useContainers} from "../hooks/useContainers";
import {format} from "date-fns";
import React, {useEffect, useState} from "react";
import {FaSort, FaSortDown, FaSortUp} from "react-icons/fa";
import {ru} from 'date-fns/locale'
import {AiOutlineCheckCircle, AiOutlineCloseCircle, AiOutlineDelete, AiOutlineLoading3Quarters} from "react-icons/ai"; // ✅ Импорт иконки загрузки


export function ContainersList() {
    const {containers, isLoading, isError, isFetching, clear, isClearing} = useContainers();
    const [sortConfig, setSortConfig] = useState<{ key: string | null; direction: "asc" | "desc" | "none" }>({
        key: null,
        direction: "none"
    });


    const [status, setStatus] = useState<"loading" | "success" | "error" | null>(null);

    useEffect(() => {
        if (isFetching) {
            setStatus("loading");
        } else if (isError) {
            setStatus("error");
        } else {
            setStatus("success");
        }
    }, [isFetching, isError]);


    const handleSort = (key: string) => {
        setSortConfig((prev) => {
            let direction: "asc" | "desc" | "none" = "asc";
            if (prev.key === key) {
                if (prev.direction === "asc") direction = "desc";
                else if (prev.direction === "desc") direction = "none";
                else direction = "asc";
            }
            return {key, direction};
        });
    };


    const sortedContainers = React.useMemo(() => {
        // if (!containers || sortConfig.direction === "none") return containers;
        if (!containers || sortConfig.direction === "none" || sortConfig.key === null) return containers;

        return [...containers].sort((a, b) => {
            if (sortConfig.key === "last_success") {
                const dateA = new Date(a[sortConfig.key] as string).getTime();
                const dateB = new Date(b[sortConfig.key] as string).getTime();
                return sortConfig.direction === "asc" ? dateA - dateB : dateB - dateA;
            }
            if (sortConfig.key === "ping_time") {
                return sortConfig.direction === "asc"
                    ? (a[sortConfig.key] as number) - (b[sortConfig.key] as number)
                    : (b[sortConfig.key] as number) - (a[sortConfig.key] as number);
            }

            return sortConfig.direction === "asc"
                // @ts-ignore
                ? String(a[sortConfig.key]).localeCompare(String(b[sortConfig.key]))
                // @ts-ignore
                : String(b[sortConfig.key]).localeCompare(String(a[sortConfig.key]));
        });
    }, [containers, sortConfig]);

    if (isLoading) {
        return (
            <div className="mt-4 card animate-pulse">
                <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
                <div className="space-y-3">
                    <div className="h-4 bg-gray-200 rounded w-full"></div>
                    <div className="h-4 bg-gray-200 rounded w-full"></div>
                    <div className="h-4 bg-gray-200 rounded w-full"></div>
                </div>
            </div>
        );
    }

    if (isError) {
        return (
            <div className="mt-4 card bg-red-50 border border-red-200">
                <p className="text-red-700">Error loading containers</p>
            </div>
        );
    }

    if (!containers) {
        return (
            <div className="mt-4 card bg-yellow-50 border border-yellow-200">
                <p className="text-yellow-700">No containers found</p>
            </div>
        );
    }

    return (
        <div className="mt-4">


            <div className="flex justify-between items-center mb-4">
                <div className="flex items-center space-x-2">
                    {status === "loading" && (
                        <>
                            <AiOutlineLoading3Quarters className="w-5 h-5 text-blue-500 animate-spin"/>
                            <span className="text-blue-500">Обновление данных...</span>
                        </>
                    )}
                    {status === "success" && (
                        <>
                            <AiOutlineCheckCircle className="w-5 h-5 text-green-500"/>
                            <span className="text-green-500">Данные обновлены</span>
                        </>
                    )}
                    {status === "error" && (
                        <>
                            <AiOutlineCloseCircle className="w-5 h-5 text-red-500"/>
                            <span className="text-red-500">Ошибка загрузки</span>
                        </>
                    )}
                </div>

                <button
                    onClick={() => clear()}
                    disabled={isClearing}
                    className=" flex items-center px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50"
                >
                    {isClearing ? (
                        <AiOutlineLoading3Quarters className="w-5 h-5 animate-spin mr-2"/>
                    ) : (
                        <AiOutlineDelete className="w-5 h-5 mr-2"/>
                    )}
                    Очистить
                </button>
            </div>


            <div className="overflow-x-auto rounded-lg">
                <table className="min-w-full divide-y divide-gray-300">
                    <thead>
                    <tr>
                        {[
                            {key: "ip_address", label: "IP Address"},
                            {key: "ping_time", label: "Ping Time (ms)"},
                            {key: "last_success", label: "Last Successful"}
                        ].map(({key, label}) => {
                            const Icon = sortConfig.key === key ? (sortConfig.direction === "asc" ? FaSortUp : sortConfig.direction === "desc" ? FaSortDown : FaSort) : FaSort;
                            return (
                                <th
                                    key={key}
                                    onClick={() => handleSort(key)}
                                    className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                                >
                                    <div className="flex items-center space-x-2">
                                        <span>{label}</span> <Icon className="w-4 h-4"/>
                                    </div>
                                </th>
                            );
                        })}
                    </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                    {sortedContainers?.map((container) => (
                        <tr key={container.ip_address} className="hover:bg-gray-100">
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                {container.ip_address}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {container.ping_time === null ? "–" : container.ping_time}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{container.last_success === null || container.last_success === "" ? "–" : format(new Date(container.last_success), 'dd MMMM yyyy, HH:mm:ss.SSS', {locale: ru})}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
