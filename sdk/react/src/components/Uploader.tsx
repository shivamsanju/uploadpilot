import React, { useEffect, useState } from 'react';
import Uppy from '@uppy/core';
import Webcam from '@uppy/webcam';
import { Dashboard } from '@uppy/react';
import RemoteSources from '@uppy/remote-sources';
import Audio from '@uppy/audio';
import ScreenCapture from '@uppy/screen-capture';
import ImageEditor from '@uppy/image-editor';
import GoldenRetriever from '@uppy/golden-retriever';
import Compressor from '@uppy/compressor';
import Informer from '@uppy/informer';
import Progress from '@uppy/progress-bar';
import StatusBar from '@uppy/status-bar';
import Tus from '@uppy/tus';

import '@uppy/core/dist/style.css';
import '@uppy/dashboard/dist/style.css';
import '@uppy/audio/dist/style.css';
import '@uppy/screen-capture/dist/style.css';
import '@uppy/image-editor/dist/style.css';

type UploaderProps = {
    uploaderId: string
    backendEndpoint: string
    h: number
    w: number
};
const Uploader: React.FC<UploaderProps> = ({ uploaderId, backendEndpoint, h, w }) => {
    const [uppy, setUppy] = useState<any>();
    const [theme, setTheme] = useState<"dark" | "light" | "auto">('auto');


    useEffect(() => {
        if (!uploaderId) return;
        const token = localStorage.getItem("token");
        fetch(`${backendEndpoint}/uploaders/${uploaderId}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => response.json())
            .then(data => {
                const config = data.config;
                const uppy = new Uppy({
                    autoProceed: true,
                    debug: true,
                    restrictions: {
                        maxFileSize: config.maxFileSize,
                        minFileSize: config.minFileSize,
                        maxNumberOfFiles: config.maxNumberOfFiles,
                        minNumberOfFiles: config.minNumberOfFiles,
                        allowedFileTypes: config.allowedFileTypes,
                        maxTotalFileSize: config.maxTotalFileSize,
                        requiredMetaFields: config.requiredMetadataFields
                    }
                });
                uppy.use(Informer);
                uppy.use(RemoteSources, {
                    companionUrl: "COMPANION_URL",
                    sources: config.allowedSources.filter((e: string) => !['FileUpload', 'Audio', 'Webcamera', 'ScreenCapture'].includes(e)),
                    companionAllowedHosts: []
                });
                uppy.use(Tus, {
                    endpoint: `${backendEndpoint}/imports`,
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'uploaderId': uploaderId
                    }
                });
                if (config.enableImageEditing) uppy.use(ImageEditor);
                if (config.useCompression) uppy.use(Compressor);
                if (config.useFaultTolerantMode) uppy.use(GoldenRetriever);
                if (config.allowedSources.includes('Audio')) uppy.use(Audio);
                if (config.allowedSources.includes('Webcamera')) uppy.use(Webcam);
                if (config.allowedSources.includes('ScreenCapture')) uppy.use(ScreenCapture);
                if (config.showProgressBar) uppy.use(Progress);
                if (config.showStatusBar) uppy.use(StatusBar);
                setUppy(uppy);
                setTheme(config.theme);
            })
    }, [uploaderId, backendEndpoint]);

    return <Dashboard
        uppy={uppy}
        height={h}
        width={w}
        theme={theme}
    />;
}

export default Uploader
