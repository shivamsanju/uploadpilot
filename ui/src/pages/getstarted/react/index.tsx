import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import {
  Box,
  Grid,
  Group,
  MantineStyleProp,
  Paper,
  Stack,
  Text,
  Timeline,
} from '@mantine/core';
import {
  IconBrandNpm,
  IconCode,
  IconConfetti,
  IconEditCircle,
} from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { Uploader } from 'uppy-react';
import 'uppy-react/dist/style.css';
import { useGetSession } from '../../../apis/user';
import { AppLoader } from '../../../components/Loader/AppLoader';
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
  const { workspaceId } = useParams();
  const { isPending: isUserPending, session } = useGetSession();
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
                  <Uploader
                    uploadEndpoint={uploadEndpoint}
                    workspaceId={workspaceId}
                    metadata={{
                      uploaderEmail: session.email,
                      uploaderName: session.name || 'Unknown',
                    }}
                    {...settingsProps}
                    note="Test your uploader"
                    headers={{
                      Authorization:
                        'Bearer ' + localStorage.getItem('uploadpilottoken'),
                    }}
                  />
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
