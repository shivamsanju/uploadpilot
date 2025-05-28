import { ActionIcon, Avatar, Badge, Group, Paper, Text } from '@mantine/core';
import { IconArrowRightDashed } from '@tabler/icons-react';
import { getRandomAvatarName } from '../../utils/avatar';
import classes from './TenantInfo.module.css';
type Props = {
  tenantId: string;
  tenantName: string;
  onChoose: (id: string) => void;
};

export const TenantInfoCard: React.FC<Props> = ({
  tenantId,
  tenantName,
  onChoose,
}) => {
  return (
    <Paper withBorder radius="md" p="sm" className={classes.tenantInfoCard}>
      <Group wrap="nowrap" justify="space-between">
        <Group wrap="nowrap">
          <Avatar size="40" src={getRandomAvatarName(tenantId)} radius="md" />
          <div>
            <Text fz="lg" fw={700}>
              {tenantName}
            </Text>

            <Badge size="xs" variant="light" color="blue">
              Free Plan
            </Badge>
          </div>
        </Group>
        <ActionIcon variant="filled" onClick={() => onChoose(tenantId)}>
          <IconArrowRightDashed size={18} stroke={1.5} />
        </ActionIcon>
      </Group>
    </Paper>
  );
};
