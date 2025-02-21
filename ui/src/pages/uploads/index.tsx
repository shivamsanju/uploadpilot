import { Box, Title, Group, Badge, Alert } from "@mantine/core";
import UploadList from "./List";
import { useState } from "react";
import { IconInfoCircle } from "@tabler/icons-react";
import { useLocalStorage } from "@mantine/hooks";

const UploadsPage = () => {
  const [totalRecords, setTotalRecords] = useState(0);
  const [alert, setAlert] = useLocalStorage({
    key: "alert",
    defaultValue: 1,
    getInitialValueInEffect: false,
  });

  return (
    <Box>
      <Group align="center" gap="xs" h="10%">
        <Title order={3} opacity={0.7}>
          Uploads
        </Title>
        <Badge variant="outline">{totalRecords}</Badge>
      </Group>
      {alert === 1 && (
        <Alert
          opacity={0.7}
          icon={<IconInfoCircle size={16} />}
          withCloseButton
          onClose={() => setAlert(0)}
        >
          Uploaded files may take a few seconds to appear here after upload.
          Please click the refresh button if you don't see any new files.
        </Alert>
      )}
      <UploadList setTotalRecords={setTotalRecords} />
    </Box>
  );
};

export default UploadsPage;
