import React, {useState} from "react";
import {Link, useParams} from "react-router-dom";
import {
	Button,
	Card,
	CircularProgress,
	FormControlLabel,
	IconButton,
	List,
	ListItem,
	ListItemSecondaryAction,
	ListItemText,
	ListSubheader,
	Skeleton,
	Switch,
	Theme
} from "@mui/material";
import {makeStyles} from "tss-react/mui";
import {ArrowLeft, ExternalLink} from "tabler-icons-react";
import InlineError from "../alert/InlineError";
import {Cluster, useClusterQuery, useKubeConfigLazyQuery} from "../../generated/graphql";
import StandardLayout from "../layout/StandardLayout";
import ClusterMetadataView from "./cluster/ClusterMetadataView";
import ClusterVersionIndicator from "./cluster/ClusterVersionIndicator";
import ClusterMetricsView from "./cluster/ClusterMetricsView";
import KubeConfigDialog from "./cluster/KubeConfigDialog";

const useStyles = makeStyles()((theme: Theme) => ({
	title: {
		fontFamily: "Manrope",
		fontSize: 24,
		fontWeight: 500,
		paddingBottom: theme.spacing(1)
	},
	button: {
		fontFamily: "Manrope",
		fontWeight: 600,
		fontSize: 12,
		textTransform: "none",
		minHeight: 24,
		height: 24
	}
}));

const ClusterView: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();
	const {classes} = useStyles();

	const clusterName = params["name"] || "";
	const tenantName = params["tenant"] || "";

	// local state
	const [refresh, setRefresh] = useState<boolean>(false);
	const [config, setConfig] = useState<string>("");

	const {data, loading, error} = useClusterQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !clusterName || !tenantName,
		pollInterval: refresh ? 30_000 : 0
	});

	const [renderKubeConfig, kubeConfig] = useKubeConfigLazyQuery();

	const onRenderConfig = (): void => {
		renderKubeConfig({
			variables: {cluster: clusterName, tenant: tenantName}
		}).then(r => {
			if (r.data != null) {
				setConfig(() => r.data?.renderKubeconfig || "");
			}
		});
	}

	return <StandardLayout>
		<ListSubheader
			sx={{display: "flex", alignItem: "center"}}>
			<IconButton
				size="small"
				centerRipple={false}
				component={Link}
				to={`/clusters/${tenantName}`}>
				<ArrowLeft
					size={18}
				/>
			</IconButton>
			Back to clusters
		</ListSubheader>
		<Card
			variant="outlined"
			sx={{p: 2}}>
			{!loading && error && <InlineError
				message="Unable to load cluster information"
				error={error}
			/>}
			<List>
				<ListItem>
					<ListItemText
						primaryTypographyProps={{className: classes.title}}
						primary={loading ? <Skeleton variant="text" height={30} width="40%"/> : data?.cluster.name}
						secondary={loading ? <Skeleton variant="text" height={20} width="20%"/> : <ClusterVersionIndicator showLabel version={data?.cluster.status.kubeVersion || ""}/>}
					/>
					<ListItemSecondaryAction>
						<FormControlLabel
							control={<Switch/>}
							label="Auto-refresh"
							labelPlacement="start"
							checked={refresh}
							onChange={(e, checked) => setRefresh(() => checked)}
						/>
					</ListItemSecondaryAction>
				</ListItem>
				<ListItem>
					<ListItemText
						primary="Client setup"
						secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : "Generate a Kubeconfig that can be used to access the cluster."}
					/>
					<ListItemSecondaryAction>
						<Button
							className={classes.button}
							disabled={loading || kubeConfig.loading}
							variant="outlined"
							startIcon={kubeConfig.loading ? <CircularProgress size={14}/> : undefined}
							onClick={onRenderConfig}>
							Open
						</Button>
					</ListItemSecondaryAction>
				</ListItem>
				<ListItem>
					<ListItemText
						primary="Dashboard"
						secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : <span>
							Manage and troubleshoot applications as well as manage the cluster itself through a web-based UI.
						</span>}
					/>
					<ListItemSecondaryAction>
						<Button
							className={classes.button}
							disabled={loading || !data?.cluster.status.webURL}
							variant="outlined"
							startIcon={<ExternalLink
								size={18}
							/>}
							href={`https://${data?.cluster.status.webURL}`}
							target="_blank">
							Open
						</Button>
					</ListItemSecondaryAction>
				</ListItem>
				<ListItem>
					<ListItemText
						primary="Supply-chain security"
						secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : <span>
							Verify and enforce container image supply-chain security.
						</span>}
					/>
					<ListItemSecondaryAction>
						<Button
							className={classes.button}
							variant="outlined"
							component={Link}
							to={`/clusters/${tenantName}/cluster/${clusterName}/-/addons`}>
							Open
						</Button>
					</ListItemSecondaryAction>
				</ListItem>
			</List>
		</Card>
		<ListSubheader>
			Metrics
		</ListSubheader>
		<ClusterMetricsView
			cluster={data?.cluster as Cluster | null}
			loading={loading}
			refresh={refresh}
		/>
		<ListSubheader>
			Metadata
		</ListSubheader>
		<ClusterMetadataView
			cluster={data?.cluster as Cluster | null}
			loading={loading}
		/>
		<KubeConfigDialog
			open={config !== ""}
			config={config}
			onClose={() => setConfig(() => "")}
		/>
	</StandardLayout>
}
export default ClusterView;
