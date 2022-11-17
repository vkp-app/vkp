import React, {useMemo, useState} from "react";
import {
	Alert,
	Box,
	Button,
	Card,
	CardHeader,
	Checkbox,
	Divider,
	FormGroup,
	FormLabel,
	List,
	ListItem,
	ListItemButton,
	ListItemIcon,
	ListItemSecondaryAction,
	ListItemText,
	Radio,
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
import Grid2 from "@mui/material/Unstable_Grid2";
import StandardLayout from "../layout/StandardLayout";
import {ReleaseTrack, useCreateClusterMutation} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import AddonChip from "./addon/AddonChip";

const useStyles = makeStyles()((theme: Theme) => ({
	button: {
		marginRight: theme.spacing(1)
	},
	formLabel: {
		fontFamily: "Manrope",
		fontWeight: "bold",
		fontSize: 13,
		paddingBottom: theme.spacing(1)
	}
}));

interface TrackItem {
	label: string;
	description: string;
	colour: string;
	icon: Icon;
	track: ReleaseTrack;
}

const TRACKS = (theme: Theme): TrackItem[] => [
	{
		label: "Stable",
		description: "Prioritise stability over new features. Updates are released on a less-frequent basis (excluding security updates) once faster tracks have extensively validated features.",
		colour: theme.palette.success.main,
		icon: Hexagon,
		track: ReleaseTrack.Stable
	},
	{
		label: "Regular (recommended)",
		description: "Access features reasonably soon after upstream release but after they have had some time to be validated in the Rapid and Beta tracks. Offers a balance between stability and new features and is recommended for most users.",
		colour: theme.palette.primary.main,
		icon: CircleSquare,
		track: ReleaseTrack.Regular
	},
	{
		label: "Rapid",
		description: "Access features as they are released upstream. The cluster will be updated frequently to stay on the latest version, and deliver new Kubernetes features and capability.",
		colour: theme.palette.warning.main,
		icon: Circles,
		track: ReleaseTrack.Rapid
	},
	{
		label: "Beta",
		description: "New versions and features hot off the press through the use of Release Candidate versions provided by the upstream. This track should not be used in a production capacity as normal SLAs can not be guaranteed.",
		colour: theme.palette.error.main,
		icon: TriangleSquareCircle,
		track: ReleaseTrack.Beta
	}
];

const nameRegex = new RegExp(/^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/);

interface ClusterTemplate {
	id: string;
	name: string;
	description: string;
	onClick: () => void;
}

const TEMPLATE_STANDARD = "standard";

const CreateCluster: React.FC = (): JSX.Element => {
	// hooks
	const {classes} = useStyles();
	const navigate = useNavigate();
	const theme = useTheme();
	const params = useParams();

	const tenantName = params["tenant"] || "";

	const [createCluster, {loading, error}] = useCreateClusterMutation();

	// local state
	const [track, setTrack] = useState<ReleaseTrack>(ReleaseTrack.Regular);
	const [trackIdx, setTrackIdx] = useState<number>(1);
	const [name, setName] = useState<string>("");
	const [ha, setHA] = useState<boolean>(false);

	const [selectedTemplate, setSelectedTemplate] = useState<string>(TEMPLATE_STANDARD);
	const trackData = TRACKS(theme);

	const nameIsValid = useMemo(() => {
		return nameRegex.test(name);
	}, [name]);

	const templates: ClusterTemplate[] = [
		{
			id: TEMPLATE_STANDARD,
			name: "Standard",
			description: "Best choice for general usage or if you are not sure what to choose.",
			onClick: () => {
				setHA(() => false);
				setTrack(() => ReleaseTrack.Regular);
				setTrackIdx(() => 1);
			}
		},
		{
			id: "production",
			name: "HA/Intensive",
			description: "Operators, CRDs or anything else that requires significant and uninterrupted interaction with the Kubernetes API.",
			onClick: () => {
				setHA(() => true);
				setTrack(() => ReleaseTrack.Stable);
				setTrackIdx(() => 0);
			}
		},
		{
			id: "pathfinder",
			name: "Pathfinder",
			description: "Our configuration for testing the latest and greatest.",
			onClick: () => {
				setHA(() => false);
				setTrack(() => ReleaseTrack.Beta);
				setTrackIdx(() => 3);
			}
		}
	];

	const onCreate = (): void => {
		createCluster({
			variables: {tenant: tenantName, input: {name, track, ha}},
		}).then(r => {
			if (!r.errors) {
				navigate(`/clusters/${tenantName}/cluster/${r.data?.createCluster.name}`);
			}
		});
	}

	const handleSelectTemplate = (t: ClusterTemplate): void => {
		setSelectedTemplate(() => t.id);
		t.onClick();
	}

	return <StandardLayout>
		<Card
			sx={{p: 2, pl: 0, pt: 0}}
			variant="outlined">
			<CardHeader
				title="Create a Kubernetes cluster"
				titleTypographyProps={{fontFamily: "Figtree"}}
			/>
			{!loading && error && <InlineError
				message="Unable to lodge cluster creation request"
				error={error}
			/>}
			<Divider/>
			<Grid2
				container>
				<Grid2
					xs={4}>
					<CardHeader
						title="Cluster templates"
						titleTypographyProps={{variant: "body1"}}
						subheader="Select a template with preconfigured settings, or customize to suit your needs."
						subheaderTypographyProps={{sx: {fontSize: 14, pt: 2}}}
					/>
					<List>
						{templates.map(t => <ListItem
							key={t.id}
							disableGutters
							disablePadding>
							<ListItemButton
								selected={t.id === selectedTemplate}
								disableRipple
								onClick={() => handleSelectTemplate(t)}>
								<ListItemIcon>
									<Radio
										value={t.id}
										checked={selectedTemplate === t.id}
										onClick={() => handleSelectTemplate(t)}
									/>
								</ListItemIcon>
								<ListItemText
									primary={t.name}
									secondary={t.description}
								/>
							</ListItemButton>
						</ListItem>)}
					</List>
				</Grid2>
				<Divider orientation="vertical" flexItem sx={{mr: "-1px"}}/>
				<Grid2
					xs={8}>
					<Box
						sx={{ml: 4, mr: 2, mt: 2}}>
						<Alert
							sx={{mb: 2, mt: 1}}
							severity="info">
							Some fields cannot be changed once the cluster has been created.
						</Alert>
						<FormGroup>
							<FormLabel
								className={classes.formLabel}>
								Name
							</FormLabel>
							<TextField
								size="small"
								value={name}
								onChange={e => setName(() => e.target.value)}
								error={!nameIsValid}
								helperText="Must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character."
							/>
						</FormGroup>
						<FormGroup
							sx={{mt: 4, mb: 4}}>
							<FormLabel
								className={classes.formLabel}>
								Release track
							</FormLabel>
							<Card
								sx={{p: 2}}
								variant="outlined">
								<Stepper
									nonLinear
									alternativeLabel
									activeStep={trackIdx}>
									{trackData.map((t, i) => <Step
										key={t.label}>
										<StepButton
											disableRipple
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
								<Box
									sx={{mt: 4}}>
									<Typography
										variant="body1"
										sx={{fontSize: 14}}>
										{trackData[trackIdx]?.description}
									</Typography>
								</Box>
							</Card>
						</FormGroup>
						<FormGroup>
							<FormLabel
								className={classes.formLabel}>
								Advanced settings
							</FormLabel>
							<List>
								<ListItem
									disablePadding>
									<ListItemText
										primary={<span>
											High availability <AddonChip label="BETA" title="This option is available for use but may not be production-ready."/>
										</span>}
										secondary="Increases the number of replicas for Kubernetes components."
									/>
									<ListItemSecondaryAction>
										<Checkbox
											checked={ha}
											onChange={(e, checked) => setHA(() => checked)}
										/>
									</ListItemSecondaryAction>
								</ListItem>
							</List>
						</FormGroup>
						<Box
							sx={{pt: 2, display: "flex", float: "right"}}>
							<Button
								className={classes.button}
								disabled={loading}
								onClick={() => navigate(-1)}>
							Cancel
							</Button>
							<Button
								className={classes.button}
								disabled={!nameIsValid || loading}
								onClick={onCreate}>
							Create
							</Button>
						</Box>
					</Box>
				</Grid2>
			</Grid2>
		</Card>
	</StandardLayout>
}
export default CreateCluster;
