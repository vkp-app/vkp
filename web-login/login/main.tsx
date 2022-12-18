import React from "react";
import ReactDOM from "react-dom/client";
import App from "../src/App";
import Login from "./Login";
import "../index.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<React.StrictMode>
		<App>
			<Login/>
		</App>
	</React.StrictMode>,
);
