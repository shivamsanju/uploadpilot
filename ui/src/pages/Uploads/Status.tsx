import { Loader, ThemeIcon, Tooltip } from "@mantine/core";
import {
    IconTrash,
    IconRosetteDiscountCheck, IconCancel,
    IconRosetteDiscountCheckFilled,
    IconProgressX, IconHelpCircle,
    IconTimeDuration0,
    IconAlertHexagon,
} from "@tabler/icons-react";

const statusConfig: Record<string, { color: string; icon: JSX.Element }> = {
    Started: { color: "blue", icon: <IconTimeDuration0 size={20} /> },
    Skipped: { color: "gray", icon: <IconProgressX size={20} /> },
    "In Progress": { color: "teal", icon: <Loader size={20} /> },
    Uploaded: { color: "green", icon: <IconRosetteDiscountCheck size={20} /> },
    Failed: { color: "red", icon: <IconAlertHexagon size={20} /> },
    Cancelled: { color: "gray", icon: <IconCancel size={20} /> },
    Processing: { color: "teal", icon: <Loader size={20} /> },
    "Processing Failed": { color: "red", icon: <IconAlertHexagon size={20} /> },
    "Processing Complete": { color: "green", icon: <IconRosetteDiscountCheckFilled size={20} /> },
    "Processing Cancelled": { color: "gray", icon: <IconCancel size={20} /> },
    Deleted: { color: "gray", icon: <IconTrash size={20} /> },
    Unknown: { color: "gray", icon: <IconHelpCircle size={20} /> },
};

export const UploadStatus = ({ status = "Unknown" }: { status?: string }) => {
    const { color, icon } = statusConfig[status] || statusConfig.Unknown;

    return (
        <Tooltip label={status}>
            <ThemeIcon color={color} variant="subtle">
                {icon}
            </ThemeIcon>
        </Tooltip>
    );
};
