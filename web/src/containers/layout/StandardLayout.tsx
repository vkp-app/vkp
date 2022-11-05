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

import {Grid} from "@mui/material";
import React, {PropsWithChildren} from "react";
import {ErrorBoundary} from "react-error-boundary";
import Error from "../alert/Error";

interface Props extends PropsWithChildren {
	size?: "medium" | "small";
}

const StandardLayout: React.FC<Props> = ({size = "medium", children}): JSX.Element => {
	return (
		<Grid
			sx={{mt: 2, mb: 2}}
			container>
			<Grid
				item
				xs={false}
				sm={1}
				md={size === "medium" ? 3 : 4}
			/>
			<Grid
				item
				xs={12}
				sm={10}
				md={size === "medium" ? 6 : 4}>
				<ErrorBoundary
					fallbackRender={p => <Error props={p}/>}>
					{children}
				</ErrorBoundary>
			</Grid>
			<Grid
				item
				xs={false}
				sm={1}
				md={size === "medium" ? 3 : 4}
			/>
		</Grid>
	);
}
export default StandardLayout;
