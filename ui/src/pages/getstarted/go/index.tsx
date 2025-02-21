import { Paper, Stack, Text, Timeline } from "@mantine/core";
import { IconBrandNpm, IconCode, IconConfetti } from "@tabler/icons-react";
import { CodeHighlight } from "@mantine/code-highlight";
import "@mantine/code-highlight/styles.css";
import { useParams } from "react-router-dom";
import { AppLoader } from "../../../components/Loader/AppLoader";
import { useGetSession } from "../../../apis/user";
import { useViewportSize } from "@mantine/hooks";

const GoSdkIntegrationPage = () => {
  const { width } = useViewportSize();
  const { workspaceId } = useParams();
  const { isPending: isUserPending } = useGetSession();

  if (!workspaceId || isUserPending) {
    return <AppLoader h="50vh" />;
  }

  // TODO: Heavy engineering: Need to find some smarter way to do this
  const style = () => {
    if (width > 768) {
      return {};
    }

    let scale = 1;
    if (width < 768 && width > 700) {
      scale = width / 768;
    } else if (width < 700 && width > 500) {
      scale = (width / 768) * 1.1;
    } else {
      scale = (width / 768) * 1.35;
    }

    return {
      transform: `scale(${scale})`,
      transformOrigin: "top left",
    };
  };

  return (
    <Stack justify="center" align="center" pt="sm" mb={50} style={style()}>
      <Timeline active={3} bulletSize={24} lineWidth={2} w="70%">
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

	gocl "github.com/uploadpilot/sdk/go-client"
)

func main() {
	err := gocl.UploadFile(
		"ddd6d6d7-5805-419d-8e86-1976d626f79b",
		"go.sum",
		&gocl.UploadOptions{
			BaseURI:  "http://localhost:8081",
			Metadata: map[string]string{"test": "test"},
			Headers:  map[string]string{"test": "test"},
		},
	)
	if err != nil {
		panic(err)
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

export default GoSdkIntegrationPage;
