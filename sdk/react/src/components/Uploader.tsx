import React, { useEffect, useState } from "react";
import Uppy from "@uppy/core";
import Webcam from "@uppy/webcam";
import { Dashboard } from "@uppy/react";
import RemoteSources from "@uppy/remote-sources";
import Audio from "@uppy/audio";
import ScreenCapture from "@uppy/screen-capture";
import ImageEditor from "@uppy/image-editor";
import GoldenRetriever from "@uppy/golden-retriever";
import Compressor from "@uppy/compressor";
import Informer from "@uppy/informer";
import Progress from "@uppy/progress-bar";
import StatusBar from "@uppy/status-bar";
import Tus from "@uppy/tus";

import "@uppy/core/dist/style.css";
import "@uppy/dashboard/dist/style.css";
import "@uppy/audio/dist/style.css";
import "@uppy/screen-capture/dist/style.css";
import "@uppy/image-editor/dist/style.css";

type UploaderProps = {
  workspaceId: string;
  backendEndpoint: string;
  height: number;
  width: number;
  theme: "auto" | "light" | "dark";
  showStatusBar?: boolean;
  showProgress?: boolean;
  metadata?: Record<string, string>;
  headers?: Record<string, string>;
  hideUploadButton?: boolean;
  hideCancelButton?: boolean;
  hideRetryButton?: boolean;
  hidePauseResumeButton?: boolean;
  hideProgressAfterFinish?: boolean;
  note?: string;
  singleFileFullScreen?: boolean;
  showSelectedFiles?: boolean;
  showRemoveButtonAfterComplete?: boolean;
  autoProceed?: boolean;
};

const Uploader: React.FC<UploaderProps> = ({
  workspaceId,
  backendEndpoint,
  height,
  width,
  theme,
  headers,
  metadata,
  showStatusBar = true,
  showProgress = true,
  hideUploadButton = false,
  hideCancelButton = false,
  hideRetryButton = false,
  hidePauseResumeButton = false,
  hideProgressAfterFinish = false,
  note = null,
  singleFileFullScreen = true,
  showSelectedFiles = true,
  autoProceed = false,
}) => {
  const [uppy, setUppy] = useState<any>();

  useEffect(() => {
    if (!workspaceId) return;
    fetch(`${backendEndpoint}/workspaces/${workspaceId}/config`)
      .then((response) => response.json())
      .then((config) => {
        const uppy = new Uppy({
          id: workspaceId,
          autoProceed: autoProceed,
          debug: false,
          restrictions: {
            maxFileSize: config?.maxFileSize,
            minFileSize: config?.minFileSize,
            maxNumberOfFiles: config?.maxNumberOfFiles,
            minNumberOfFiles: config?.minNumberOfFiles,
            allowedFileTypes:
              config?.allowedFileTypes && config.allowedFileTypes.length > 0
                ? config.allowedFileTypes
                : undefined,
            maxTotalFileSize: config?.maxTotalFileSize,
            requiredMetaFields:
              config?.requiredMetadataFields &&
              config.requiredMetadataFields.length > 0
                ? config.requiredMetadataFields
                : [],
          },
        });
        uppy.use(Informer);
        uppy.use(Progress);
        uppy.use(StatusBar);
        uppy.use(RemoteSources, {
          companionUrl: `${backendEndpoint}/remote`,
          sources: config.allowedSources.filter(
            (e: string) =>
              !["FileUpload", "Audio", "Webcamera", "ScreenCapture"].includes(
                e,
              ),
          ),
          companionAllowedHosts: [backendEndpoint],
        });
        uppy.use(Tus, {
          endpoint: `${backendEndpoint}/upload`,
          headers: {
            workspaceId: workspaceId,
            ...headers,
          },
        });
        if (config.allowedSources.includes("Audio")) uppy.use(Audio);
        if (config.allowedSources.includes("Webcamera")) uppy.use(Webcam);
        if (config.allowedSources.includes("ScreenCapture"))
          uppy.use(ScreenCapture);
        if (config.enableImageEditing) uppy.use(ImageEditor);
        if (config.useCompression) uppy.use(Compressor);
        if (config.useFaultTolerantMode) uppy.use(GoldenRetriever);
        if (metadata) uppy.setMeta(metadata);
        setUppy(uppy);
      });
  }, [workspaceId, backendEndpoint]);

  return (
    uppy && (
      <Dashboard
        uppy={uppy}
        height={height}
        width={width}
        theme={theme}
        hideUploadButton={hideUploadButton}
        hideCancelButton={hideCancelButton}
        hideRetryButton={hideRetryButton}
        hidePauseResumeButton={hidePauseResumeButton}
        hideProgressAfterFinish={hideProgressAfterFinish}
        note={note}
        singleFileFullScreen={singleFileFullScreen}
        showSelectedFiles={showSelectedFiles}
        showRemoveButtonAfterComplete={showSelectedFiles}
        showProgressDetails={showProgress}
        disableStatusBar={!showStatusBar}
        proudlyDisplayPoweredByUppy={false}
      />
    )
  );
};

export default Uploader;
