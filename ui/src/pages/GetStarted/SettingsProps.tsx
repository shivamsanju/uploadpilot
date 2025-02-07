import { useMantineColorScheme } from "@mantine/core";
import { useViewportSize } from "@mantine/hooks";
import { useEffect, useState } from "react";

const _useSettingsProps = () => {
  const { width: screenWidth } = useViewportSize();

  const { colorScheme } = useMantineColorScheme();
  const [height, setHeight] = useState<number>(screenWidth > 600 ? 700 : 500);
  const [width, setWidth] = useState<number>(
    screenWidth > 600 ? 600 : screenWidth
  );
  const [theme, setTheme] = useState<"auto" | "light" | "dark">(colorScheme);
  const [primaryColor, setPrimaryColor] = useState<string | undefined>();
  const [textColor, setTextColor] = useState<string | undefined>();
  const [hoverColor, setHoverColor] = useState<string | undefined>();
  const [noteColor, setNoteColor] = useState<string | undefined>();
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

  useEffect(() => {
    setWidth(screenWidth > 600 ? 600 : screenWidth);
    setHeight(screenWidth > 600 ? 700 : 500);
  }, [screenWidth]);

  return {
    height,
    setHeight,
    width,
    setWidth,
    theme,
    setTheme,
    primaryColor,
    setPrimaryColor,
    textColor,
    setTextColor,
    hoverColor,
    setHoverColor,
    noteColor,
    setNoteColor,
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

export const useSettingsProps = () => {
  const x = _useSettingsProps();

  return x;
};
