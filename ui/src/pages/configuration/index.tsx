import { Box, Title } from "@mantine/core";
import UploaderConfigForm from "./Config";
import { useGetUploaderConfig } from "../../apis/uploader";
import { useParams } from "react-router-dom";
import { ErrorLoadingWrapper } from "../../components/ErrorLoadingWrapper";

const ConfigurationPage = () => {
  const { workspaceId } = useParams();

  let { isPending, error, config } = useGetUploaderConfig(
    workspaceId as string,
  );

  if (!isPending && !error && !config) {
    error = new Error("No config found for this workspace");
  }

  return (
    <ErrorLoadingWrapper error={error} isPending={isPending}>
      <Title order={3} opacity={0.7} mb="xl">
        Configuration
      </Title>
      {/* <Text c="dimmed" size="xs" mt={2} mb="lg">
                Configure your uploader to match your requirements
            </Text> */}
      <Box mb={50}>
        <UploaderConfigForm config={config} />
      </Box>
    </ErrorLoadingWrapper>
  );
};

export default ConfigurationPage;
