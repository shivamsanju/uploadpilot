import { Box, Group, NumberInput, SegmentedControl, Stack, Text } from '@mantine/core';
import UploaderConfigForm from '../../../New/Form/ConfigurationForm';
import { CreateUploaderForm } from '../../../../../types/uploader';
import { UseFormReturnType } from '@mantine/form';
import classes from "./Settings.module.css";

type SettingsProps = {
    height: number
    setHeight: React.Dispatch<React.SetStateAction<number>>
    width: number
    setWidth: React.Dispatch<React.SetStateAction<number>>
    theme: 'light' | 'dark' | 'auto'
    setTheme: React.Dispatch<React.SetStateAction<'light' | 'dark' | 'auto'>>
    form: UseFormReturnType<CreateUploaderForm, (values: CreateUploaderForm) => CreateUploaderForm>;
    editMode: boolean
}

const w = "15vw";

const Settings: React.FC<SettingsProps> = ({
    height,
    setHeight,
    width,
    setWidth,
    theme,
    setTheme,
    form,
    editMode
}) => {
    return (
        <Box style={{ overflow: "auto" }} pr="xl">
            <Stack px="md" pt="sm">
                <Group justify="space-between" pb="sm" className={classes.settingsItem} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Height</Text>
                        <Text size="xs" c="dimmed">
                            Set the height of the file uploader in px
                        </Text>
                    </div>
                    <NumberInput
                        w={w}
                        size="xs"
                        placeholder="Enter height in px"
                        value={height}
                        onChange={(e) => setHeight(Number(e))}
                    />
                </Group>
                <Group justify="space-between" pb="sm" className={classes.settingsItem} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Width</Text>
                        <Text size="xs" c="dimmed">
                            Set the width of the file uploader in px
                        </Text>
                    </div>
                    <NumberInput
                        w={w}
                        size="xs"
                        placeholder="Enter width in px"
                        value={width}
                        onChange={(e) => setWidth(Number(e))}
                    />
                </Group>
                {/* Theme */}
                <Group justify="space-between" pb="sm" className={classes.settingsItem} wrap="nowrap" gap="xl">
                    <div>
                        <Text size="sm">Choose Theme</Text>
                        <Text size="xs" c="dimmed">
                            Set the theme of the file uploader
                        </Text>
                    </div>
                    <SegmentedControl
                        w={w}
                        size="xs"
                        onChange={(value) => setTheme(value as 'light' | 'dark' | 'auto')}
                        value={theme}
                        data={[
                            {
                                value: 'auto',
                                label: 'Auto',
                            },
                            {
                                value: 'dark',
                                label: 'Dark',
                            },
                            {
                                value: 'light',
                                label: 'Light',
                            },
                        ]}
                    />
                </Group>
            </Stack>
            <UploaderConfigForm form={form} type={editMode ? "edit" : "view"} />
        </Box>
    );
}

export default Settings