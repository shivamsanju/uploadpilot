import { Card, Group, Paper, SegmentedControl, SimpleGrid, Text, useMantineColorScheme } from "@mantine/core"
import { Uploader } from "uploadpilotreact"
import { useState } from "react";
import { IconCode, IconEye } from "@tabler/icons-react";
import { CodeHighlightTabs } from "@mantine/code-highlight";
import '@mantine/code-highlight/styles.css';
import { useParams } from "react-router-dom";
import AppLoader from "../../../../components/Loader/AppLoader";
import Settings from "./Settings/Settings";

const getCode = (uploaderId: string, backendEndpoint: string, h: number, w: number) => `
import Uploader from "upload-pilot"

const Component = () => {
    return (
        <Uploader 
            uploaderId="${uploaderId}"
            backendEndpoint="${backendEndpoint}"
            h={${h}} 
            w={${w}}
        />
    )
}

export default Component
`


const ConfigurationUI = ({ uploaderDetails }: { uploaderDetails: any }) => {
    const [viewMode, setViewMode] = useState<string>('settings');
    const [height, setHeight] = useState<number>(600);
    const [width, setWidth] = useState<number>(500);
    const [backendEndpoint, setBackendEndpoint] = useState<string>('');

    const { uploaderId } = useParams();
    const { colorScheme } = useMantineColorScheme();


    const code = getCode(uploaderId as string, backendEndpoint, height, width);


    return !uploaderId ? <AppLoader h="70vh" /> : (
        <SimpleGrid cols={2}>
            <Card withBorder shadow="xs" radius="sm" h="70vh">
                <Card.Section withBorder inheritPadding py="xs">
                    <Group justify="space-between">
                        <Text fw={500}>Configuration</Text>
                        <SegmentedControl
                            w="250"
                            value={viewMode}
                            onChange={(e) => setViewMode(e)}
                            data={[
                                {
                                    value: 'settings',
                                    label: (
                                        <Group align="center" justify="center" gap="sm">
                                            <IconEye size={16} stroke={1.5} />
                                            <div>Settings</div>
                                        </Group>
                                    ),
                                },
                                {
                                    value: 'code',
                                    label: (
                                        <Group align="center" justify="center" gap="sm">
                                            <IconCode size={16} stroke={1.5} />
                                            <div>Code</div>
                                        </Group>
                                    ),
                                },
                            ]}
                        />
                    </Group>

                </Card.Section>
                <Card.Section py="xs">

                    {viewMode === 'settings' ? (
                        <Settings
                            height={height}
                            setHeight={setHeight}
                            width={width}
                            setWidth={setWidth}
                            backendEndpoint={backendEndpoint}
                            setBackendEndpoint={setBackendEndpoint}
                            uploaderDetails={uploaderDetails}
                        />
                    ) :
                        <CodeHighlightTabs
                            m="sm"
                            code={[
                                { fileName: 'React', code: code, language: 'tsx' },
                                { fileName: 'Angular', code: code, language: 'tsx' },
                                { fileName: 'Vue', code: code, language: 'tsx' },
                                { fileName: 'Svelte', code: code, language: 'tsx' },
                            ]}
                        />
                    }
                </Card.Section>
            </Card>
            <Paper
                withBorder
                shadow="xs"
                radius="sm"
                style={{
                    backgroundImage: colorScheme === 'light' ? 'radial-gradient(circle,rgb(219, 219, 219) 1px, transparent 1px)' : 'radial-gradient(circle, #3e3e3e 1px, transparent 1px)',
                    backgroundSize: '10px 10px',
                }}>
                <Group justify="center" align="center" style={{ overflow: 'auto' }} h="70vh">
                    <Uploader
                        backendEndpoint="http://localhost:8080"
                        uploaderId={uploaderId}
                        h={height}
                        w={width}
                    />
                </Group>
            </Paper>
        </SimpleGrid>
    )
}

export default ConfigurationUI