import {
  ColorInput,
  Grid,
  NumberInput,
  SegmentedControl,
  SimpleGrid,
  Switch,
  Text,
} from '@mantine/core';
import classes from './Preview.module.css';
const w = '150px';

type SettingsProps = {
  height: number;
  setHeight: React.Dispatch<React.SetStateAction<number>>;
  width: number;
  setWidth: React.Dispatch<React.SetStateAction<number>>;
  theme: 'light' | 'dark' | 'auto';
  setTheme: React.Dispatch<React.SetStateAction<'light' | 'dark' | 'auto'>>;
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
}) => {
  const items: any = [
    {
      label: 'Height',
      description: 'Set the height of the file uploader in px',
      value: height,
      setter: setHeight,
      type: 'number',
    },
    {
      label: 'Theme',
      description: 'Change the theme',
      value: theme,
      setter: setTheme,
      type: 'segmented',
    },
    {
      label: 'Width',
      description: 'Set the width of the file uploader in px',
      value: width,
      setter: setWidth,
      type: 'number',
    },
    {
      label: 'Auto Proceed',
      description: 'Start uploading automatically',
      value: autoProceed,
      setter: setAutoProceed,
      type: 'switch',
    },
    {
      label: 'Primary Color',
      description: 'Change the primary color of the uploader',
      value: primaryColor,
      setter: setPrimaryColor,
      type: 'color',
    },
    {
      label: 'Show Status Bar',
      description: 'Toggle to show the status bar',
      value: showStatusBar,
      setter: setShowStatusBar,
      type: 'switch',
    },
    {
      label: 'Primary Hover Color',
      description: 'Change the color when hovered',
      value: hoverColor,
      setter: setHoverColor,
      type: 'color',
    },
    {
      label: 'Show Progress Bar',
      description: 'Toggle to show progress',
      value: showProgress,
      setter: setShowProgress,
      type: 'switch',
    },
    {
      label: 'Primary Text Color',
      description: 'Change the primary text color',
      value: textColor,
      setter: setTextColor,
      type: 'color',
    },
    {
      label: 'Hide Upload Button',
      description: 'Toggle to hide upload button',
      value: hideUploadButton,
      setter: setHideUploadButton,
      type: 'switch',
    },
    {
      label: 'Note Text Color',
      description: 'Change the note text color',
      value: noteColor,
      setter: setNoteColor,
      type: 'color',
    },
    {
      label: 'Hide Cancel Button',
      description: 'Toggle to hide cancel button',
      value: hideCancelButton,
      setter: setHideCancelButton,
      type: 'switch',
    },
  ];

  return (
    <SimpleGrid cols={{ sm: 1, lg: 2 }} spacing="xl">
      {items.map(
        ({ label, description, value, setter, type }: any, index: number) => (
          <Grid key={index}>
            <Grid.Col span={6}>
              <Text size="sm">{label}</Text>
              <Text c="dimmed">{description}</Text>
            </Grid.Col>
            <Grid.Col span={6} className={classes.inputContainer}>
              {type === 'number' && (
                <NumberInput
                  w={w}
                  value={value}
                  onChange={e => setter(Number(e))}
                />
              )}
              {type === 'color' && (
                <ColorInput w={w} value={value} onChange={setter} />
              )}
              {type === 'switch' && (
                <Switch
                  className={classes.customSwitch}
                  checked={value}
                  onChange={e => setter(e.target.checked)}
                />
              )}
              {type === 'segmented' && (
                <SegmentedControl
                  w={250}
                  h="34"
                  size="xs"
                  withItemsBorders={false}
                  value={value}
                  onChange={setter}
                  data={[
                    { label: 'Auto', value: 'auto' },
                    { label: 'Light', value: 'light' },
                    { label: 'Dark', value: 'dark' },
                  ]}
                />
              )}
            </Grid.Col>
          </Grid>
        ),
      )}
    </SimpleGrid>
  );
};

export default Settings;
