import {Chip, Tooltip} from "@mui/material";
import React from "react";

interface Props {
	title?: string;
	label?: string;
}

const AddonChip: React.FC<Props> = ({title, label}): JSX.Element => {
	return <Tooltip
		title={title || "This addon has been installed."}>
		<Chip
			sx={{ml: 1, maxHeight: 20, height: 20, fontFamily: "Manrope", fontWeight: "bold", fontSize: 12}}
			variant="outlined"
			label={label || "Installed"}
			color="primary"
		/>
	</Tooltip>
}
export default AddonChip;
