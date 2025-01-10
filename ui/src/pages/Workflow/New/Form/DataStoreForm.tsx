import { Stack, Select, TextInput } from "@mantine/core";
import { optionsFilter } from "../../../../utils/filter";
import { useGetStorageConnectors } from "../../../../apis/storage";
import { Connector } from "../../../../types/connector";
import { useState } from "react";


const DataStorePage = ({ form }: { form: any }) => {
    const [selectedConnector, setSelectedConnector] = useState<Connector>();

    const { isPending, error, connectors } = useGetStorageConnectors();


    return (<Stack gap="md">
        <Select
            withAsterisk
            filter={optionsFilter}
            label="Storage connector"
            placeholder="Select a storage connector"
            data={isPending ? [] : connectors.map((connector: Connector) => ({ ...connector, value: connector.id, label: connector.name }))}
            onChange={(value) => {
                setSelectedConnector(connectors.find((connector: Connector) => connector.id === value))
                form.setFieldValue("connectorId", value);
            }}
        />

        {selectedConnector?.id && (
            <>
                <TextInput
                    withAsterisk
                    label="Datastore Name"
                    placeholder="Enter a datastore name"
                    {...form.getInputProps("dataStoreName")}

                />
                <TextInput
                    withAsterisk
                    label={selectedConnector.type === "azure" ? "Container Name" : "Bucket Name"}
                    placeholder={`Enter a ${selectedConnector.type === "azure" ? "container" : "bucket"} name`}
                    {...form.getInputProps("bucket")}
                />
            </>
        )}
    </Stack>
    )
};

export default DataStorePage;
