import {
  IconArrowDownRight,
  IconArrowUpRight,
  IconCircleCheck,
  IconExclamationCircle,
  IconNumber,
  IconServer2,
} from "@tabler/icons-react";
import { Box, Group, Paper, SimpleGrid, Text, Title } from "@mantine/core";
import classes from "./dashboard.module.css";
import { LineChart } from "@mantine/charts";

const icons = {
  total: IconNumber,
  success: IconCircleCheck,
  fail: IconExclamationCircle,
  file: IconServer2,
};

const data = [
  { title: "Total Uploads", icon: "total", value: "40,215", diff: 14 },
  { title: "Successful", icon: "success", value: "40,145", diff: 13 },
  { title: "Failed", icon: "fail", value: "70", diff: -18 },
  { title: "Total size", icon: "file", value: "75 GB", diff: 4 },
] as const;

const data2 = [
  { date: "Jan", uploads: 25000 },
  { date: "Feb", uploads: 10000 },
  { date: "Mar", uploads: 5000 },
  { date: "Apr", uploads: 15000 },
  { date: "May", uploads: 30000 },
  { date: "Jun", uploads: 15000 },
  { date: "Jul", uploads: 30000 },
  { date: "Aug", uploads: 40000 },
  { date: "Sep", uploads: 15000 },
  { date: "Oct", uploads: 20000 },
  { date: "Nov", uploads: 10000 },
  { date: "Dec", uploads: 10000 },
];

const DashboardPage = () => {
  const stats = data.map((stat) => {
    const Icon = icons[stat.icon];
    const DiffIcon = stat.diff > 0 ? IconArrowUpRight : IconArrowDownRight;

    return (
      <Paper withBorder p="md" radius="md" key={stat.title}>
        <Group justify="space-between">
          <Text c="dimmed" className={classes.title}>
            {stat.title}
          </Text>
          <Icon className={classes.icon} size={22} stroke={1.5} />
        </Group>

        <Group align="flex-end" gap="xs" mt={25}>
          <Text className={classes.value}>{stat.value}</Text>
          <Text
            c={stat.diff > 0 ? "green" : "red"}
            fz="sm"
            fw={500}
            className={classes.diff}
          >
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
    <Box mb={50}>
      <Title order={3} opacity={0.7} mb="xl">
        Analytics
      </Title>
      {/* <Text c="dimmed" size="xs" mt={2} mb="lg">
                Analyze file upload trends over the last 6 months
            </Text> */}
      <SimpleGrid cols={{ base: 1, xs: 2, md: 4 }}>{stats}</SimpleGrid>
      <Paper mt="xl" p="sm">
        <LineChart
          title="File upload trends"
          h={400}
          data={data2}
          series={[{ name: "uploads", label: "Number of uploads" }]}
          dataKey="date"
          type="gradient"
          gradientStops={[
            { offset: 0, color: "red.6" },
            { offset: 20, color: "appcolor.6" },
            { offset: 40, color: "yellow.5" },
            { offset: 70, color: "lime.5" },
            { offset: 80, color: "cyan.5" },
            { offset: 100, color: "appcolor.5" },
          ]}
          strokeWidth={5}
          curveType="natural"
          yAxisProps={{ domain: [0, 60000] }}
          valueFormatter={(value) => `${value / 1000}K`}
        />
      </Paper>
    </Box>
  );
};

export default DashboardPage;
