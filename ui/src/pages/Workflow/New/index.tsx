import { Container } from '@mantine/core';
import { CreateWorkflowForm, } from '../../../types/workflow';
import NewWorkflowForm from './Form';
import { useNavigate } from 'react-router-dom';
import { useCreateDataStoreMutation } from '../../../apis/storage';
import { useCreateWorkflowMutation } from '../../../apis/workflow';

const WorkflowsPage = () => {
    const navigate = useNavigate();

    const { mutateAsync } = useCreateDataStoreMutation();
    const { mutateAsync: mutateAsyncWorkflow } = useCreateWorkflowMutation();

    const handleSubmit = async (c: CreateWorkflowForm) => {
        mutateAsync({
            name: c.dataStoreName,
            connectorId: c.connectorId,
            bucket: c.bucket
        }).then((data) => {
            mutateAsyncWorkflow({
                name: c.name,
                description: c.description,
                tags: c.tags,
                importPolicyId: c.importPolicyId,
                dataStoreId: data.id
            }).then(() => {
                navigate('/workflows');
            })
        })
    }

    return (
        <Container p="sm" >
            <NewWorkflowForm onSubmit={handleSubmit} />
        </Container>
    );
}

export default WorkflowsPage;