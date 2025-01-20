import { Box, Stack } from '@mantine/core';
import UploaderConfigForm from './Config';
import { useGetUploaderConfig } from '../../../apis/uploader';
import { useParams } from 'react-router-dom';
import { ErrorLoadingWrapper } from '../../../components/ErrorLoadingWrapper';


const ConfigurationPage = () => {
    const { workspaceId } = useParams();

    let { isPending, error, config } = useGetUploaderConfig(workspaceId as string);

    if (!isPending && !error && !config) {
        error = new Error("No config found for this workspace");
    }

    return (
        <ErrorLoadingWrapper error={error} isPending={isPending} >
            <Stack gap="md" justify='space-between'>
                <Box mb="xl">
                    <UploaderConfigForm config={config} />
                </Box>

            </Stack>
        </ErrorLoadingWrapper>
    );
}

export default ConfigurationPage