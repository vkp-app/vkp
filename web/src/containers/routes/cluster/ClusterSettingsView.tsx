import React, {useState} from "react";
import {
	Box,
	Button,
	Card,
	CircularProgress,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	FormControlLabel,
	List,
	ListItem,
	ListItemSecondaryAction,
	ListItemText,
	Switch
} from "@mui/material";
import {Link, useNavigate} from "react-router-dom";
import {formatDistance} from "date-fns";
import cronstrue from "cronstrue";
import {Cluster, MaintenanceWindow, useDeleteClusterMutation} from "../../../generated/graphql";
import InlineError from "../../alert/InlineError";

interface Props {
	cluster: Cluster | null;
	maintenanceWindow: MaintenanceWindow | null;
	readOnly: boolean;
}

const ClusterSettingsView: React.FC<Props> = ({cluster, maintenanceWindow, readOnly}): JSX.Element => {
	// hooks
	const navigate = useNavigate();
	const [deleteCluster, {loading, error}] = useDeleteClusterMutation();

	const [showDelete, setShowDelete] = useState<boolean>(false);
	const [allowDelete, setAllowDelete] = useState<boolean>(false);

	const onDeleteCluster = (): void => {
		if (!cluster)
			return;
		deleteCluster({
			variables: {tenant: cluster.tenant, cluster: cluster.name}
		}).then(r => {
			if (r.data != null) {
				navigate(`/clusters/${cluster.tenant}`);
			}
		});
	}

	const handleClose = (): void => {
		setShowDelete(() => false);
		setAllowDelete(() => false);
	}

	const cronExpr = cronstrue.toString(maintenanceWindow?.schedule || "* * * * *", {throwExceptionOnParseError: false});

	return <Card
		sx={{p: 2}}>
		<List>
			<ListItem>
				<ListItemText
					primary="Maintenance automation"
					secondary={<span>
						Schedule: {maintenanceWindow?.schedule === "* * * * *" ? "at any time" : cronExpr.charAt(0).toLocaleLowerCase() + cronExpr.slice(1)}<br/>
						{maintenanceWindow?.schedule !== "* * * * *" && <React.Fragment>
							Next window: in {formatDistance((maintenanceWindow?.next || 0) * 1000, Date.now())}
						</React.Fragment>}
					</span>}
				/>
				<ListItemSecondaryAction>
					<Button
						disabled={readOnly || cluster == null}
						component={Link}
						to={`/clusters/${cluster?.tenant}/cluster/${cluster?.name}/-/maintenance`}>
						Open
					</Button>
				</ListItemSecondaryAction>
			</ListItem>
			<ListItem>
				<ListItemText
					primary="Permissions"
					secondary="Control who can access this cluster and what they can do."
				/>
				<ListItemSecondaryAction>
					<Button
						disabled={readOnly || cluster == null}
						component={Link}
						to={`/clusters/${cluster?.tenant}/cluster/${cluster?.name}/-/accessors`}>
						Open
					</Button>
				</ListItemSecondaryAction>
			</ListItem>
			<ListItem>
				<ListItemText
					primary="Delete cluster"
					secondary="Permanently remove this cluster and all its resources."
				/>
				<ListItemSecondaryAction>
					<Button
						color="error"
						disabled={readOnly || cluster == null}
						onClick={() => setShowDelete(() => true)}>
						Delete
					</Button>
				</ListItemSecondaryAction>
			</ListItem>
		</List>
		<Dialog
			open={showDelete}
			onClose={handleClose}>
			<DialogTitle>
				Delete cluster
			</DialogTitle>
			<DialogContent>
				All cluster and workload resources will be removed.
				<b>If you have running applications, they will be permanently deleted.</b>
				<br/>
				<FormControlLabel
					value={allowDelete}
					onChange={(event, checked) => setAllowDelete(() => checked)}
					control={<Switch color="error"/>}
					label="I understand"
				/>
				{!loading && <Box
					sx={{mt: 1}}>
					<InlineError
						message="Unable to delete cluster"
						error={error}
					/>
				</Box>}
			</DialogContent>
			<DialogActions>
				<Button
					onClick={() => setShowDelete(() => false)}
					disabled={loading}>
					Cancel
				</Button>
				<Button
					onClick={() => onDeleteCluster()}
					color="error"
					disabled={loading || !allowDelete}
					startIcon={loading && <CircularProgress
						size={16}
						color="error"
					/>}>
					Delete cluster and resources
				</Button>
			</DialogActions>
		</Dialog>
	</Card>
}
export default ClusterSettingsView;
