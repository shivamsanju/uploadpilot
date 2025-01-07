export type Codebase = {
    id: string;
    name: string;
    description: string;
    tags: string[];
    status: string;
    lang: string;
    basePath: string
    updatedAt: number;
    createdAt: number;
    createdBy: string;
    updatedBy: string;
    members: string[];
    gitUrl?: string;
}
