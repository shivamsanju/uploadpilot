import { Anchor, Container, Divider, List, Text, Title } from '@mantine/core';
import React from 'react';

const PrivacyPolicy: React.FC = () => {
  return (
    <Container size="sm" mt="lg">
      <Title order={1}>Privacy Policy</Title>
      <Text size="sm" c="dimmed">
        Effective Date: [01/01/2025]
      </Text>

      <Divider my="lg" />

      <Section title="1. Introduction">
        <Text size="sm">
          UploadPilot values your privacy and is committed to protecting your
          personal data. This Privacy Policy explains how we collect, use,
          store, and protect your personal information when you use our
          services.
        </Text>
      </Section>

      <Section title="2. Information We Collect">
        <Text size="sm">
          We collect the following personal data when you register and use our
          services:
        </Text>
        <List>
          <List.Item>
            <Text size="sm">
              Personal Information: Name, email address, phone number, location,
              avatar URL, description, and user ID.
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              Usage Data: Data related to your interactions with our service,
              including upload history, processing requests, and limited file
              metadata like size and type.
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              Technical Data: Information about the devices and networks you use
              to access our services.
            </Text>
          </List.Item>
        </List>
      </Section>

      <Section title="3. Use of Your Information">
        <Text size="sm">
          We use your personal data for the following purposes:
        </Text>
        <List>
          <List.Item>
            <Text size="sm">To provide and maintain our services</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">To process files you upload</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">To improve our platform and services</Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              To communicate with you regarding updates, billing, or support
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">To respond to data deletion requests</Text>
          </List.Item>
        </List>
      </Section>

      <Section title="4. Data Storage and Security">
        <Text size="sm">
          All data, including personal data and uploaded files, are stored
          securely in Amazon S3 or similar cloud storage services. We implement
          appropriate technical and organizational measures to protect your data
          from unauthorized access or loss.
        </Text>
      </Section>

      <Section title="5. Data Sharing and Third-Party Services">
        <Text size="sm">
          We do not share your personal data with third parties unless required
          by law or for the purpose of providing our services (e.g., third-party
          API integrations). We do not sell your personal data to any third
          parties.
        </Text>
      </Section>

      <Section title="6. User Rights">
        <Text size="sm">
          As a user, you have the following rights regarding your personal data:
        </Text>
        <List>
          <List.Item>
            <Text size="sm">
              Right to access: You can request a copy of your data at any time.
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              Right to correct: You can request corrections to any inaccurate or
              incomplete information.
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              Right to delete: You can request to delete your data by emailing
              support@uploadpilot.app.
            </Text>
          </List.Item>
          <List.Item>
            <Text size="sm">
              Right to object: You can object to the processing of your data in
              certain circumstances.
            </Text>
          </List.Item>
        </List>
      </Section>

      <Section title="7. International Data Transfers">
        <Text size="sm">
          Our services are based in India, and by using our platform, you
          consent to the transfer of your data to India. We take appropriate
          steps to ensure your data is protected in accordance with this policy.
        </Text>
      </Section>

      <Section title="8. Cookies and Tracking Technologies">
        <Text size="sm">
          We may use cookies and other tracking technologies to enhance your
          experience on our platform. You can control cookies through your
          browser settings.
        </Text>
      </Section>

      <Section title="9. Changes to the Privacy Policy">
        <Text size="sm">
          We may update this Privacy Policy from time to time. Any changes will
          be posted on this page, and we will notify you if required by law.
        </Text>
      </Section>

      <Section title="10. Contact Us">
        <Text size="sm">
          If you have any questions or concerns about this Privacy Policy or our
          data practices, please contact us at{' '}
          <Anchor size="sm" href="mailto:support@uploadpilot.app">
            support@uploadpilot.app
          </Anchor>
          .
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

export default PrivacyPolicy;
