import React, {useMemo, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {
	Button,
	Card,
	CircularProgress,
	FormControlLabel,
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
import Icon from "@mdi/react";
import {mdiOpenInNew} from "@mdi/js";
import InlineError from "../alert/InlineError";
import {Cluster, useCanEditClusterQuery, useClusterQuery, useKubeConfigLazyQuery} from "../../generated/graphql";
import StandardLayout from "../layout/StandardLayout";
import BackButton from "../layout/BackButton";
import ClusterMetadataView from "./cluster/ClusterMetadataView";
import ClusterVersionIndicator from "./cluster/ClusterVersionIndicator";
import ClusterMetricsView from "./cluster/ClusterMetricsView";
import KubeConfigDialog from "./cluster/KubeConfigDialog";
import ClusterSettingsView from "./cluster/ClusterSettingsView";

const useStyles = makeStyles()((theme: Theme) => ({
	title: {
		fontFamily: "Figtree",
		fontSize: 30,
		fontWeight: 400,
		paddingBottom: theme.spacing(1)
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

	const canEditCluster = useCanEditClusterQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !tenantName || !clusterName
	});

	const installedAddons: string[] = useMemo(() => {
		if (data == null)
			return [];
		return data.clusterInstalledAddons.map(a => a.name);
	}, [data]);

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
		<BackButton
			title="Back to clusters"
			to={`/clusters/${tenantName}`}
		/>
		<Card
			sx={{p: 2, pt: 0}}>
			{!loading && error && <InlineError
				message="Unable to load cluster information"
				error={error}
			/>}
			<List>
				<ListItem>
					<ListItemText
						primaryTypographyProps={{className: classes.title, textTransform: "capitalize"}}
						primary={loading ? <Skeleton variant="text" height={50} width="40%"/> : data?.cluster.name}
						secondary={loading ? <Skeleton variant="text" height={25} width="20%"/> : <ClusterVersionIndicator
							showLabel
							version={data?.cluster.status.kubeVersion || ""}
							platform={data?.cluster.status.platformVersion}
						/>}
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
							disabled={loading || kubeConfig.loading}
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
							Manage and troubleshoot applications as well as manage the cluster itself through a web-based UI (requires installing a Dashboard addon).
						</span>}
					/>
					<ListItemSecondaryAction>
						<Button
							disabled={loading || !data?.cluster.status.webURL || installedAddons.find(i => i.startsWith(`${clusterName}-dashboard-`)) == null}
							startIcon={<Icon
								path={mdiOpenInNew}
								size={0.7}
							/>}
							href={`https://${data?.cluster.status.webURL}`}
							target="_blank">
							Open
						</Button>
					</ListItemSecondaryAction>
				</ListItem>
				<ListItem>
					<ListItemText
						primary="Addons"
						secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : <span>
							Pre-packaged applications, configuration and capability provided by us, your administrators and the community.
						</span>}
					/>
					<ListItemSecondaryAction>
						<Button
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
			clusterError={error}
			refresh={refresh}
		/>
		<ListSubheader>
			Advanced settings
		</ListSubheader>
		<ClusterSettingsView
			cluster={data?.cluster as Cluster | null}
			readOnly={!canEditCluster.data?.hasClusterAccess ?? true}
		/>
		{canEditCluster.data?.hasClusterAccess === true && <React.Fragment>
			<ListSubheader>
				Metadata
			</ListSubheader>
			<ClusterMetadataView
				cluster={data?.cluster as Cluster | null}
				loading={loading}
			/>
		</React.Fragment>}
		<KubeConfigDialog
			open={config !== ""}
			config={config}
			onClose={() => setConfig(() => "")}
		/>
	</StandardLayout>
}
export default ClusterView;
