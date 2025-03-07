import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import { MantineStyleProp, Stack, Text, Timeline } from '@mantine/core';
import { IconBrandNpm, IconCode, IconConfetti } from '@tabler/icons-react';
import { useParams } from 'react-router-dom';
import { getTenantId } from '../../../utils/config';

const NodejsIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
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
            Install the npm library
          </Text>
          <CodeHighlight m="sm" code={`npm i uploadpilot-uploader`} />
        </Timeline.Item>

        <Timeline.Item bullet={<IconCode size={12} />} title="Code">
          <Text opacity={0.7} size="sm" mb="lg">
            Use the javascript/typescipt sdk
          </Text>
          <CodeHighlight
            m="sm"
            code={`
import { Uploader } from 'uploadpilot-uploader';
import fs from 'fs';
import path from 'path';

const uploader = new Uploader(
  '${getTenantId()}',
  '${workspaceId}',
  'YOUR_API_KEY',
);

async function main() {
  try {
    const filePath = path.resolve('path/to/your/file.ext'); // Replace with actual file path
    const fileBuffer = fs.readFileSync(filePath);

    const file = new File([fileBuffer], path.basename(filePath), {
      type: 'application/octet-stream', // Change MIME type if necessary
    });

    await uploader.uploadMultiple([file], {
      "metadata-key": "metadata-value",
    });
    
    console.log('File uploaded successfully');
  } catch (error) {
    console.error('Error uploading file:', error);
  }
}

main();
`}
          />
        </Timeline.Item>

        <Timeline.Item bullet={<IconConfetti size={12} />} title="Next steps">
          <Text opacity={0.7} size="sm" mb="lg">
            You did it, Start uploading and check your uploaded files in the
            uploads section or configure from the configuration section
          </Text>
        </Timeline.Item>
      </Timeline>
    </Stack>
  );
};

export default NodejsIntegrationPage;
