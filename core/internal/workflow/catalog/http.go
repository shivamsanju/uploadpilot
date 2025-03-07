package catalog

var HTTP_V_01 = &ActivityMetadata{
	Name:        "HTTP_V_01",
	DisplayName: "Sends Http Request",
	Description: `
This Activity sends an HTTP request to a specified endpoint. 
The request can include headers, parameters, and a body. 
The response can be captured for further processing or analysis.
  `,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: HTTP_V_01
    with:
      url: "https://api.restful-api.dev/objects"
      method: "POST"
      headers:
        Content-Type: "application/json"
      body:
        name: "Apple MacBook Pro 16"
        data:
          year: 2019
          price: 1849.99
          CPU model: "Intel Core i9"
          Hard disk size: "1 TB"
`,
}
