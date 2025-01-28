import { CodeHighlight } from "@mantine/code-highlight";
import { LoadingOverlay, Modal } from "@mantine/core";
import React from "react";
import { UseGetUploadLogs } from "../../apis/upload";
import { ErrorCard } from "../../components/ErrorCard/ErrorCard";

const formatLogs = (logs: any[]) => {
    let str = ""
    for (const log of logs) {
        str += `${log?.timestamp} | [${log?.level?.toUpperCase()}] ${log?.message}\n`;
    }
    return str
}
type LogsModalProps = {
    open: boolean
    onClose: () => void
    uploadId: string
    workspaceId: string
}
export const LogsModal: React.FC<LogsModalProps> = ({ open, onClose, uploadId, workspaceId }) => {

    const { isPending, error, data } = UseGetUploadLogs(workspaceId, uploadId)


    return (
        <Modal
            opened={open}
            title="Logs"
            onClose={onClose}
            size="100%"
        >
            {error ? <ErrorCard message={error.message} title={error.name} /> : (
                <>
                    <LoadingOverlay visible={isPending} overlayProps={{ radius: 'sm', blur: 1 }} />
                    <CodeHighlight mih={300} code={formatLogs(data || [])} />
                </>
            )}
        </Modal>
    )
};
