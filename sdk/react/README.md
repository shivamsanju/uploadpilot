## How to use

```ts
import { Uploader } from "uppy-react";
import "uppy-react/dist/style.css";

export default function Home() {
  return (
    <SomeModal>
      <Uploader
        workspaceId="c640d2ad-9f99-41fa-a348-d78a60f782d7"
        uploadEndpoint="http://localhost:8081"
        height={400}
        width={350}
        theme={"auto"}
        metadata={{ user: "johndoe@example.com" }}
        headers={{ Authorization: "Bearer mysecrettoken" }}
        note="Upload your files here"
      />
    </SomeModal>
  );
}
```
