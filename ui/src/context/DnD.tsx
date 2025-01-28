import { createContext, useState, ReactNode } from 'react';

interface DataTransfer {
    [key: string]: any;
}

interface DnDContextType {
    type: string;
    setType: (type: string) => void;
    dataTransfer: DataTransfer;
    setDataTransfer: (dataTransfer: DataTransfer) => void;
}

export const DnDContext = createContext<DnDContextType>({
    type: "",
    setType: () => { },
    dataTransfer: {},
    setDataTransfer: () => { },
});

export const DnDProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [type, setType] = useState<string>("");
    const [dataTransfer, setDataTransfer] = useState<DataTransfer>({});

    return (
        <DnDContext.Provider value={{ type, setType, dataTransfer, setDataTransfer }}>
            {children}
        </DnDContext.Provider>
    );
}
