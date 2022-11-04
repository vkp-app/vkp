import React from "react";
import StandardLayout from "../layout/StandardLayout";
import {Alert, Box, Button, Card, FormGroup, ListSubheader, TextField, Theme} from "@mui/material";
import {makeStyles} from "tss-react/mui";
import {useNavigate} from "react-router-dom";

const useStyles = makeStyles()((theme: Theme) => ({
	button: {
		fontFamily: "Manrope",
		fontWeight: 600,
		fontSize: 13,
		textTransform: "none",
		minHeight: 24,
		height: 24,
		marginRight: theme.spacing(1)
	}
}));

const CreateCluster: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();
	const navigate = useNavigate();

	return <StandardLayout>
		<ListSubheader>
			Create cluster
		</ListSubheader>
		<Card
			sx={{p: 2}}
			variant={"outlined"}>
			<Alert
				sx={{mb: 2}}
				severity={"info"}>
				Some fields cannot be changed once the cluster has been created.
			</Alert>
			<FormGroup>
				<TextField
					label={"Name"}
				/>
			</FormGroup>
			<Box
				sx={{pt: 2, display: "flex", float: "right"}}>
				<Button
					className={classes.button}
					variant={"outlined"}
					onClick={() => navigate(-1)}>
					Cancel
				</Button>
				<Button
					className={classes.button}
					variant={"outlined"}
					disabled>
					Create
				</Button>
			</Box>
		</Card>
	</StandardLayout>
}
export default CreateCluster;
