import {
  ColorInput,
  Group,
  NumberInput,
  SegmentedControl,
  Stack,
  Switch,
  Text,
} from "@mantine/core";
import classes from "./Preview.module.css";

type SettingsProps = {
  height: number;
  setHeight: React.Dispatch<React.SetStateAction<number>>;
  width: number;
  setWidth: React.Dispatch<React.SetStateAction<number>>;
  theme: "light" | "dark" | "auto";
  setTheme: React.Dispatch<React.SetStateAction<"light" | "dark" | "auto">>;
  primaryColor: string | undefined;
  setPrimaryColor: React.Dispatch<React.SetStateAction<string | undefined>>;
  textColor: string | undefined;
  setTextColor: React.Dispatch<React.SetStateAction<string | undefined>>;
  hoverColor: string | undefined;
  setHoverColor: React.Dispatch<React.SetStateAction<string | undefined>>;
  noteColor: string | undefined;
  setNoteColor: React.Dispatch<React.SetStateAction<string | undefined>>;
  autoProceed: boolean;
  setAutoProceed: React.Dispatch<React.SetStateAction<boolean>>;
  showStatusBar: boolean;
  setShowStatusBar: React.Dispatch<React.SetStateAction<boolean>>;
  showProgress: boolean;
  setShowProgress: React.Dispatch<React.SetStateAction<boolean>>;
  hideUploadButton: boolean;
  setHideUploadButton: React.Dispatch<React.SetStateAction<boolean>>;
  hideCancelButton: boolean;
  setHideCancelButton: React.Dispatch<React.SetStateAction<boolean>>;
  hideRetryButton: boolean;
  setHideRetryButton: React.Dispatch<React.SetStateAction<boolean>>;
  hidePauseResumeButton: boolean;
  setHidePauseResumeButton: React.Dispatch<React.SetStateAction<boolean>>;
  hideProgressAfterFinish: boolean;
  setHideProgressAfterFinish: React.Dispatch<React.SetStateAction<boolean>>;
  note: string;
  setNote: React.Dispatch<React.SetStateAction<string>>;
  singleFileFullScreen: boolean;
  setSingleFileFullScreen: React.Dispatch<React.SetStateAction<boolean>>;
  showSelectedFiles: boolean;
  setShowSelectedFiles: React.Dispatch<React.SetStateAction<boolean>>;
  showRemoveButtonAfterComplete: boolean;
  setShowRemoveButtonAfterComplete: React.Dispatch<
    React.SetStateAction<boolean>
  >;
};

const w = "300px";

