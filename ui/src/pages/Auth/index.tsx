import {
  Button,
  Image,
  Stack,
  Text,
  Modal,
  Anchor,
  MantineProvider,
  ScrollArea,
  Box,
  Title,
} from "@mantine/core";
import classes from "./Auth.module.css";
import { getApiDomain } from "../../utils/config";
import TokenCatcher from "./TokenCatcher";
import { useEffect, useState } from "react";
import { AppLoader } from "../../components/Loader/AppLoader";
import { useNavigate } from "react-router-dom";
import { Logo2 } from "../../components/Logo/Logo2";
import GoogleIcon from "../../assets/icons/google.svg";
import GithubIcon from "../../assets/icons/github.svg";
import PrivacyPolicy from "../../components/Policy/PrivacyPolicy";
import TermsOfService from "../../components/Policy/TermsOfService";
import { Lines } from "../../components/Lines";

const AuthPage = () => {
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<"terms" | "privacy">(
    "terms"
  );
  const navigate = useNavigate();

  const handleLogin = (provider: string) => {
    window.location.href = getApiDomain() + `/auth/${provider}/authorize`;
  };

  useEffect(() => {
    const token = localStorage.getItem("uploadpilottoken");
    if (token) {
      navigate("/", { replace: true });
    }
    setLoading(false);
  }, [navigate]);

  const openModal = (content: "terms" | "privacy") => {
    setModalContent(content);
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
  };

  if (loading) {
    return <AppLoader h="100vh" />;
  }

  return (
    <div className={classes.wrapper}>
      <Lines />
      <TokenCatcher />
      <MantineProvider forceColorScheme="dark">
        <Box w={{ base: "100vw", sm: 600 }} className={classes.form}>
          <Title order={1} c="dimmed" ta="center" mb="30">
            Welcome to
          </Title>
          <Stack gap="xs" mb="60">
            <Logo2 enableOnClick={false} />
          </Stack>

          <Stack mt="xl">
            <Text size="xs" ta="center" c="dimmed">
              By continuing, you agree to our{" "}
              <Anchor
                onClick={() => openModal("terms")}
                size="sm"
                style={{ cursor: "pointer" }}
              >
                Terms of Service
              </Anchor>{" "}
              and acknowledge you have read our{" "}
              <Anchor
                onClick={() => openModal("privacy")}
                size="sm"
                style={{ cursor: "pointer" }}
              >
                Privacy Policy
              </Anchor>
            </Text>
            <Button
              c="dark"
              variant="white"
              leftSection={<Image src={GoogleIcon} width={20} height={20} />}
              onClick={() => handleLogin("google")}
              size="sm"
            >
              Google
            </Button>
            <Text ta="center" c="dimmed">
              or
            </Text>
            <Button
              c="dark"
              variant="white"
              leftSection={<Image src={GithubIcon} width={25} height={25} />}
              onClick={() => handleLogin("github")}
              size="sm"
            >
              Github
            </Button>
          </Stack>
        </Box>

        {/* Modal for Terms of Service or Privacy Policy */}
        <Modal
          opened={modalOpen}
          onClose={closeModal}
          title={
            modalContent === "terms" ? "Terms of Service" : "Privacy Policy"
          }
          size="xl"
          scrollAreaComponent={(props) => (
            <ScrollArea.Autosize {...props} scrollbarSize={10} />
          )}
        >
          {modalContent === "terms" ? <TermsOfService /> : <PrivacyPolicy />}
        </Modal>
      </MantineProvider>
    </div>
  );
};

export default AuthPage;
