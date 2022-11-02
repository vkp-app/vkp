import React from "react";
import {Alert, AlertTitle} from "@mui/material";
import StandardLayout from "../layout/StandardLayout";

const NotFound: React.FC = (): JSX.Element => {
	return (
		<StandardLayout>
			<div>
				<Alert
					severity="error">
					<AlertTitle>
						Page not found
					</AlertTitle>
					This page cannot be found, or the server refused to disclose it.
				</Alert>
			</div>
		</StandardLayout>
	);
}
export default NotFound;
