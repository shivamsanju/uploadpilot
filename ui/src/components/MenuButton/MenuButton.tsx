import {
  Group,
  ThemeIcon,
  Title,
  Transition,
  UnstyledButton,
} from '@mantine/core';
import { IconMenu4 } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useNavbar } from '../../context/NavbarContext';

type Props = {};

export const MenuButton: React.FC<Props> = () => {
  const { toggle, opened } = useNavbar();
  const navigate = useNavigate();
  return (
    <Group align="center" gap="sm" justify="center" wrap="nowrap">
      <ThemeIcon
        onClick={toggle}
        variant="subtle"
        style={{ cursor: 'pointer' }}
        opacity={0.9}
      >
        <IconMenu4 size={22} />
      </ThemeIcon>
      <Transition
        mounted={opened}
        transition="fade"
        duration={200}
        timingFunction="ease"
      >
        {styles =>
          opened ? (
            <UnstyledButton onClick={() => navigate('/')}>
              <Title
                order={3}
                style={{ ...styles }}
                opacity={0.9}
                lineClamp={1}
              >
                Upload Pilot
              </Title>
            </UnstyledButton>
          ) : (
            <></>
          )
        }
      </Transition>
    </Group>
  );
};
