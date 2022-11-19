import {FormControlLabel, IconButton, ListSubheader, Typography} from "@mui/material";
import {Link} from "react-router-dom";
import React from "react";
import Icon from "@mdi/react";
import {mdiArrowLeft} from "@mdi/js";

interface Props {
	title: string;
	to: string;
}

const BackButton: React.FC<Props> = ({title, to}): JSX.Element => {
	return <ListSubheader
		sx={{display: "flex", alignItem: "center", mb: 1, mt: 1}}>
		<FormControlLabel
			sx={{cursor: "initial", ml: -1}}
			disableTypography
			control={<IconButton
				size="small"
				component={Link}
				to={to}>
				<Icon
					path={mdiArrowLeft}
					size={0.75}
				/>
			</IconButton>}
			label={<Typography
				sx={{ml: 1, fontSize: 14, fontWeight: 500}}
				variant="body1">
				{title}
			</Typography>}/>
	</ListSubheader>
}
export default BackButton;
