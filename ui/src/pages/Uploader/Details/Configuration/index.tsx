import { Button, Card, Group, SegmentedControl, SimpleGrid, useMantineColorScheme } from "@mantine/core"
import { Uploader } from "@uploadpilot/react"
import { useState } from "react";
import { IconCancel, IconCode, IconDeviceFloppy, IconEdit, IconEye } from "@tabler/icons-react";
import { CodeHighlightTabs } from "@mantine/code-highlight";
import '@mantine/code-highlight/styles.css';
import { useParams } from "react-router-dom";
import AppLoader from "../../../../components/Loader/AppLoader";
import Settings from "./Settings/Settings";
import { getApiDomain } from "../../../../utils/config";
import { CreateUploaderForm } from "../../../../types/uploader";
import { useForm } from "@mantine/form";
import { useUpdateUploaderConfigMutation } from "../../../../apis/uploader";
import { useGetCurrentUserDetails } from "../../../../apis/user";

const getCode = (uploaderId: string, backendEndpoint: string, h: number, w: number, theme: 'auto' | 'light' | 'dark' = 'auto') => `
import { Uploader } from "@uploadpilot/react"

const Component = () => {
    return (
        <Uploader 
            uploaderId="${uploaderId}"
            backendEndpoint="${backendEndpoint}"
            height={${h}} 
            width={${w}}
            theme={${theme}}
            metadata={{"key": "value"}}
            headers={{"key": "value"}}
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
    const [theme, setTheme] = useState<'auto' | 'light' | 'dark'>('auto');
    const [editMode, setEditMode] = useState<boolean>(false);
    const [key, refreshKey] = useState<string>("0");

    const { uploaderId } = useParams();
    const { mutateAsync } = useUpdateUploaderConfigMutation(uploaderId as string)
    const { isPending: isUserPending, me } = useGetCurrentUserDetails();
    const { colorScheme } = useMantineColorScheme();
    const code = getCode(uploaderId as string, backendEndpoint, height, width, theme);
    const form = useForm<CreateUploaderForm>({
        initialValues: { ...uploaderDetails.config, requiredMetadataFields: uploaderDetails.requiredMetadataFields || [] },
    });


    const handleEditAndSaveButton = async () => {
        if (editMode) {
            mutateAsync(form.values)
                .then(() => refreshKey(prev => prev === "0" ? "1" : "0"))
        }
        setEditMode(prev => !prev)
    }

    const handleCancelButton = () => {
        setEditMode(false)
    }

    return (!uploaderId || isUserPending) ? <AppLoader h="50vh" /> : (
        <Card withBorder >
            <Card.Section withBorder inheritPadding py="xs">
                <Group justify="space-between">
                    <Group justify="flex-end" gap="md">
                        {editMode && <Button
                            variant="light"
                            onClick={handleCancelButton}
                            leftSection={<IconCancel size={18} />}
                        >
                            Cancel
                        </Button>}
                        <Button
                            variant={editMode ? undefined : "subtle"}
                            onClick={handleEditAndSaveButton}
                            leftSection={editMode ? <IconDeviceFloppy size={18} /> : <IconEdit size={18} />}
                        >
                            {editMode ? "Save" : "Edit"}
                        </Button>
                    </Group>
                    <Group align="center" gap="md">

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
                </Group>
            </Card.Section>
            <Card.Section>
                <SimpleGrid cols={2} h="72vh">
                    {viewMode === 'settings' ? (
                        <Settings
                            height={height}
                            setHeight={setHeight}
                            width={width}
                            setWidth={setWidth}
                            theme={theme}
                            setTheme={setTheme}
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
                    <Group justify="center" align="center" style={{
                        overflow: 'auto',
                        backgroundImage: colorScheme === 'light' ? 'radial-gradient(circle,rgb(219, 219, 219) 1px, transparent 1px)' : 'radial-gradient(circle, #3e3e3e 1px, transparent 1px)',
                        backgroundSize: '10px 10px',
                    }}>
                        <Uploader
                            key={key}
                            backendEndpoint={backendEndpoint}
                            uploaderId={uploaderId}
                            height={height}
                            width={width}
                            theme={theme}
                            metadata={{
                                "uploaderEmail": me.email,
                                "uploaderName": me.firstName + " " + me.lastName
                            }}
                        />
                    </Group>
                </SimpleGrid>
            </Card.Section>
        </Card>
    )
}

export default ConfigurationUI