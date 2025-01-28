import { CodeHighlight } from "@mantine/code-highlight";
import { Modal } from "@mantine/core";
import React from "react";

type MetadataModalProps = {
    open: boolean
    onClose: () => void
    metadata: Record<string, string>
}
export const MetadataModal: React.FC<MetadataModalProps> = ({ open, onClose, metadata }) => {
    const code = JSON.stringify(metadata, null, 2)

    return (
        <Modal
            opened={open}
            title="Metadata"
            onClose={onClose}
            size="xl"
        >
            <CodeHighlight code={code} language="json" />
        </Modal>
    )
};