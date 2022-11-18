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
	List,
	ListItem,
	ListItemSecondaryAction,
	ListItemText
} from "@mui/material";
import {Link, useNavigate} from "react-router-dom";
import {Cluster, useDeleteClusterMutation} from "../../../generated/graphql";
import InlineError from "../../alert/InlineError";

interface Props {
	cluster: Cluster | null;
	readOnly: boolean;
}

const ClusterSettingsView: React.FC<Props> = ({cluster, readOnly}): JSX.Element => {
	// hooks
	const navigate = useNavigate();
	const [deleteCluster, {loading, error}] = useDeleteClusterMutation();

	const [showDelete, setShowDelete] = useState<boolean>(false);

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

	return <Card
		sx={{p: 2}}>
		<List>
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
			onClose={() => setShowDelete(() => false)}>
			<DialogTitle>
				Delete cluster
			</DialogTitle>
			<DialogContent>
				All cluster and workload resources will be removed.
				<b>If you have running applications, they will be permanently deleted.</b>
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
					disabled={loading}
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
