import Head from "next/head";
import { Uploader } from "uppy-react"

export default function Home() {
  return (
    <>
      <Head>
        <title>UploadPilot Next.js Demo</title>
        <meta name="description" content="UploadPilot Next.js Demo" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh" }}>
        <Uploader
          workspaceId="678e470bb0470cf5d9e6aff4"
          backendEndpoint="http://localhost:8080"
          height={400}
          width={350}
          theme={"auto"}
          metadata={{ "user": "johndoe@example.com" }}
          headers={{ "Authorization": "Bearer mysecrettoken" }}
          note="Upload your files here"
        />
      </main>
    </>
  );
}
