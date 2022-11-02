/*
 *    Copyright 2022 Django Cass
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

import React from "react";
import {Alert, AlertTitle, Button} from "@mui/material";
import {FallbackProps} from "react-error-boundary";

interface ErrorProps {
	props: FallbackProps;
}

const Error: React.FC<ErrorProps> = ({props}): JSX.Element => {
	return (
		<div>
			<Alert
				severity="error">
				<AlertTitle>
					Something went wrong
				</AlertTitle>
				This component was unable to render properly. Sorry about that!
				<br/>
				<br/>
				<code>{props.error.name}: {props.error.message}</code>
				<br/>
				<Button
					sx={{mt: 2, textTransform: "none"}}
					variant="contained"
					color="primary"
					onClick={() => window.location.reload()}>
					Reload page
				</Button>
			</Alert>
		</div>
	);
}
export default Error;
