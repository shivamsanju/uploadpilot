import { Card, Title } from '@mantine/core';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../../config';
import { useNavigate } from 'react-router-dom';
import NewConnectorForm, { Connector } from '../../../components/NewConnectorForm';

const NewConnectorsPage = () => {
    const navigate = useNavigate();

    const handleSubmit = async (c: Connector) => {
        try {
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
            const resp = await axios.post(getApiDomain() + "/storage/connectors", body);
            if (resp.status === 200) {
                notifications.show({
                    title: 'Success',
                    message: 'Connector added successfully',
                    color: 'green',
                })
                navigate('/storage/connectors');
            } else {
                notifications.show({
                    title: 'Error',
                    message: 'Failed to add Connector',
                    color: 'red',
                })
            }
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to add Connector',
                color: 'red',
            })
        }
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