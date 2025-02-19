import {
  Button,
  Container,
  Image,
  Stack,
  Text,
  Title,
  useMantineColorScheme,
} from "@mantine/core";
import imageDark from "../../assets/images/error-dark.png";
import imageLight from "../../assets/images/error-light.png";
import classes from "./Error.module.css";

type ErrorCardProps = {
  title: string;
  message: string;
  h?: string;
};

export const ErrorCard: React.FC<ErrorCardProps> = ({ title, message, h }) => {
  const refreshPage = () => window.location.reload();
  const { colorScheme } = useMantineColorScheme();
  return (
    <Container className={classes.root}>
      <Stack align="center" h={h}>
        <Image
          h={300}
          w={500}
          src={colorScheme === "dark" ? imageDark : imageLight}
          className={classes.mobileImage}
        />
        <Title className={classes.title}>{title}</Title>
        <Text c="dimmed" size="lg">
          {message}
        </Text>
        <Button
          variant="outline"
          size="md"
          mt="xl"
          className={classes.control}
          onClick={refreshPage}
        >
          Refresh the page
        </Button>
      </Stack>
    </Container>
  );
};
