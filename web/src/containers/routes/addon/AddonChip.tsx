import {Chip, Tooltip} from "@mui/material";
import React from "react";

const AddonChip: React.FC = (): JSX.Element => {
	return <Tooltip
		title="This addon has been installed.">
		<Chip
			sx={{ml: 1, maxHeight: 20, height: 20}}
			variant="outlined"
			label="Installed"
			color="primary"
		/>
	</Tooltip>
}
export default AddonChip;
