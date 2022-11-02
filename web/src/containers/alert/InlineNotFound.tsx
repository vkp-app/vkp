import {Box, Typography} from "@mui/material";
import React from "react";

interface Props {
	title?: string;
	subtitle?: string;
}

const InlineNotFound: React.FC<Props> = ({title, subtitle}): JSX.Element => <Box sx={{m: 1, display: "flex", alignItems: "center", justifyContent: "center"}}>
	<div>
		<Typography
			sx={{display: "block", fontWeight: 500}}
			align="center"
			color="textSecondary">
			{title || "No resources"}
		</Typography>
		<Typography
			sx={{display: "block", fontSize: 14}}
			align="center"
			color="textSecondary">
			{subtitle || "Nothing could be found."}
		</Typography>
	</div>
</Box>;
export default InlineNotFound;
