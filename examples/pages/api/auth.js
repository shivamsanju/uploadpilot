export default function handler(req, res) {
  if (req.method !== "GET") {
    return res.status(405).json({ error: "Method not allowed" });
  }

  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    return res.status(401).json({ error: "Unauthorized" });
  }

  const token = authHeader.split(" ")[1];

  const SECRET_TOKEN = "mysecrettoken";

  if (token === SECRET_TOKEN) {
    return res.status(200).json({ message: "Authorized" });
  } else {
    return res.status(401).json({ error: "Unauthorized" });
  }
}
