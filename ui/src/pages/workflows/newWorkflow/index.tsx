import { Container, Title } from '@mantine/core';
import { Workflow } from '../../../types/workflow';
import NewWorkflowForm from '../../../components/NewWorkflowForm';
import axios from 'axios';
import { notifications } from '@mantine/notifications';
import { getApiDomain } from '../../../config';
import { useNavigate } from 'react-router-dom';

const WorkflowsPage = () => {
    const navigate = useNavigate();

    const handleSubmit = async (c: Workflow) => {
        try {
            const resp = await axios.post(getApiDomain() + "/workflows", c);
            if (resp.status === 200) {
                notifications.show({
                    title: 'Success',
                    message: 'Workflow added successfully',
                    color: 'green',
                })
                navigate('/');
            }
        } catch {
            notifications.show({
                title: 'Error',
                message: 'Failed to add workflow',
                color: 'red',
            })
        }
    }

    return (
        <Container size="lg">
            <Title opacity={0.8} mt={8} mb={20}>New Workflow</Title>
            <NewWorkflowForm onSubmit={handleSubmit} />
        </Container>
    );
}

export default WorkflowsPage;