import React, {useMemo, useState} from "react";
import {
	Alert,
	Box,
	Button,
	Card,
	CardHeader,
	FormGroup,
	Step,
	StepButton,
	Stepper,
	TextField,
	Theme,
	Typography
} from "@mui/material";
import {makeStyles} from "tss-react/mui";
import {useNavigate, useParams} from "react-router-dom";
import {useTheme} from "@mui/material/styles";
import {Circles, CircleSquare, Hexagon, Icon, TriangleSquareCircle} from "tabler-icons-react";
import StandardLayout from "../layout/StandardLayout";
import {Track, useCreateClusterMutation} from "../../generated/graphql";
import InlineError from "../alert/InlineError";

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

interface TrackItem {
	label: string;
	description: string;
	colour: string;
	icon: Icon;
	track: Track;
}

const TRACKS = (theme: Theme): TrackItem[] => [
	{
		label: "Stable",
		description: "Prioritise stability over new features. Updates are released on a less-frequent basis (excluding security updates) once faster tracks have extensively validated features.",
		colour: theme.palette.success.main,
		icon: Hexagon,
		track: Track.Stable
	},
	{
		label: "Regular (recommended)",
		description: "Access features reasonably soon after upstream release but after they have had some time to be validated in the Rapid and Beta tracks. Offers a balance between stability and new features and is recommended for most users.",
		colour: theme.palette.primary.main,
		icon: CircleSquare,
		track: Track.Regular
	},
	{
		label: "Rapid",
		description: "Access features as they are released upstream. The cluster will be updated frequently to stay on the latest version, and deliver new Kubernetes features and capability.",
		colour: theme.palette.warning.main,
		icon: Circles,
		track: Track.Rapid
	},
	{
		label: "Beta",
		description: "New versions and features hot off the press through the use of Release Candidate versions provided by the upstream. This track should not be used in a production capacity as normal SLAs can not be guaranteed.",
		colour: theme.palette.error.main,
		icon: TriangleSquareCircle,
		track: Track.Beta
	}
];

const nameRegex = new RegExp(/^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/);

const CreateCluster: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();
	const navigate = useNavigate();
	const theme = useTheme();
	const params = useParams();

	const tenantName = params["tenant"] || "";

	const [createCluster, {loading, error}] = useCreateClusterMutation();

	// local state
	const [track, setTrack] = useState<Track>(Track.Regular);
	const [trackIdx, setTrackIdx] = useState<number>(1);
	const [name, setName] = useState<string>("");
	const trackData = TRACKS(theme);

	const nameIsValid = useMemo(() => {
		return nameRegex.test(name);
	}, [name]);

	const onCreate = (): void => {
		createCluster({
			variables: {tenant: tenantName, input: {name, track}},
		}).then(r => {
			if (!r.errors) {
				navigate(`/clusters/${tenantName}/cluster/${r.data?.createCluster.name}`);
			}
		});
	}

	return <StandardLayout>
		<Card
			sx={{p: 2}}
			variant="outlined">
			<CardHeader
				title="Create cluster"
			/>
			{!loading && error && <InlineError
				message="Unable to lodge cluster creation request"
				error={error}
			/>}
			<Alert
				sx={{mb: 2, mt: 1}}
				severity="info">
				Some fields cannot be changed once the cluster has been created.
			</Alert>
			<FormGroup>
				<TextField
					label="Name"
					value={name}
					onChange={e => setName(() => e.target.value)}
					error={!nameIsValid}
					helperText="Must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character."
				/>
			</FormGroup>
			<FormGroup
				sx={{m: 1, mt: 4, mb: 4}}>
				<Stepper
					nonLinear
					alternativeLabel
					activeStep={trackIdx}>
					{trackData.map((t, i) => <Step
						key={t.label}>
						<StepButton
							icon={<t.icon
								color={t.colour}
							/>}
							onClick={() => {
								setTrack(() => t.track);
								setTrackIdx(() => i);
							}}>
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
					{trackData[trackIdx]?.description}
				</Typography>
			</Box>
			<Box
				sx={{pt: 2, display: "flex", float: "right"}}>
				<Button
					className={classes.button}
					variant="outlined"
					disabled={loading}
					onClick={() => navigate(-1)}>
					Cancel
				</Button>
				<Button
					className={classes.button}
					variant="outlined"
					disabled={!nameIsValid || loading}
					onClick={onCreate}>
					Create
				</Button>
			</Box>
		</Card>
	</StandardLayout>
}
export default CreateCluster;
