import { Group, Paper, SegmentedControl, SimpleGrid, Stack, useMantineColorScheme } from "@mantine/core"
import { Uploader } from "@uploadpilot/react"
import { useState } from "react";
import { IconCode, IconEye } from "@tabler/icons-react";
import { CodeHighlightTabs } from "@mantine/code-highlight";
import '@mantine/code-highlight/styles.css';
import { useParams } from "react-router-dom";
import AppLoader from "../../../../components/Loader/AppLoader";
import { getApiDomain } from "../../../../utils/config";
import { useGetSession } from "../../../../apis/user";
import Settings from "./Settings";
import { useSettingsProps } from "./SettingsProps";

const getCode = (uploaderId: string, backendEndpoint: string, settingsProps: any) => {
    const properties = [
        `uploaderId="${uploaderId}"`,
        `backendEndpoint="${backendEndpoint}"`,
        `height={${settingsProps.height}}`,
        `width={${settingsProps.width}}`,
        `theme={${settingsProps.theme}}`,
        `metadata={{"key": "value"}}`,
        `headers={{"key": "value"}}`,
        !settingsProps.showStatusBar && 'showStatusBar={false}',
        !settingsProps.showProgress && 'showProgress={false}',
        settingsProps.hideUploadButton && 'hideUploadButton={true}',
        settingsProps.hideCancelButton && 'hideCancelButton={true}',
        settingsProps.hideRetryButton && 'hideRetryButton={true}',
        settingsProps.hidePauseResumeButton && 'hidePauseResumeButton={true}',
        settingsProps.hideProgressAfterFinish && 'hideProgressAfterFinish={true}',
        settingsProps.note && `note="${settingsProps.note}"`,
        !settingsProps.singleFileFullScreen && 'singleFileFullScreen={false}',
        !settingsProps.showSelectedFiles && 'showSelectedFiles={false}',
    ]
        .filter(Boolean) // Remove any falsy values
        .join('\n            '); // Join with proper indentation

    const code = `
import { Uploader } from "@uploadpilot/react"

const UploaderComponent = () => {
    return (
        <Uploader
            ${properties}
        />
    )
}

export default UploaderComponent
`;

    return code.replace(/[\r\n]+/g, '\n').trim();
};


const backendEndpoint = getApiDomain();

const PreviewTab = () => {
    const [viewMode, setViewMode] = useState<string>('settings');

    // Uploader props
    const settingsProps = useSettingsProps();


    const { uploaderId } = useParams();
    const { isPending: isUserPending, session } = useGetSession();
    const { colorScheme } = useMantineColorScheme();
    const code = getCode(uploaderId as string, backendEndpoint, settingsProps);


    return (!uploaderId || isUserPending) ? <AppLoader h="50vh" /> : (
        <SimpleGrid cols={2} h="75vh" >
            <Stack gap="md" pt="md">
                <Group justify="center">
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
                {viewMode === 'settings' ? (
                    <Settings
                        {...settingsProps}
                    />
                ) :
                    <Paper>
                        <CodeHighlightTabs
                            m="sm"
                            code={[
                                { fileName: 'React', code: code, language: 'tsx' },
                                { fileName: 'Angular', code: code, language: 'tsx' },
                                { fileName: 'Vue', code: code, language: 'tsx' },
                                { fileName: 'Svelte', code: code, language: 'tsx' },
                            ]}
                        />
                    </Paper>
                }
            </Stack>
            <Group justify="center" align="center" style={{
                overflow: 'auto',
                backgroundImage: colorScheme === 'light' ? 'radial-gradient(circle,rgb(219, 219, 219) 1px, transparent 1px)' : 'radial-gradient(circle, #3e3e3e 1px, transparent 1px)',
                backgroundSize: '10px 10px',
            }}>
                <Uploader
                    backendEndpoint={backendEndpoint}
                    uploaderId={uploaderId}
                    metadata={{
                        "uploaderEmail": session.email,
                        "uploaderName": session.name || "sss"
                    }}
                    {...settingsProps}
                />
            </Group>
        </SimpleGrid>
    )
}

export default PreviewTab