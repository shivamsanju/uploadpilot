import { Group, ThemeIcon, Title, Transition } from '@mantine/core';
import { IconMenu4 } from '@tabler/icons-react';

type Props = {
  toggle: React.Dispatch<React.SetStateAction<boolean>>;
  collapsed: boolean;
};

export const MenuButton: React.FC<Props> = ({ toggle, collapsed }) => {
  return (
    <Group align="center" gap="sm" justify="center" wrap="nowrap">
      <ThemeIcon
        onClick={() => toggle(prev => !prev)}
        variant="subtle"
        style={{ cursor: 'pointer' }}
        opacity={0.9}
      >
        <IconMenu4 size={22} />
      </ThemeIcon>
      <Transition
        mounted={!collapsed}
        transition="fade"
        duration={200}
        timingFunction="ease"
      >
        {styles =>
          collapsed ? (
            <></>
          ) : (
            <Title order={3} style={{ ...styles }} opacity={0.9} lineClamp={1}>
              Upload Pilot
            </Title>
          )
        }
      </Transition>
    </Group>
  );
};
