import { Box, Group, NumberInput, Stack, Text, TextInput } from '@mantine/core';
import UploaderConfigForm from '../../../New/Form/ConfigurationForm';
import { CreateUploaderForm } from '../../../../../types/uploader';
import { useForm } from '@mantine/form';
import classes from "./Settings.module.css";

type SettingsProps = {
    height: number
    setHeight: React.Dispatch<React.SetStateAction<number>>
    width: number
    setWidth: React.Dispatch<React.SetStateAction<number>>
    uploaderDetails: any
}

const w = "15vw";

const Settings: React.FC<SettingsProps> = ({
    height,
    setHeight,
    width,
    setWidth,
    uploaderDetails
}) => {
    const form = useForm<CreateUploaderForm>({
        initialValues: { ...uploaderDetails.config, requiredMetadataFields: uploaderDetails.requiredMetadataFields || [] },
    });

    return (
        <Box style={{ overflow: "auto", height: '62vh' }}>
            <Stack p="md" pt="xs">
                <Group justify="space-between" mb="md" pb="xl" className={classes.settingsItem} wrap="nowrap" gap="xl">
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
                <Group justify="space-between" mb="md" pb="xl" className={classes.settingsItem} wrap="nowrap" gap="xl">
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
            </Stack>
            <UploaderConfigForm form={form} type="view" />
        </Box>
    );
}

export default Settings