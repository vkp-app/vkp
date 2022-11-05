import React, {useState} from "react";
import {
	Alert,
	Box,
	Button,
	Card,
	FormGroup,
	ListSubheader,
	Step,
	StepButton,
	Stepper,
	TextField,
	Theme,
	Typography
} from "@mui/material";
import {makeStyles} from "tss-react/mui";
import {useNavigate} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";

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

const TRACKS = [
	{
		label: "Stable",
		description: "Prioritise stability over new features. Updates are released on a less-frequent basis (excluding security updates) once faster tracks have extensively validated features."
	},
	{
		label: "Regular (recommended)",
		description: "Access features reasonably soon after upstream release but after they have had some time to be validated in the Rapid and Beta tracks. Offers a balance between stability and new features and is recommended for most users."
	},
	{
		label: "Rapid",
		description: "Access features as they are released upstream. The cluster will be updated frequently to stay on the latest version, and deliver new Kubernetes features and capability."
	},
	{
		label: "Beta",
		description: "New versions and features hot off the press through the use of Release Candidate versions provided by the upstream. This track should not be used in a production capacity as normal SLAs can not be guaranteed."
	}
];

const CreateCluster: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();
	const navigate = useNavigate();

	// local state
	const [track, setTrack] = useState<number>(1);

	return <StandardLayout>
		<ListSubheader>
			Create cluster
		</ListSubheader>
		<Card
			sx={{p: 2}}
			variant="outlined">
			<Alert
				sx={{mb: 2}}
				severity="info">
				Some fields cannot be changed once the cluster has been created.
			</Alert>
			<FormGroup>
				<TextField
					label="Name"
				/>
			</FormGroup>
			<FormGroup
				sx={{m: 1, mt: 4, mb: 4}}>
				<Stepper
					nonLinear
					activeStep={track}>
					{TRACKS.map((t, i) => <Step
						key={t.label}>
						<StepButton
							onClick={() => setTrack(() => i)}>
							{t.label}
						</StepButton>
					</Step>)}
				</Stepper>
			</FormGroup>
			<Box
				sx={{mt: 4}}>
				<Typography
					variant="body1"
					sx={{fontSize: 14}}>
					{TRACKS[track]?.description}
				</Typography>
			</Box>
			<Box
				sx={{pt: 2, display: "flex", float: "right"}}>
				<Button
					className={classes.button}
					variant="outlined"
					onClick={() => navigate(-1)}>
					Cancel
				</Button>
				<Button
					className={classes.button}
					variant="outlined"
					disabled>
					Create
				</Button>
			</Box>
		</Card>
	</StandardLayout>
}
export default CreateCluster;
