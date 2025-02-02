import React from "react";
import { Container, Title, Text, List, Anchor, Divider } from "@mantine/core";

const TermsOfService: React.FC = () => {
  return (
    <Container size="sm" mt="lg">
      <Title order={1}>Terms of Service</Title>
      <Text size="sm" c="dimmed">
        Effective Date: [01/01/2025]
      </Text>

      <Divider my="lg" />

      <Section title="1. Acceptance of Terms">
        <Text size="sm">
          By accessing or using UploadPilot, you agree to comply with and be
          bound by these Terms of Service. If you do not agree to these terms,
          do not use our services.
        </Text>
      </Section>

      <Section title="2. Service Description">
        <Text size="sm">
          Our platform allows you to upload, process, and manage files from a
          variety of sources, including local files, cloud storage services, and
          social media platforms. Our service processes files in various ways
          depending on the type of file and the selected processing options.
        </Text>
      </Section>

      <Section title="3. File Upload Sources">
        <Text size="sm">You may upload files from the following sources:</Text>
        <List>
          <List.Item>
            <Text size="sm">FileUpload</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Audio</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Webcamera</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">ScreenCapture</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Box</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Dropbox</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Facebook</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">GoogleDrive</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">GooglePhotos</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Instagram</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">OneDrive</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Unsplash</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Url</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Zoom</Text>
          </List.Item>
        </List>
      </Section>

      <Section title="4. File Processing">
        <Text size="sm">
          We provide file processing services that may include file conversion,
          data extraction, or other forms of manipulation depending on the type
          of file uploaded. Processing may vary for each file type.
        </Text>
      </Section>

      <Section title="5. User Responsibilities">
        <Text size="sm">
          You are responsible for the content of the files you upload.
        </Text>
        <List>
          <List.Item>
            <Text size="sm">
              You must not upload files that contain harmful, offensive, or
              illegal content.
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              You agree to comply with all applicable laws regarding file
              uploading and processing.
            </Text>
          </List.Item>
        </List>
      </Section>

      <Section title="6. Data Storage">
        <Text size="sm">
          All uploaded and processed files are stored securely in our
          infrastructure using Amazon S3 or similar cloud storage services. By
          using our services, you consent to the storage and processing of your
          data in India.
        </Text>
      </Section>

      <Section title="7. User Data">
        <Text size="sm">
          We collect personal data to provide and improve our services. The
          information we collect includes:
        </Text>
        <List>
          <List.Item>
            <Text size="sm">Email</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">First Name</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Last Name</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Phone Number</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Avatar URL</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Provider (e.g., Google, Facebook)</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">Location</Text>
          </List.Item>
        </List>
        <Text size="sm">
          We may also collect non-personally identifiable data, such as usage
          patterns and technical data.
        </Text>
      </Section>

      <Section title="8. Data Deletion Requests">
        <Text size="sm">
          You may request to delete your personal data by emailing us at{" "}
          <Anchor size="sm" href="mailto:support@uploadpilot.app">
            support@uploadpilot.app
          </Anchor>
          . We will process your request within a reasonable timeframe.
        </Text>
      </Section>

      <Section title="9. Termination">
        <Text size="sm">
          We may suspend or terminate your access to the service if you violate
          any terms of service or if we determine that continued access could
          harm the system or other users.
        </Text>
      </Section>

      <Section title="10. Limitation of Liability">
        <Text size="sm">
          Our liability for any claim related to the use of the service is
          limited to the amount you have paid for the service in the past 12
          months, if any.
        </Text>
      </Section>

      <Section title="11. Changes to the Terms">
        <Text size="sm">
          We may update these terms from time to time. Any changes will be
          posted on this page, and we will notify you if required by law.
        </Text>
      </Section>
    </Container>
  );
};

const Section: React.FC<{ title: string; children: React.ReactNode }> = ({
  title,
  children,
}) => (
  <div>
    <Title order={3}>{title}</Title>
    {children}
    <Divider my="lg" />
  </div>
);

export default TermsOfService;
