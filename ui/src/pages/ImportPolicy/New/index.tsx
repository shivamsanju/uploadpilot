import { Container, Group, Paper, Title } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import ImportPolicyForm from './Form';
import { useCreateImportPolicyMutation } from '../../../apis/importPolicy';
import { ImportPolicy } from '../../../types/importpolicy';

const NewImportPolicyPage = () => {
    const navigate = useNavigate();

    const { mutateAsync } = useCreateImportPolicyMutation();

    const handleSubmit = async (c: ImportPolicy) => {
        mutateAsync(c).then(() => {
            navigate('/importPolicies');
        })
    }

    const handleCancel = () => {
        navigate('/importPolicies');
    }

    return (
        <Container p="sm">
            <Group justify='center' mb="lg">
                <Title order={3} mb="lg" opacity={0.8}>Create new import policy</Title>
            </Group>

            <Paper shadow='xs' p="lg" radius="md" withBorder>
                <ImportPolicyForm onSubmit={handleSubmit} onCancel={handleCancel} showCancelButton showSubmitButton />
            </Paper>
        </Container>
    );
}

export default NewImportPolicyPage;