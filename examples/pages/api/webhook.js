export default function handler(req, res) {
    if (req.method !== 'POST') {
        return res.status(405).json({ error: 'Method not allowed' });
    }

    console.log("headers: ", req.headers);
    const authHeader = req.headers.secret;
    const SECRET_TOKEN = 'webhooksecret';


    if (authHeader === SECRET_TOKEN) {
        const body = req.body;
        console.log("Recieved a webhook: ", body);
        return res.status(200).json({ message: 'Authorized' });
    } else {
        return res.status(401).json({ error: 'Unauthorized' });
    }
}
