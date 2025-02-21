import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

const TokenHandler = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get("uploadpilottoken");

    if (token) {
      console.log("Token found:", token);
      localStorage.setItem("uploadpilottoken", token);
      navigate("/");
    }
  }, [navigate]);

  return <></>;
};

export default TokenHandler;
