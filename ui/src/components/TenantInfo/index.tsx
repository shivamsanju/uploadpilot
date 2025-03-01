import { Avatar, Group, Paper, Text } from '@mantine/core';
import classes from './TenantInfo.module.css';
type Props = {
  tenantId: string;
  tenantName: string;
};

export const TenantInfoCard: React.FC<Props> = ({ tenantId, tenantName }) => {
  return (
    <Paper withBorder radius="md" p="sm" className={classes.tenantInfoCard}>
      <Group wrap="nowrap" justify="space-between">
        <Group wrap="nowrap">
          <Avatar
            size="40"
            src="https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-2.png"
            radius="md"
          />
          <div>
            <Text fz="lg" fw={700} c="dimmed">
              {tenantName}
            </Text>

            <Text fz="xs" c="dimmed">
              Free Plan
            </Text>
          </div>
        </Group>
      </Group>
    </Paper>
  );
};
