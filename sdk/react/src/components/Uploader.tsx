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
    height: number
    width: number
    theme: 'auto' | 'light' | 'dark'
    metadata?: Record<string, string>
    headers?: Record<string, string>
};
const Uploader: React.FC<UploaderProps> = ({ uploaderId, backendEndpoint, height, width, theme, headers, metadata }) => {
    const [uppy, setUppy] = useState<any>();



    useEffect(() => {
        if (!uploaderId) return;
        fetch(`${backendEndpoint}/uploaders/${uploaderId}`)
            .then(response => response.json())
            .then(data => {
                const config = data.config;
                console.log(config);
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
                    companionUrl: `${backendEndpoint}/remote`,
                    sources: config.allowedSources.filter((e: string) => !['FileUpload', 'Audio', 'Webcamera', 'ScreenCapture'].includes(e)),
                    companionAllowedHosts: [
                        backendEndpoint
                    ],
                });
                uppy.use(Tus, {
                    endpoint: `${backendEndpoint}/upload`,
                    headers: {
                        'uploaderId': uploaderId,
                        ...headers
                    },
                });
                if (metadata) uppy.setMeta(metadata);
                if (config.enableImageEditing) uppy.use(ImageEditor);
                if (config.useCompression) uppy.use(Compressor);
                if (config.useFaultTolerantMode) uppy.use(GoldenRetriever);
                if (config.allowedSources.includes('Audio')) uppy.use(Audio);
                if (config.allowedSources.includes('Webcamera')) uppy.use(Webcam);
                if (config.allowedSources.includes('ScreenCapture')) uppy.use(ScreenCapture);
                if (config.showProgress) uppy.use(Progress);
                if (config.showStatusBar) uppy.use(StatusBar);
                setUppy(uppy);
            })
    }, [uploaderId, backendEndpoint]);

    return uppy && <Dashboard
        uppy={uppy}
        height={height}
        width={width}
        theme={theme}
        proudlyDisplayPoweredByUppy={false}
    />;
}

export default Uploader
