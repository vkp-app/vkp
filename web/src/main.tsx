import React from "react"
import ReactDOM from "react-dom/client"
import "./index.css"
import {BrowserRouter} from "react-router-dom";
import {ApolloProvider} from "@apollo/client";
import App from "./App"
import Client from "./graph";
import "typeface-roboto";

ReactDOM.createRoot(document.getElementById("root")!).render(
	<React.StrictMode>
		<BrowserRouter>
			<ApolloProvider client={Client}>
				<App/>
			</ApolloProvider>
		</BrowserRouter>
	</React.StrictMode>
);
