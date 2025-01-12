import { useNavigate, useParams } from 'react-router-dom';
import { Tabs } from '@mantine/core';
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
            <Tabs.List mb="sm" grow>
                <Tabs.Tab value=".">Configuration</Tabs.Tab>
                <Tabs.Tab value="imports">Imports</Tabs.Tab>
            </Tabs.List>
            <Tabs.Panel value=".">
                <Configuration uploaderDetails={uploaderDetails} />
            </Tabs.Panel>
            <Tabs.Panel value="imports">
                <Imports uploaderDetails={uploaderDetails} />
            </Tabs.Panel>
        </Tabs>
    );
}

export default UploaderTabs;