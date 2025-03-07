import { CodeHighlight } from '@mantine/code-highlight';
import '@mantine/code-highlight/styles.css';
import {
  Button,
  FileInput,
  Group,
  MantineStyleProp,
  Stack,
  Text,
  Timeline,
} from '@mantine/core';
import {
  IconBrandNpm,
  IconCircleCheck,
  IconCode,
  IconConfetti,
  IconExclamationCircle,
} from '@tabler/icons-react';
import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { getTenantId } from '../../../utils/config';
import { Uploader } from '../../../utils/uploader';

const BrowserIntegrationPage = ({ style }: { style: MantineStyleProp }) => {
  const [files, setFiles] = useState<File[]>([]);
  const [status, setStatus] = useState('');
  const { workspaceId } = useParams();

  const uploader = new Uploader(getTenantId()!, workspaceId!, '');

  const handleUpload = async () => {
    try {
      setStatus('Uploading...');
      await uploader.uploadMultiple(files, {
        sub: 'testing',
      });
      setStatus('Upload complete');
    } catch (error) {
      setStatus('Upload failed');
    }
  };

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

const uploader = new Uploader(
  '${getTenantId()}',
  '${workspaceId}',
  'YOUR_API_KEY',
);

async function handleFileUpload(event) {
  try {
    const files = event.target.files;
    if (!files.length) {
      console.error('No files selected');
      return;
    }

    await uploader.uploadMultiple(files, {
      "metadata-key": "metadata-value",
    });
    
    console.log('Files uploaded successfully');
  } catch (error) {
    console.error('Error uploading file:', error);
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.multiple = true;
  input.addEventListener('change', handleFileUpload);
  document.body.appendChild(input);
});
`}
          />
        </Timeline.Item>

        <Timeline.Item
          bullet={<IconBrandNpm size={12} />}
          title="Test your uploader"
        >
          <Text opacity={0.7} size="sm" mb="lg">
            Upload Files
          </Text>
          <Group
            align="center"
            ml="md"
            mb={4}
            gap={3}
            c={
              status === 'Upload complete'
                ? 'green'
                : status === 'Upload failed'
                  ? 'red'
                  : 'yellow'
            }
          >
            {status === 'Upload complete' && <IconCircleCheck size={16} />}
            {status === 'Upload failed' && <IconExclamationCircle size={16} />}
            <Text>{status}</Text>
          </Group>
          <Group align="center" justify="space-between" gap="md" px="sm">
            <FileInput
              flex={1}
              multiple
              placeholder="Upload files"
              radius="md"
              onChange={setFiles}
            />
            <Button radius="md" onClick={handleUpload}>
              Upload
            </Button>
          </Group>
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

export default BrowserIntegrationPage;
