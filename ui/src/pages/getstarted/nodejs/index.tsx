import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import { MantineStyleProp, Stack, Text, Timeline } from '@mantine/core';
import { IconBrandNpm, IconCode, IconConfetti } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { getTenantId } from '../../../utils/config';

const NodejsIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
  const { workspaceId } = useParams();

  return (
    <Stack mb={50}>
      <Timeline active={3} bulletSize={24} lineWidth={2}>
        <Timeline.Item
          bullet={<IconBrandNpm size={12} />}
          title="Install package"
        >
          <Text opacity={0.7} mb="lg">
            Install the npm library
          </Text>
          <CodeHighlight
            m="sm"
            language="bash"
            code={`npm i uploadpilot-uploader`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconCode size={12} />} title="Code">
          <Text opacity={0.7} mb="lg">
            Use the javascript/typescipt sdk
          </Text>
          <CodeHighlight
            m="sm"
            language="javascript"
            code={`
import fs from "fs";
import path from "path";
import { NodeFile, Uploader } from 'uploadpilot-uploader';

const uploader = new Uploader(
  '${getTenantId()}',
  '${workspaceId}',
  'YOUR_API_KEY',
);

async function main(): Promise<void> {
  try {
    const filePath: string = path.resolve("./files/a.txt"); // Replace with actual file path
    const fileBuffer: Buffer = fs.readFileSync(filePath);

    const file = new NodeFile(
      fileBuffer,
      path.basename(filePath),
      "application/octet-stream"
    );

    await uploader.uploadMultiple([file], {
      file_name: "a.txt",
    });

    console.log("File uploaded successfully");
  } catch (error) {
    console.error(
      "Error uploading file:",
      error instanceof Error ? error.message : error
    );
  }
}

main();
`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconConfetti size={12} />} title="Next steps">
          <Text opacity={0.7} mb="lg">
            You did it, Start uploading and check your uploaded files in the
            uploads section or configure from the configuration section
          </Text>
        </Timeline.Item>
      </Timeline>
    </Stack>
  );
};

export default NodejsIntegrationPage;
