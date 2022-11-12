import {AlertColor, Chip, Tooltip} from "@mui/material";
import React from "react";
import {AddonSource} from "../../../generated/graphql";

interface Props {
	source: AddonSource;
}

const AddonSourceChip: React.FC<Props> = ({source}): JSX.Element => {
	const colour = (): AlertColor => {
		switch (source) {
			case AddonSource.Official:
				return "success";
			case AddonSource.Platform:
				return "success";
			case AddonSource.Community:
				return "info";
			default:
			case AddonSource.Unknown:
				return "warning";
		}
	}

	const tip = (): string => {
		switch (source) {
			case AddonSource.Official:
				return "This addon is provided by KubeGlass.";
			case AddonSource.Platform:
				return "This addon is provided by your system administrators.";
			case AddonSource.Community:
				return "This addon is provided by the community.";
			default:
			case AddonSource.Unknown:
				return "This addon is provided by an unknown source (probably you!).";
		}
	}

	return <Tooltip
		title={tip()}
		arrow>
		<Chip
			sx={{ml: 1, maxHeight: 20, height: 20}}
			variant="outlined"
			label={source}
			color={colour()}
		/>
	</Tooltip>
}
export default AddonSourceChip;
