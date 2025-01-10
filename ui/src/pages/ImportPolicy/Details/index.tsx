import { Button, Container, Group, Loader, Paper, Title } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useNavigate, useParams } from 'react-router-dom';
import axiosInstance from '../../../utils/axios';
import { ImportPolicy } from '../../../types/importpolicy';
import ImportPolicyForm from '../New/Form';
import { useState } from 'react';
import { useGetImportPolicyDetails } from '../../../apis/importPolicy';

const ImportPolicyDetailsPage = () => {
    const navigate = useNavigate();
    const { importPolicyId } = useParams();
    const [editMode, setEditMode] = useState(false);

    const { isPending, error, importPolicy } = useGetImportPolicyDetails(importPolicyId || '');

    const handleSubmit = async (c: ImportPolicy) => {
        try {
            const resp = await axiosInstance.post("/importPolicies", c);
            if (resp.status === 200) {
                notifications.show({
                    title: 'Success',
                    message: 'Import policy added successfully',
                    color: 'green',
                })
                navigate('/importPolicies');
            }
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to add import policy',
                color: 'red',
            })
        }
    }

    const handleCancel = () => {
        setEditMode(false);
    }


    console.log(importPolicy)


    return (
        <Container p="sm">
            <Group justify='center' mb="lg">
                <Title order={3} mb="lg" opacity={0.8}>Import Policy : {importPolicy?.name}</Title>
            </Group>

            <Paper shadow='xs' p="lg" radius="md" withBorder>
                {
                    isPending ? <Loader /> : (importPolicy && <ImportPolicyForm onSubmit={handleSubmit} onCancel={handleCancel} showCancelButton showSubmitButton type={editMode ? 'edit' : 'view'} initialValues={importPolicy} />)
                }
                {!editMode && <Group justify="flex-end" mt="md" gap="md"><Button mt="md" onClick={() => setEditMode(true)}>Edit</Button></Group>}
            </Paper>
        </Container>
    );
}

export default ImportPolicyDetailsPage;