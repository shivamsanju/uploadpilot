import { useParams } from 'react-router-dom';
import { useGetUploaderDetailsById } from '../../../apis/uploader';
import { Badge, Group, Paper, Text, Title } from '@mantine/core';
import AppLoader from '../../../components/Loader/AppLoader';
import ErrorCard from '../../../components/ErrorCard/ErrorCard';
import { timeAgo } from '../../../utils/datetime';
import UploaderTabs from './Tabs/Tabs';

const CodeMapPage = () => {
    const { uploaderId } = useParams();

    const { isPending, error, uploader } = useGetUploaderDetailsById(uploaderId as string);

    return isPending ? <AppLoader h="70vh" /> : error ? <ErrorCard title={error.name} message={error.message} h="70vh" /> : (
        <Paper shadow="xs" p="sm" radius="xs" withBorder h={"90vh"}>
            <Group justify='space-between'>
                <Title order={3} mb="xs" opacity={0.8}>{uploader.name}</Title>
                <Group align='center' >
                    {uploader.tags && uploader.tags.length > 0 && uploader.tags.map((tag: string) => (<Badge size="sm" variant="light" ta="center" >{tag}</Badge>))}
                </Group>
            </Group>
            <Group align='center' mb="sm">
                <Text size="xs">Created By: {uploader.createdBy}</Text>
                <Text size="xs">Updated: {timeAgo.format(new Date(uploader.updatedAt))}</Text>
            </Group>
            <Group align='center' mb="sm">
                <Text lineClamp={2} size="xs" opacity={0.8} >{uploader.description}</Text>
            </Group>
            <UploaderTabs uploaderDetails={uploader} />
        </Paper>
    );
}

export default CodeMapPage;