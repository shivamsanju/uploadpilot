import {
  Button,
  Container,
  Image,
  Stack,
  Text,
  Title,
  useMantineColorScheme,
} from '@mantine/core';
import imageDark from '../../assets/images/error-dark.png';
import imageLight from '../../assets/images/error-light.png';

type ErrorCardProps = {
  title: string;
  message: string;
  h?: string;
};

export const ErrorCard: React.FC<ErrorCardProps> = ({ title, message, h }) => {
  const refreshPage = () => window.location.reload();
  const { colorScheme } = useMantineColorScheme();
  return (
    <Container p="sm">
      <Stack align="center" h={h}>
        <Image
          alt="error"
          h="300"
          w="400"
          src={colorScheme === 'dark' ? imageDark : imageLight}
        />
        <Title order={1} fw="bolder" p={0}>
          {title}
        </Title>
        <Text c="dimmed" size="md">
          {message}
        </Text>
        <Button variant="outline" size="md" mt="sm" onClick={refreshPage}>
          Refresh the page
        </Button>
      </Stack>
    </Container>
  );
};
