import {
  IconAi,
  IconBadge,
  IconBrandGit,
  IconCloud,
  IconDeviceFloppy,
  IconDualScreen,
  IconFileInfo,
  IconImageInPicture,
  IconJoinRound,
  IconLock,
  IconLockAccessOff,
  IconMail,
  IconSearch,
  IconShieldCheck,
  IconSortDescendingSmallBig,
  IconStopwatch,
  IconTag,
  IconTagsOff,
  IconTextRecognition,
  IconTool,
  IconTransformPoint,
} from "@tabler/icons-react";
import { Paper, SimpleGrid, Text, useMantineTheme } from "@mantine/core";
import classes from "./Utils.module.css";

const mockdata = [
  {
    title: "Virus scanning",
    description: "Automatically scan uploaded files for viruses and malware.",
    icon: IconShieldCheck,
  },
  {
    title: "Metadata storage",
    description:
      "Store upload metadata such as file size, type, and creation date on elasticsearch.",
    icon: IconFileInfo,
  },
  {
    title: "Convert file format",
    description:
      "Convert uploaded files into different formats like PDF, JPG, or DOCX.",
    icon: IconTransformPoint,
  },
  {
    title: "Copy to cloud storage",
    description:
      "Copy the imports to a cloud storage service like S3, GCS or Azure",
    icon: IconCloud,
  },
  {
    title: "Reduce file size",
    description: "Compress files before they are uploaded to your workspace",
    icon: IconSortDescendingSmallBig,
  },
  {
    title: "OCR processing",
    description:
      "Extract text from images or scanned documents using Optical Character Recognition.",
    icon: IconTextRecognition,
  },
  {
    title: "Content tagging",
    description:
      "Automatically tag files based on their content for easy categorization.",
    icon: IconTag,
  },
  {
    title: "AI-based classification",
    description:
      "Classify files into predefined categories using AI and machine learning.",
    icon: IconAi,
  },
  {
    title: "Duplicate detection",
    description: "Identify and flag duplicate files in the system.",
    icon: IconDualScreen,
  },
  {
    title: "Encryption",
    description: "Encrypt files to secure sensitive information.",
    icon: IconLock,
  },
  {
    title: "File preview generation",
    description: "Generate previews for images, videos, and documents.",
    icon: IconImageInPicture,
  },
  {
    title: "Transcoding",
    description:
      "Transcode videos or audio files into different resolutions or bitrates.",
    icon: IconTool,
  },
  {
    title: "Automatic backup",
    description: "Backup uploaded files to a secure storage location.",
    icon: IconDeviceFloppy,
  },
  {
    title: "Email notifications",
    description:
      "Send notifications upon successful processing of uploaded files.",
    icon: IconMail,
  },
  {
    title: "Version control",
    description: "Maintain a history of changes and updates to files.",
    icon: IconBrandGit,
  },
  {
    title: "File expiration",
    description:
      "Set expiration dates for files to automatically delete old data.",
    icon: IconStopwatch,
  },
  {
    title: "Watermarking",
    description:
      "Add watermarks to images or documents for branding or copyright protection.",
    icon: IconTagsOff,
  },
  {
    title: "Audit logging",
    description:
      "Log all actions taken on files for compliance and monitoring.",
    icon: IconBadge,
  },
  {
    title: "Access control",
    description:
      "Define permissions to control who can view or edit uploaded files.",
    icon: IconLockAccessOff,
  },
  {
    title: "Integration with third-party apps",
    description:
      "Sync files with third-party services like Dropbox, Google Drive, or Slack.",
    icon: IconJoinRound,
  },
  {
    title: "Searchable file repository",
    description: "Index files to enable fast and efficient searching.",
    icon: IconSearch,
  },
  {
    title: "Thumbnail generation",
    description: "Generate thumbnails for uploaded images and videos.",
    icon: IconImageInPicture,
  },
];

export const ToolsGrid = () => {
  const theme = useMantineTheme();
  const features = mockdata.map((feature) => (
    <Paper key={feature.title} radius="md" p="xl">
      <feature.icon size={50} stroke={2} color={theme.colors.appcolor[6]} />
      <Text fz="lg" fw={500} className={classes.cardTitle} mt="md">
        {feature.title}
      </Text>
      <Text fz="sm" c="dimmed" mt="sm">
        {feature.description}
      </Text>
    </Paper>
  ));

  return (
    <SimpleGrid cols={{ base: 1, md: 3, xl: 4 }} spacing="xl">
      {features}
    </SimpleGrid>
  );
};
