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

import {Alert} from "@mui/material";
import React from "react";
import {ApolloError} from "@apollo/client";
import InlineNotFound from "./InlineNotFound";

interface Props {
	error: ApolloError | undefined;
	message?: string;
}

const InlineError: React.FC<Props> = ({error, message}): JSX.Element | null => {
	if (!error)
		return null;

	switch (error?.message) {
		case "unauthorised":
		case "unauthorized":
			return <InlineNotFound
				title="Not logged in."
				subtitle="Please login to view this resource."
			/>
		default:
			return <Alert
				severity="error">
				{message || "Something went wrong"}
				{error && ` (${error?.message})`}
			</Alert>
	}
}
export default InlineError;
