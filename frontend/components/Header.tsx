import React from 'react';

export const Header = () => {
    return (
        <header className="bg-white shadow-sm">
            <div className="mx-auto max-w-8xl px-4 sm:px-6 lg:px-8">
                <div className="flex h-16 justify-between items-center">
                    <div className="flex">
                        <div className="flex flex-shrink-0 items-center">
                            <h1 className="text-xl font-bold text-primary">Docker Monitor</h1>
                        </div>
                    </div>
                    <div>
                        <h2><a className="hover:text-accent" href="https://github.com/VinGP">Github</a></h2>
                    </div>
                </div>

            </div>
        </header>
    );
};
