import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import { MantineStyleProp, Stack, Text, Timeline } from '@mantine/core';
import { IconBrandNpm, IconCode, IconConfetti } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { getTenantId } from '../../../utils/config';

const PythonIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
  const { workspaceId } = useParams();

  return (
    <Stack mb={50}>
      <Timeline active={3} bulletSize={24} lineWidth={2}>
        <Timeline.Item
          bullet={<IconBrandNpm size={12} />}
          title="Install package"
        >
          <Text opacity={0.7} mb="lg">
            Install the go sdk
          </Text>
          <CodeHighlight
            m="sm"
            language="bash"
            code={`go get github.com/uploadpilot/sdk/go-client`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconCode size={12} />} title="Code">
          <Text opacity={0.7} mb="lg">
            Use the go sdk
          </Text>
          <CodeHighlight
            m="sm"
            language="go"
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
    &client.UploaderOpts{
      BaseURL: "http://localhost:8080",
    }
  )

  fileData, err := os.ReadFile("./files/image.png")
  if err != nil {
    fmt.Println("Error reading file:", err)
    return
  }

  File := &client.File{
    Name:         "image.png",
    Data:         fileData,
    ContentType:  "image/png",
  }

  success, err := up.Upload(File, map[string]interface{}{"description": "Test file"})
  if err != nil {
    fmt.Println("Upload failed:", err)
  } else {
    fmt.Println("Upload successful:", success)
  }
}
`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconConfetti size={12} />} title="Cheers">
          <Text opacity={0.7} mb="lg">
            You did it, Start uploading and check your uploaded files in the
            uploads section or configure from the configuration section
          </Text>
        </Timeline.Item>
      </Timeline>
    </Stack>
  );
};

export default PythonIntegrationPage;
