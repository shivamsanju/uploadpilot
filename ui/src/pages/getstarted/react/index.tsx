import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import {
  ActionIcon,
  Box,
  Button,
  Grid,
  Group,
  MantineStyleProp,
  Paper,
  PasswordInput,
  Stack,
  Text,
  Timeline,
  Tooltip,
} from '@mantine/core';
import {
  IconBrandNpm,
  IconCode,
  IconConfetti,
  IconEditCircle,
  IconKey,
} from '@tabler/icons-react';
import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Uploader } from 'uppy-react';
import 'uppy-react/dist/style.css';
import { useGetUserDetails } from '../../../apis/user';
import { AppLoader } from '../../../components/Loader/AppLoader';
import { TEMP_API_KEY, TENANT_ID_KEY } from '../../../constants/tenancy';
import { getUploadApiDomain } from '../../../utils/config';
import Settings from './Settings';
import { useSettingsProps } from './SettingsProps';

const getCode = (
  workspaceId: string,
  backendEndpoint: string,
  settingsProps: any,
) => {
  const properties = [
    `workspaceId="${workspaceId}"`,
    `uploadEndpoint="${backendEndpoint}"`,
    `height={${settingsProps.height}}`,
    `width={${settingsProps.width}}`,
    `theme={${settingsProps.theme}}`,
    settingsProps.primaryColor &&
      `primaryColor="${settingsProps.primaryColor}"`,
    settingsProps.textColor && `textColor="${settingsProps.textColor}"`,
    settingsProps.hoverColor && `hoverColor="${settingsProps.hoverColor}"`,
    settingsProps.noteColor && `noteColor="${settingsProps.noteColor}"`,
    `metadata={{"key": "value"}}`,
    `headers={{"key": "value"}}`,
    !settingsProps.showStatusBar && 'showStatusBar={false}',
    !settingsProps.showProgress && 'showProgress={false}',
    settingsProps.hideUploadButton && 'hideUploadButton={true}',
    settingsProps.hideCancelButton && 'hideCancelButton={true}',
    settingsProps.hideRetryButton && 'hideRetryButton={true}',
    settingsProps.hidePauseResumeButton && 'hidePauseResumeButton={true}',
    settingsProps.hideProgressAfterFinish && 'hideProgressAfterFinish={true}',
    settingsProps.note && `note="${settingsProps.note}"`,
    !settingsProps.singleFileFullScreen && 'singleFileFullScreen={false}',
    !settingsProps.showSelectedFiles && 'showSelectedFiles={false}',
  ]
    .filter(Boolean) // Remove any falsy values
    .join('\n            '); // Join with proper indentation

  const code = `
import { Uploader } from "uppy-react"

const UploaderComponent = () => {
    return (
        <Uploader
            ${properties}
        />
    )
}

export default UploaderComponent
`;

  return code.replace(/[\r\n]+/g, '\n').trim();
};

const uploadEndpoint = getUploadApiDomain();

const ReactIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
  const settingsProps = useSettingsProps();
  const [apiKey, setApiKey] = useState('');
  const [changeApiKey, setChangeApiKey] = useState(false);
  const { workspaceId } = useParams();
  const { isPending: isUserPending, user } = useGetUserDetails();
  const code = getCode(workspaceId as string, uploadEndpoint, settingsProps);

  if (!workspaceId || isUserPending) {
    return <AppLoader h="50vh" />;
  }

  return (
    <Stack justify="center" align="center" pt="sm" mb={50}>
      <Timeline
        active={3}
        bulletSize={24}
        lineWidth={2}
        w={{ sm: '100vw', md: '70vw', lg: '60vw' }}
      >
        <Timeline.Item bullet={<IconEditCircle size={12} />} title="Customize">
          <Text opacity={0.7} size="sm" mb="lg">
            Customize your uploader to match your brand and requirements
          </Text>
          <Paper>
            <Grid>
              <Grid.Col span={12}>
                <Group justify="flex-end" p={0} m={0} mr="md">
                  <ActionIcon
                    variant="subtle"
                    c="dimmed"
                    onClick={() => setChangeApiKey(!changeApiKey)}
                  >
                    <Tooltip label="Change API Key">
                      <IconKey size={16} />
                    </Tooltip>
                  </ActionIcon>
                </Group>
                <Group
                  mx="sm"
                  mt="xs"
                  p="30"
                  h="100%"
                  justify="center"
                  align="center"
                  style={{
                    overflow: 'auto',
                    background:
                      settingsProps.theme === 'light' ? '#ccc' : '#1e1e1e',
                    backgroundSize: '10px 10px',
                    borderRadius: '20px',
                  }}
                >
                  {sessionStorage.getItem(TEMP_API_KEY) && !changeApiKey ? (
                    <Uploader
                      uploadEndpoint={uploadEndpoint}
                      workspaceId={workspaceId}
                      metadata={{
                        uploaderEmail: user.email,
                        uploaderName: user.name || 'User',
                      }}
                      headers={{
                        'X-Tenant-Id':
                          sessionStorage.getItem(TENANT_ID_KEY) || '',
                        'X-Api-Key': sessionStorage.getItem(TEMP_API_KEY) || '',
                      }}
                      {...settingsProps}
                      note="Test your uploader"
                    />
                  ) : (
                    <Box w="90%">
                      <PasswordInput
                        label="Please enter your API key"
                        placeholder="Please enter your API key"
                        required
                        onChange={event => setApiKey(event.target.value)}
                      />
                      <Button
                        mt="md"
                        onClick={() => {
                          sessionStorage.setItem(TEMP_API_KEY, apiKey);
                          window.location.reload();
                        }}
                      >
                        Submit
                      </Button>
                    </Box>
                  )}
                </Group>
              </Grid.Col>
              <Grid.Col span={12}>
                <Box p="xl">
                  <Settings {...settingsProps} />
                </Box>
              </Grid.Col>
            </Grid>
          </Paper>
        </Timeline.Item>
        <Timeline.Item
          bullet={<IconBrandNpm size={12} />}
          title="Install package"
        >
          <Text opacity={0.7} size="sm" mb="lg">
            Install our library from npm
          </Text>
          <Paper p="lg">
            <CodeHighlight m="sm" code={`npm install uppy-react`} />
          </Paper>
        </Timeline.Item>

        <Timeline.Item bullet={<IconCode size={12} />} title="Code">
          <Text opacity={0.7} size="sm" mb="lg">
            Based on your framework, add the code to your page
          </Text>
          <Paper p="lg">
            <CodeHighlight m="sm" code={code} language="tsx" />
          </Paper>
        </Timeline.Item>

        <Timeline.Item bullet={<IconConfetti size={12} />} title="Cheers">
          <Text opacity={0.7} size="sm" mb="lg">
            You did it, Check your imported files in the import section or
            configure from the configuration section
          </Text>
        </Timeline.Item>
      </Timeline>
    </Stack>
  );
};

export default ReactIntegrationPage;
