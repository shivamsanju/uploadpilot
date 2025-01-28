import { Grid } from "@mantine/core"
import { useParams } from "react-router-dom"
import { AppLoader } from "../../components/Loader/AppLoader"
import { ProcessorCanvas } from "./Canvas";
import { Toolbar } from "./Toolbar";
import classes from "./Processor.module.css"
import { ProcEditorHeader } from "./Header";
import { ReactFlowProvider } from "@xyflow/react";
import { DnDProvider } from "../../context/DnD";
import { ProcEditorProvider } from "../../context/EditorCtx";

const ProcessorPage = () => {
    const { workspaceId } = useParams();

    if (!workspaceId) {
        return <AppLoader h="70vh" />
    }

    return (
        <ReactFlowProvider>
            <ProcEditorProvider>
                <DnDProvider>
                    <Grid>
                        <Grid.Col span={12} className={classes.header}>
                            <ProcEditorHeader />
                        </Grid.Col>
                        <Grid.Col span={3} className={classes.toolbar}>
                            <Toolbar />
                        </Grid.Col>
                        <Grid.Col span={9} className={classes.canvas} p={0} m={0}>
                            <ProcessorCanvas />
                        </Grid.Col>
                    </Grid>
                </DnDProvider>
            </ProcEditorProvider>
        </ReactFlowProvider>
    )
}

export default ProcessorPage