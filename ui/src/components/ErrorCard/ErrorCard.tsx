import {
  Box,
  Button,
  Divider,
  Group,
  Image,
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
    <Group align="center" justify="center" h={h}>
      <Image
        alt="error"
        h="300"
        w="400"
        src={colorScheme === 'dark' ? imageDark : imageLight}
      />
      <Box>
        <Title order={1} fw="bolder" p={0} m={0}>
          {getFormattedTitle(title, message)}
        </Title>
        <Divider mb="sm" size="lg" />
        <Text c="dimmed" size="md" p={0} m={0} mt="3">
          {getFormattedMessage(message)}
        </Text>
        <Button variant="outline" mt="sm" onClick={refreshPage}>
          Refresh the page
        </Button>
      </Box>
    </Group>
  );
};

const getFormattedTitle = (title: string, message: string) => {
  if (message === 'Request failed with status code 403') {
    return 'Unauthorized';
  }

  if (message === 'Network Error') {
    return 'Network Error';
  }

  return title.replace('AxiosError', 'Error');
};

const getFormattedMessage = (message: string) => {
  switch (message) {
    case 'Request failed with status code 403': {
      return 'Access Denied';
    }
    case 'Network Error': {
      return 'Please check your internet connection';
    }
  }
  return message;
};
