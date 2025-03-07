import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import { MantineStyleProp, Stack, Text, Timeline } from '@mantine/core';
import { IconBrandNpm, IconCode, IconConfetti } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { getTenantId } from '../../../utils/config';

const GoIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
  const { workspaceId } = useParams();

  return (
    <Stack justify="center" align="center" pt="sm" mb={50}>
      <Timeline
        active={3}
        bulletSize={24}
        lineWidth={2}
        w={{ sm: '100vw', md: '70vw', lg: '60vw' }}
      >
        <Timeline.Item
          bullet={<IconBrandNpm size={12} />}
          title="Install package"
        >
          <Text opacity={0.7} size="sm" mb="lg">
            Install the go sdk
          </Text>
          <CodeHighlight
            m="sm"
            code={`go get github.com/uploadpilot/sdk/go-client`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconCode size={12} />} title="Code">
          <Text opacity={0.7} size="sm" mb="lg">
            Use the go sdk
          </Text>
          <CodeHighlight
            m="sm"
            code={`
package main

import (
  "log"
  "os"

  "github.com/uploadpilot/go-sdk/client"
)

func main() {
  uploader, err := client.NewUploader(
    "${getTenantId()}",
    "${workspaceId}",
    "YOUR_API_KEY",
  )
  if err != nil {
    log.Fatal("Failed to create uploader:", err)
    os.Exit(1)
  }

  err = uploader.UploadFile(
    "../wf.yaml",
    &client.UploadOptions{
      FileName: "wf.yaml",
      Metadata: map[string]string{
        "key": "value",
      },
    },
  )
  if err != nil {
    log.Fatal("Failed to upload file:", err)
    os.Exit(1)
  }
}
`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconConfetti size={12} />} title="Cheers">
          <Text opacity={0.7} size="sm" mb="lg">
            You did it, Start uploading and check your uploaded files in the
            uploads section or configure from the configuration section
          </Text>
        </Timeline.Item>
      </Timeline>
    </Stack>
  );
};

export default GoIntegrationPage;
