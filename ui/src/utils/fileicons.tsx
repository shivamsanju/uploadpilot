import {
    IconFileTypePng,
    IconFileTypePdf,
    IconFileTypeJpg,
    IconFileTypeCsv,
    IconFileTypeZip,
    IconFileTypeDoc,
    IconFileTypeDocx,
    IconFileTypeXls,
    IconFileTypePpt,
    IconFileTypeSvg,
    IconGif,
    IconImageInPicture,
    IconVideo,
    IconDeviceAudioTape,
    IconAutomation,
    IconFileUnknown,
} from "@tabler/icons-react";
import { ReactNode } from "react";

export const MIME_TYPE_ICONS: { [key: string]: any } = {
    "image/png": IconFileTypePng,
    "image/gif": IconGif,
    "image/jpeg": IconFileTypeJpg,
    "image/svg+xml": IconFileTypeSvg,
    "image/webp": IconImageInPicture,
    "image/avif": IconImageInPicture,
    "image/heic": IconImageInPicture,
    "image/heif": IconImageInPicture,
    "video/mp4": IconVideo,
    "video/webm": IconVideo,
    "video/quicktime": IconVideo,
    "video/x-matroska": IconVideo,
    "video/x-msvideo": IconVideo,
    "video/x-flv": IconVideo,
    "video/3gpp": IconVideo,
    "audio/mpeg": IconDeviceAudioTape,
    "audio/wav": IconDeviceAudioTape,
    "audio/ogg": IconDeviceAudioTape,
    "audio/flac": IconDeviceAudioTape,
    "application/zip": IconFileTypeZip,
    "application/x-rar": IconFileTypeZip,
    "application/x-7z-compressed": IconFileTypeZip,
    "text/csv": IconFileTypeCsv,
    "application/pdf": IconFileTypePdf,
    "application/msword": IconFileTypeDocx,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": IconFileTypeDoc,
    "application/vnd.ms-excel": IconFileTypeXls,
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": IconFileTypeXls,
    "application/vnd.ms-powerpoint": IconFileTypePpt,
    "application/vnd.openxmlformats-officedocument.presentationml.presentation": IconFileTypePpt,
    "application/vnd.microsoft.portable-executable": IconAutomation,
};


export const getFileIcon = (mimeType: string, size?: number): ReactNode => {
    const IconComponent = MIME_TYPE_ICONS[mimeType];
    return IconComponent ? <IconComponent size={size} /> : <IconFileUnknown size={size} />;
};