import { Button, Card, Group, Badge, Avatar, Tooltip, Text } from '@mantine/core';
import { IconCode, IconPin, IconPinned } from '@tabler/icons-react';
import { Codebase } from '../../types/codebase';
import { useNavigate } from 'react-router-dom';
import classes from './CodebaseCard.module.css';

const avatars = [
    'https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-2.png',
    'https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-4.png',
    'https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-7.png',
];

const CodeBaseCard: React.FC<{ codebase: Codebase }> = ({ codebase }) => {
    const navigate = useNavigate();

    const handleClick = () => {
        navigate(`/codebases/${codebase.id}`);
    }

    return (
        <Card key={codebase.id} shadow="sm" padding="md" radius="md" withBorder>
            <Group justify='space-between' align='flex-start'>
                <div className={classes.cardcontent}>
                    <Group gap="xs" align='center' >
                        <IconCode size={20} />
                        <Text className={classes.cardtitle} onClick={handleClick}> {codebase.name}</Text>
                        <Badge variant='default' size="xs">{codebase.status}</Badge>
                    </Group>
                    <Text size="xs" c="dimmed" className={classes.carddescription}>{codebase.description}</Text>
                    <Group gap="xs" mt="xs">
                        {codebase.tags.map((tag) => (
                            <Badge key={tag} variant="light" size="xs">{tag}</Badge>
                        ))}
                    </Group>
                    <Text size="xs" c="dimmed" mt="xs">
                        Language: {codebase.lang} | Updated At: {new Date(codebase.updatedAt).toLocaleString('en-CA')}
                    </Text>
                </div>
                <div>
                    <Tooltip __size='xs' label={false ? 'Unpin' : 'Pin'}>
                        <Button
                            variant="light"
                            color={false ? 'blue' : 'gray'}
                        >
                            {false ? <IconPinned size={16} /> : <IconPin size={16} />}
                        </Button>
                    </Tooltip>
                </div>
            </Group>

            {/* Members Section */}
            <Avatar.Group spacing="sm" mt="sm">
                <Avatar src={avatars[0]} radius="xl" size={30} />
                <Avatar src={avatars[1]} radius="xl" size={30} />
                <Avatar src={avatars[2]} radius="xl" size={30} />
                <Avatar radius="xl" size={30}>+5</Avatar>
            </Avatar.Group>
        </Card>
    );
}

export default CodeBaseCard;