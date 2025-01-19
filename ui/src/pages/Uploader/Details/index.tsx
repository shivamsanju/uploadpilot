import { useParams } from 'react-router-dom';
import { useGetUploaderDetailsById } from '../../../apis/uploader';
import { Badge, Group, Box, Text, Title } from '@mantine/core';
import AppLoader from '../../../components/Loader/AppLoader';
import ErrorCard from '../../../components/ErrorCard/ErrorCard';
import { timeAgo } from '../../../utils/datetime';
import UploaderTabs from './Tabs/Tabs';

const CodeMapPage = () => {
    const { uploaderId } = useParams();

    const { isPending, error, uploader } = useGetUploaderDetailsById(uploaderId as string);

    return isPending ? <AppLoader h="70vh" /> : error ? <ErrorCard title={error.name} message={error.message} h="70vh" /> : (
        <Box h={"90vh"}>
            <Group justify='space-between'>
                <Title order={3} opacity={0.7}>{uploader.name}</Title>
                <Group align='center' >
                    {uploader.tags && uploader.tags.length > 0 && uploader.tags.map((tag: string) => (<Badge size="sm" variant="light" ta="center" >{tag}</Badge>))}
                </Group>
            </Group>
            <Group align='center' mb="sm">
                <Text size="xs" opacity={0.5}>Created By: {uploader.createdBy}</Text>
                <Text size="xs" opacity={0.5}>Updated: {timeAgo.format(new Date(uploader.updatedAt))}</Text>
            </Group>
            <UploaderTabs uploaderDetails={uploader} />
        </Box>
    );
}

export default CodeMapPage;