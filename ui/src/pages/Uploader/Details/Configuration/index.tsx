import { Box, Button, Card, Group, Paper, SegmentedControl, SimpleGrid, Text, useMantineColorScheme } from "@mantine/core"
import { Uploader } from "@uploadpilot/react"
import { useState } from "react";
import { IconCode, IconDeviceFloppy, IconEdit, IconEye } from "@tabler/icons-react";
import { CodeHighlightTabs } from "@mantine/code-highlight";
import '@mantine/code-highlight/styles.css';
import { useParams } from "react-router-dom";
import AppLoader from "../../../../components/Loader/AppLoader";
import Settings from "./Settings/Settings";
import { getApiDomain } from "../../../../utils/config";
import { CreateUploaderForm } from "../../../../types/uploader";
import { useForm } from "@mantine/form";
import { useUpdateUploaderConfigMutation } from "../../../../apis/uploader";

const getCode = (uploaderId: string, backendEndpoint: string, h: number, w: number) => `
import { Uploader } from "@uploadpilot/react"

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

const backendEndpoint = getApiDomain();

const ConfigurationUI = ({ uploaderDetails }: { uploaderDetails: any }) => {
    const [viewMode, setViewMode] = useState<string>('settings');
    const [height, setHeight] = useState<number>(600);
    const [width, setWidth] = useState<number>(500);
    const [editMode, setEditMode] = useState<boolean>(false);
    const [key, refreshKey] = useState<string>("0");

    const { uploaderId } = useParams();
    const { mutateAsync } = useUpdateUploaderConfigMutation(uploaderId as string)
    const { colorScheme } = useMantineColorScheme();
    const code = getCode(uploaderId as string, backendEndpoint, height, width);
    const form = useForm<CreateUploaderForm>({
        initialValues: { ...uploaderDetails.config, requiredMetadataFields: uploaderDetails.requiredMetadataFields || [] },
    });


    const handleEditButton = async () => {
        if (editMode) {
            mutateAsync(form.values)
                .then(() => refreshKey(prev => prev === "0" ? "1" : "0"))
        }
        setEditMode(prev => !prev)
    }

    return !uploaderId ? <AppLoader h="67vh" /> : (
        <Box>
            <Group justify="flex-end" mb="xs">
                <Button
                    size="xs"
                    variant="default"
                    onClick={handleEditButton}
                    leftSection={editMode ? <IconDeviceFloppy size={20} /> : <IconEdit size={20} />}
                >
                    {editMode ? "Save" : "Edit"}
                </Button>
            </Group>

            <SimpleGrid cols={2}>
                <Card withBorder shadow="xs" radius="sm" h="67vh">
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
                                editMode={editMode}
                                form={form}
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
                    <Group justify="center" align="center" style={{ overflow: 'auto' }} h="67vh">
                        <Uploader
                            key={key}
                            backendEndpoint={backendEndpoint}
                            uploaderId={uploaderId}
                            h={height}
                            w={width}
                        />
                    </Group>
                </Paper>
            </SimpleGrid>
        </Box>
    )
}

export default ConfigurationUI