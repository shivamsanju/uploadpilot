import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import { MantineStyleProp, Paper, Stack, Text, Timeline } from '@mantine/core';
import { IconBrandNpm, IconCode, IconConfetti } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { useGetSession } from '../../../apis/user';
import { AppLoader } from '../../../components/Loader/AppLoader';
import { getUploadApiDomain } from '../../../utils/config';

const uploadEndpoint = getUploadApiDomain();

const GoIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
  const { workspaceId } = useParams();
  const { isPending: isUserPending } = useGetSession();

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
        <Timeline.Item
          bullet={<IconBrandNpm size={12} />}
          title="Install package"
        >
          <Text opacity={0.7} size="sm" mb="lg">
            Install the go sdk
          </Text>
          <Paper p="lg">
            <CodeHighlight
              m="sm"
              code={`go get github.com/uploadpilot/sdk/go-client`}
            />
          </Paper>
        </Timeline.Item>

        <Timeline.Item bullet={<IconCode size={12} />} title="Code">
          <Text opacity={0.7} size="sm" mb="lg">
            Use the go sdk
          </Text>
          <Paper p="lg">
            <CodeHighlight
              m="sm"
              language="go"
              code={`
package main

import (
	"log"
	"os"

	gocl "github.com/uploadpilot/sdk/go-client"
)

func main() {
	log.Println("uploading file start...")
	cl, err := gocl.NewUploader(
          "${workspaceId}", // replace with your workspace id
          "${uploadEndpoint}", // uploadpilot endpoint
          "your_api_key", // replace with your api key
        )
	if err != nil {
		log.Fatal("failed to create uploader:", err)
		os.Exit(1)
	}
	err = cl.UploadFile(
		"../wf.yaml",
		&gocl.UploadOptions{
			Metadata: map[string]string{
				"filename": "wf.yaml",
				"filetype": "text/plain",
			},
		},
	)
	if err != nil {
		log.Fatal("failed to upload file:", err)
		os.Exit(1)
	}
	log.Println("file uploaded successfully!")
}`}
            />
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

export default GoIntegrationPage;
