export const isDataNode = (key: string) => {
    return ["webhook"].includes(key);
}