const Settings: React.FC<SettingsProps> = ({
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
  autoProceed,
  setAutoProceed,
  showStatusBar,
  setShowStatusBar,
  showProgress,
  setShowProgress,
  hideUploadButton,
  setHideUploadButton,
  hideCancelButton,
  setHideCancelButton,
  hideRetryButton,
  setHideRetryButton,
  /* Will implement later - as not to complicate the code 
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
    setShowRemoveButtonAfterComplete
    */
}) => {
  return (
    <Stack align="space-between" h="100%" gap="xl">
      <Group justify="space-between" wrap="nowrap" gap="xl">
        <div>
          <Text size="sm">Height</Text>
          <Text c="dimmed">Set the height of the file uploader in px</Text>
        </div>
        <NumberInput
          w={w}
          placeholder="Enter height in px"
          value={height}
          onChange={(e) => setHeight(Number(e))}
        />
      </Group>
      <Group justify="space-between" wrap="nowrap" gap="xl">
        <div>
          <Text size="sm">Width</Text>
          <Text c="dimmed">Set the width of the file uploader in px</Text>
        </div>
        <NumberInput
          w={w}
          placeholder="Enter width in px"
          value={width}
          onChange={(e) => setWidth(Number(e))}
        />
      </Group>
      {/* Theme */}
      <Group justify="space-between" wrap="nowrap" gap="xl">
        <div>
          <Text size="sm">Choose Theme</Text>
          <Text c="dimmed">Set the theme of the file uploader</Text>
        </div>
        <SegmentedControl
          w={w}
          onChange={(value) => setTheme(value as "light" | "dark" | "auto")}
          value={theme}
          defaultValue="auto"
          data={[
            {
              value: "auto",
              label: "Auto",
            },
            {
              value: "dark",
              label: "Dark",
            },
            {
              value: "light",
              label: "Light",
            },
          ]}
        />
      </Group>

      {/* Primary Color */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Primary Color</Text>
          <Text c="dimmed">Change the primary color of the uploader</Text>
        </div>
        <ColorInput
          w={w}
          className={classes.color}
          value={primaryColor}
          onChange={setPrimaryColor}
        />
      </Group>

      {/* Primary Hover Color */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Primary Hover Color</Text>
          <Text c="dimmed">
            Change the color when the primary items are hovered
          </Text>
        </div>
        <ColorInput
          w={w}
          className={classes.color}
          value={hoverColor}
          onChange={setHoverColor}
        />
      </Group>

      {/* Primary Text Color */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Primary Text Color</Text>
          <Text c="dimmed">Change the primary text color</Text>
        </div>
        <ColorInput
          w={w}
          className={classes.color}
          value={textColor}
          onChange={setTextColor}
        />
      </Group>

      {/* Note Color */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Note Text Color</Text>
          <Text c="dimmed">Change the color of the note text</Text>
        </div>
        <ColorInput
          w={w}
          className={classes.color}
          value={noteColor}
          onChange={setNoteColor}
        />
      </Group>

      {/* Auto Proceed */}
      <Group justify="space-between" wrap="nowrap" gap="xl">
        <div>
          <Group align="center">
            <Text size="sm">Auto Proceed</Text>
            <Text c="red" opacity={0.9}>
              {hideUploadButton &&
                !autoProceed &&
                "* (If you hide the upload button, you must enable auto proceed)"}
            </Text>
          </Group>
          <Text c="dimmed">
            Toggle to start uploading file as soon as it is selected
          </Text>
        </div>
        <Switch
          className={classes.cusomSwitch}
          onLabel="ON"
          offLabel="OFF"
          checked={autoProceed}
          onChange={(e) => setAutoProceed(e.target.checked)}
        />
      </Group>

      {/* Status Bar */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Show status bar</Text>
          <Text c="dimmed">Toggle to show the status bar in the uploader</Text>
        </div>
        <Switch
          className={classes.cusomSwitch}
          onLabel="ON"
          offLabel="OFF"
          checked={showStatusBar}
          onChange={(e) => setShowStatusBar(e.target.checked)}
        />
      </Group>

      {/* Progress Bar */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Show progress bar</Text>
          <Text c="dimmed">
            Toggle to show the progress bar in the uploader
          </Text>
        </div>
        <Switch
          className={classes.cusomSwitch}
          onLabel="ON"
          offLabel="OFF"
          checked={showProgress}
          onChange={(e) => setShowProgress(e.target.checked)}
        />
      </Group>

      {/* Hide Upload Button */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Hide upload button</Text>
          <Text c="dimmed">
            Toggle to hide the upload button in the uploader
          </Text>
        </div>
        <Switch
          className={classes.cusomSwitch}
          onLabel="ON"
          offLabel="OFF"
          checked={hideUploadButton}
          onChange={(e) => setHideUploadButton(e.target.checked)}
        />
      </Group>

      {/* Hide Cancel Button */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Hide cancel button</Text>
          <Text c="dimmed">
            Toggle to hide the cancel button in the uploader
          </Text>
        </div>
        <Switch
          className={classes.cusomSwitch}
          onLabel="ON"
          offLabel="OFF"
          checked={hideCancelButton}
          onChange={(e) => setHideCancelButton(e.target.checked)}
        />
      </Group>

      {/* Hide Retry Button */}
      <Group
        justify="space-between"
        className={classes.item}
        wrap="nowrap"
        gap="xl"
      >
        <div>
          <Text size="sm">Hide retry button</Text>
          <Text c="dimmed">
            Toggle to hide the retry button in the uploader
          </Text>
        </div>
        <Switch
          className={classes.cusomSwitch}
          onLabel="ON"
          offLabel="OFF"
          checked={hideRetryButton}
          onChange={(e) => setHideRetryButton(e.target.checked)}
        />
      </Group>
    </Stack>
  );
};

export default Settings;
