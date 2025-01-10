import { Card, Title } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import NewConnectorForm from './Form';
import { useCreateStorageConnectorMutation } from '../../../../apis/storage';

const NewConnectorsPage = () => {
    const navigate = useNavigate();
    const { mutateAsync } = useCreateStorageConnectorMutation()

    const handleSubmit = async (c: any) => {
        const body: any = {
            name: c.name,
            type: c.type,
            tags: c.tags
        }
        if (c.type === 's3') {
            body.s3Config = {
                region: c.s3Region,
                accessKey: c.s3AccessKey,
                secretKey: c.s3SecretKey
            }
        } else if (c.type === 'gcs') {
            body.gcsConfig = {
                apiKey: c.gcsApiKey
            }
        } else if (c.type === 'azure') {
            body.azureConfig = {
                accountName: c.azureAccountName,
                accountKey: c.azureAccountKey
            }
        }

        mutateAsync(body).then(() => {
            navigate('/storage/connectors');
        })

    }

    const handleCancel = () => {
        navigate('/storage/connectors');
    }

    return (
        <>
            <Title order={3} opacity={0.8} mt={8} mb={20}>New Connector</Title>
            <Card w="100%" withBorder shadow="xs" radius="md" padding="xl">
                <NewConnectorForm showSubmitButton onSubmit={handleSubmit} showCancelButton onCancel={handleCancel} />
            </Card>
        </>
    );
}

export default NewConnectorsPage;