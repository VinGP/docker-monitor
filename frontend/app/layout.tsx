import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Header } from "../components/Header";
import Providers from "./providers";
import React from "react";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
    title: "Docker Container Monitor",
    description: "Real-time monitoring of Docker containers status",
};

export default function RootLayout({
                                       children,
                                   }: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en">
        <body className={inter.className}>
        <Providers>
            <div className="min-h-screen bg-gray-50 sm:px-6 lg:px-8">
                <Header />
                <div className="flex ">
                    <main className="flex-1 p-6 pt-0">{children}</main>
                </div>
            </div>
        </Providers>
        </body>
        </html>
    );
}
