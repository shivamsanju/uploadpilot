import { Stack } from "@mantine/core"
import ImportsList from "./List"

type ImportsProps = {
    uploaderDetails: any
}
const Imports: React.FC<ImportsProps> = ({ uploaderDetails }) => {
    return (
        <Stack gap="md" >
            <ImportsList />
        </Stack>
    )
}

export default Imports