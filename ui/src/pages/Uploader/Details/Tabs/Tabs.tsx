import { useNavigate, useParams } from 'react-router-dom';
import { Divider, Tabs } from '@mantine/core';
import Configuration from '../Configuration';
import Imports from '../Imports/Imports';
import classes from "./Tabs.module.css"

type UploaderTabsProps = {
    uploaderDetails: any
}

const UploaderTabs: React.FC<UploaderTabsProps> = ({ uploaderDetails }) => {
    const navigate = useNavigate();
    const { tabValue, uploaderId } = useParams();

    const handleTabChange = (value: string | null) => {
        if (!value || value === tabValue) return;
        if (value === ".") {
            navigate(`/uploaders/${uploaderId}`);
        }

        if (value && value !== tabValue) {
            navigate(`/uploaders/${uploaderId}/${value}`);
        }
    }

    return (
        <Tabs
            defaultValue="."
            value={tabValue || "."}
            onChange={handleTabChange}
            classNames={classes}
        >
            <Tabs.List mb="sm" grow={false} style={{ gap: "0.5rem" }}>
                <Tabs.Tab w={150} value=".">Configuration</Tabs.Tab>
                <Tabs.Tab w={150} value="imports">Imports</Tabs.Tab>
                <Tabs.Tab w={150} value="hooks">Hooks</Tabs.Tab>
            </Tabs.List>
            <Tabs.Panel value=".">
                <Configuration uploaderDetails={uploaderDetails} />
            </Tabs.Panel>
            <Tabs.Panel value="imports">
                <Imports uploaderDetails={uploaderDetails} />
            </Tabs.Panel>
            <Tabs.Panel value="hooks">
                Coming soon
            </Tabs.Panel>
        </Tabs>
    );
}

export default UploaderTabs;