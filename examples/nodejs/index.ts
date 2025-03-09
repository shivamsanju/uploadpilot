import fs from "fs";
import path from "path";
import { NodeFile, Uploader } from "./uploader";

const uploader = new Uploader(
  "a3660691-0412-4da7-a537-e63f500899e3",
  "a1e4d1f9-7d0a-4d85-8b25-ea8f90267dc7",
  "up-PoDTH92DEDKN78WD20250316165824"
);

async function main(): Promise<void> {
  try {
    const filePath: string = path.resolve("./files/a.txt"); // Replace with actual file path
    const fileBuffer: Buffer = fs.readFileSync(filePath);

    const file = new NodeFile(
      fileBuffer,
      path.basename(filePath),
      "application/octet-stream"
    );

    await uploader.uploadMultiple([file], {
      file_name: "a.txt",
    });

    console.log("File uploaded successfully");
  } catch (error) {
    console.error(
      "Error uploading file:",
      error instanceof Error ? error.message : error
    );
  }
}

main();
