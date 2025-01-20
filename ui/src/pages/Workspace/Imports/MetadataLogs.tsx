import { CodeHighlight } from "@mantine/code-highlight";
import { Modal } from "@mantine/core";
import React from "react";

type MetadataLogsModalProps = {
    open: boolean
    onClose: () => void
    variant: 'metadata' | 'logs'
    logs?: string[],
    metadata?: Record<string, string>
}
const MetadataLogsModal: React.FC<MetadataLogsModalProps> = ({ open, onClose, variant, logs, metadata }) => {
    let code = ""
    if (variant === 'metadata') {
        code = JSON.stringify(metadata, null, 2)
    } else {
        code = JSON.stringify(logs, null, 2)
    }

    return (
        <Modal
            opened={open}
            title={variant === 'metadata' ? 'Metadata' : 'Logs'}
            onClose={onClose}
            size="xl"
        >
            <CodeHighlight code={code} language="json" />
        </Modal>
    )
};

export default MetadataLogsModal