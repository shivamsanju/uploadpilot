import { useMantineColorScheme } from "@mantine/core";
import { useState } from "react";

export const useSettingsProps = () => {
  const { colorScheme } = useMantineColorScheme();
  const [height, setHeight] = useState<number>(400);
  const [width, setWidth] = useState<number>(350);
  const [theme, setTheme] = useState<"auto" | "light" | "dark">(colorScheme);
  const [showStatusBar, setShowStatusBar] = useState<boolean>(true);
  const [autoProceed, setAutoProceed] = useState<boolean>(false);
  const [showProgress, setShowProgress] = useState<boolean>(true);
  const [hideUploadButton, setHideUploadButton] = useState<boolean>(false);
  const [hideCancelButton, setHideCancelButton] = useState<boolean>(false);
  const [hideRetryButton, setHideRetryButton] = useState<boolean>(false);
  const [hidePauseResumeButton, setHidePauseResumeButton] =
    useState<boolean>(false);
  const [hideProgressAfterFinish, setHideProgressAfterFinish] =
    useState<boolean>(false);
  const [note, setNote] = useState<string>("");
  const [singleFileFullScreen, setSingleFileFullScreen] =
    useState<boolean>(true);
  const [showSelectedFiles, setShowSelectedFiles] = useState<boolean>(true);
  const [showRemoveButtonAfterComplete, setShowRemoveButtonAfterComplete] =
    useState<boolean>(true);

  return {
    height,
    setHeight,
    width,
    setWidth,
    theme,
    setTheme,
    showStatusBar,
    setShowStatusBar,
    autoProceed,
    setAutoProceed,
    showProgress,
    setShowProgress,
    hideUploadButton,
    setHideUploadButton,
    hideCancelButton,
    setHideCancelButton,
    hideRetryButton,
    setHideRetryButton,
    hidePauseResumeButton,
    setHidePauseResumeButton,
    hideProgressAfterFinish,
    setHideProgressAfterFinish,
    note,
    setNote,
    singleFileFullScreen,
    setSingleFileFullScreen,
    showSelectedFiles,
    setShowSelectedFiles,
    showRemoveButtonAfterComplete,
    setShowRemoveButtonAfterComplete,
  };
};
