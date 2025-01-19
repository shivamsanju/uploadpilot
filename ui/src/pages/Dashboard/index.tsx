import {
    IconArrowDownRight,
    IconArrowUpRight,
    IconCoin,
    IconDiscount2,
    IconReceipt2,
    IconUserPlus,
} from '@tabler/icons-react';
import { Box, Card, Group, Paper, SimpleGrid, Text, Title } from '@mantine/core';
import classes from './dashboard.module.css';
import { LineChart } from '@mantine/charts';

const icons = {
    user: IconUserPlus,
    discount: IconDiscount2,
    receipt: IconReceipt2,
    coin: IconCoin,
};

const data = [
    { title: 'Revenue', icon: 'receipt', value: '13,456', diff: 34 },
    { title: 'Profit', icon: 'coin', value: '4,145', diff: -13 },
    { title: 'Coupons usage', icon: 'discount', value: '745', diff: 18 },
    { title: 'New customers', icon: 'user', value: '188', diff: -30 },
] as const;

const data2 = [
    { date: 'Jan', temperature: -25 },
    { date: 'Feb', temperature: -10 },
    { date: 'Mar', temperature: 5 },
    { date: 'Apr', temperature: 15 },
    { date: 'May', temperature: 30 },
    { date: 'Jun', temperature: 15 },
    { date: 'Jul', temperature: 30 },
    { date: 'Aug', temperature: 40 },
    { date: 'Sep', temperature: 15 },
    { date: 'Oct', temperature: 20 },
    { date: 'Nov', temperature: 0 },
    { date: 'Dec', temperature: -10 },
];

const DashboardPage = () => {
    const stats = data.map((stat) => {
        const Icon = icons[stat.icon];
        const DiffIcon = stat.diff > 0 ? IconArrowUpRight : IconArrowDownRight;

        return (
            <Paper withBorder p="md" radius="md" key={stat.title}>

                <Group justify="space-between">
                    <Text size="xs" c="dimmed" className={classes.title}>
                        {stat.title}
                    </Text>
                    <Icon className={classes.icon} size={22} stroke={1.5} />
                </Group>

                <Group align="flex-end" gap="xs" mt={25}>
                    <Text className={classes.value}>{stat.value}</Text>
                    <Text c={stat.diff > 0 ? 'grape' : 'red'} fz="sm" fw={500} className={classes.diff}>
                        <span>{stat.diff}%</span>
                        <DiffIcon size={16} stroke={1.5} />
                    </Text>
                </Group>

                <Text fz="xs" c="dimmed" mt={7}>
                    Compared to previous month
                </Text>
            </Paper>
        );
    });
    return (
        <Box px="sm">
            <Title order={3} mb="lg" opacity={0.7}>Dashboard</Title>
            <SimpleGrid cols={{ base: 1, xs: 2, md: 4 }}>{stats}</SimpleGrid>
            <Card mt="xl">
                <LineChart
                    title='File upload trends'
                    h={500}
                    data={data2}
                    series={[{ name: 'temperature', label: 'Avg. Temperature' }]}
                    dataKey="date"
                    type="gradient"
                    gradientStops={[
                        { offset: 0, color: 'red.6' },
                        { offset: 20, color: 'orange.6' },
                        { offset: 40, color: 'yellow.5' },
                        { offset: 70, color: 'lime.5' },
                        { offset: 80, color: 'cyan.5' },
                        { offset: 100, color: 'grape.5' },
                    ]}
                    strokeWidth={5}
                    curveType="natural"
                    yAxisProps={{ domain: [-25, 60] }}
                    valueFormatter={(value) => `${value}°C`}
                />
            </Card>
        </Box>
    );
}

export default DashboardPage