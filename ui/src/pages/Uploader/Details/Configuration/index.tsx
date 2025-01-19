import { Box, Button, Group, Stack, Transition } from '@mantine/core';
import UploaderConfigForm from './Config';
import { CreateUploaderForm } from '../../../../types/uploader';
import { useForm } from '@mantine/form';
import { useUpdateUploaderConfigMutation } from '../../../../apis/uploader';
import { useParams } from 'react-router-dom';
import { IconCopy, IconDeviceFloppy, IconRestore } from '@tabler/icons-react';


const ConfigurationTab = ({ uploaderDetails }: { uploaderDetails: any }) => {

    const { uploaderId } = useParams();
    const { mutateAsync } = useUpdateUploaderConfigMutation(uploaderId as string)
    const form = useForm<CreateUploaderForm>({
        initialValues: { ...uploaderDetails.config, requiredMetadataFields: uploaderDetails.requiredMetadataFields || [] },
    });


    const handleEditAndSaveButton = async () => {
        if (form.isDirty()) {
            mutateAsync(form.values)
        }
    }

    const handleResetButton = () => {
        form.reset();
        form.resetDirty();
    }


    return (
        <Stack gap="md" justify='space-between'>
            <Box style={{ overflow: "auto" }} h="65vh" pr="xl">
                <UploaderConfigForm form={form} />
            </Box>
            <Transition
                mounted={form.isDirty()}
                transition="fade-up"
                duration={400}
                timingFunction="ease"
            >
                {(styles) => <div style={styles}>
                    <Group justify="center" gap="md" mt="xl">
                        <Button
                            variant="light"
                            onClick={handleResetButton}
                            leftSection={<IconRestore size={18} />}
                        >
                            Reset
                        </Button>
                        <Button
                            onClick={handleEditAndSaveButton}
                            leftSection={<IconDeviceFloppy size={18} />}
                        >
                            Save
                        </Button>
                    </Group>
                </div>}
            </Transition>
        </Stack>
    );
}

export default ConfigurationTab