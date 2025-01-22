import { IconCloud, IconSortDescendingSmallBig, IconTable, IconUser } from '@tabler/icons-react';
import {
    Paper,
    SimpleGrid,
    Text,
    useMantineTheme,
} from '@mantine/core';
import classes from './Hooks.module.css';

const mockdata = [
    {
        title: 'Copy to cloud storage',
        description:
            'Copy the imports to a cloud storage service like S3, GCS or Azure',
        icon: IconCloud,
    },
    {
        title: 'Authenticate users',
        description:
            'Add authentication to users before they can upload files to your workspace',
        icon: IconUser,
    },
    {
        title: 'Reduce file size',
        description: 'Compress files before they are uploaded to your workspace',
        icon: IconSortDescendingSmallBig,
    },
    {
        title: 'Validate CSV schema',
        description: 'Validate CSV files schema before they are uploaded to your workspace',
        icon: IconTable,
    },
];

export const HooksMarketPlace = () => {
    const theme = useMantineTheme();
    const features = mockdata.map((feature) => (
        <Paper key={feature.title} radius="md" p="xl">
            <feature.icon size={50} stroke={2} color={theme.colors.grape[6]} />
            <Text fz="lg" fw={500} className={classes.cardTitle} mt="md">
                {feature.title}
            </Text>
            <Text fz="sm" c="dimmed" mt="sm">
                {feature.description}
            </Text>
        </Paper>
    ));

    return (
        <SimpleGrid cols={{ base: 1, md: 4 }} spacing="xl">
            {features}
        </SimpleGrid>
    );
}