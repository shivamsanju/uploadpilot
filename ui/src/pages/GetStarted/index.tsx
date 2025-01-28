import { Box, Group, Paper, SimpleGrid, Stack, Text, Timeline, useMantineColorScheme } from "@mantine/core"
import { Uploader } from "uppy-react"
import { IconBrandNpm, IconCode, IconConfetti, IconEditCircle } from "@tabler/icons-react";
import { CodeHighlight, CodeHighlightTabs } from "@mantine/code-highlight";
import '@mantine/code-highlight/styles.css';
import { useParams } from "react-router-dom";
import { AppLoader } from "../../components/Loader/AppLoader";
import { getApiDomain } from "../../utils/config";
import { useGetSession } from "../../apis/user";
import Settings from "./Settings";
import { useSettingsProps } from "./SettingsProps";
import { useViewportSize } from "@mantine/hooks";

const getCode = (workspaceId: string, backendEndpoint: string, settingsProps: any) => {
    const properties = [
        `workspaceId="${workspaceId}"`,
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
import { Uploader } from "uppy-react"

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

const UploaderPreviewPage = () => {
    const settingsProps = useSettingsProps();
    const { width } = useViewportSize();
    const { workspaceId } = useParams();
    const { isPending: isUserPending, session } = useGetSession();
    const { colorScheme } = useMantineColorScheme();
    const code = getCode(workspaceId as string, backendEndpoint, settingsProps);

    if (!workspaceId || isUserPending) {
        return <AppLoader h="50vh" />
    }

    // TODO: Heavy engineering: Need to find some smarter way to do this
    const style = () => {
        if (width > 768) {
            return {};
        }

        let scale = 1;
        if (width < 768 && width > 700) {
            scale = (width / 768);
        } else if (width < 700 && width > 500) {
            scale = (width / 768) * 1.1;
        } else {
            scale = (width / 768) * 1.35;
        }

        return {
            transform: `scale(${scale})`,
            transformOrigin: 'top left',
        };
    }

    return (
        <Stack justify="center" align="center" pt="sm" mb={50} style={style()}>
            <Timeline active={3} bulletSize={24} lineWidth={2}>
                <Timeline.Item bullet={<IconEditCircle size={12} />} title="Customize">
                    <Text opacity={0.7} size="sm" mb="lg">Customize your uploader to match your requirements</Text>
                    <Paper  >
                        <SimpleGrid cols={{ sm: 1, md: 2 }}>
                            <Box p="xl">
                                <Settings
                                    {...settingsProps}
                                />
                            </Box>
                            <Group justify="center" align="center" style={{
                                overflow: 'auto',
                                backgroundImage: colorScheme === 'light' ? 'radial-gradient(circle,rgb(219, 219, 219) 1px, transparent 1px)' : 'radial-gradient(circle, #3e3e3e 1px, transparent 1px)',
                                backgroundSize: '10px 10px',
                            }}>

                                <Uploader
                                    backendEndpoint={backendEndpoint}
                                    workspaceId={workspaceId}
                                    metadata={{
                                        "uploaderEmail": session.email,
                                        "uploaderName": session.name || "sss"
                                    }}
                                    {...settingsProps}
                                    note="Test your uploader"
                                    headers={{ "Authorization": "Bearer mysecrettoken" }}
                                />
                            </Group>
                        </SimpleGrid>
                    </Paper>
                </Timeline.Item>
                <Timeline.Item bullet={<IconBrandNpm size={12} />} title="Install package" >
                    <Text opacity={0.7} size="sm" mb="lg">Install our library from npm</Text>
                    <Paper p="lg">
                        <CodeHighlight
                            m="sm"
                            code={`npm install uppy-react`}
                        />
                    </Paper>
                </Timeline.Item>

                <Timeline.Item bullet={<IconCode size={12} />} title="Code" >
                    <Text opacity={0.7} size="sm" mb="lg">Based on your framework, add the code to your page</Text>
                    <Paper p="lg">
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
                </Timeline.Item>

                <Timeline.Item bullet={<IconConfetti size={12} />} title="Cheers" >
                    <Text opacity={0.7} size="sm" mb="lg">You did it, Check your imported files in the import section or configure from the configuration section</Text>
                </Timeline.Item>

            </Timeline>
        </Stack>
    )
}

export default UploaderPreviewPage