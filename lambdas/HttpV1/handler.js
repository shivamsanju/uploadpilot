const axios = require("axios");

exports.handler = async (event) => {
  try {
    const url = event.url || "";
    const method = event.method || "GET";
    const headers = event.headers || {};
    const body = event.body || null;

    const response = await axios({
      method: method.toUpperCase(),
      url: url,
      headers: headers,
      data: body,
    });

    return {
      statusCode: response.status,
      headers: response.headers,
      body: response.data,
    };
  } catch (error) {
    return {
      statusCode: error.response ? error.response.status : 500,
      body: error.message,
    };
  }
};
