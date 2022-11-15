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
import {useNavigate} from "react-router-dom";
import {makeStyles} from "tss-react/mui";
import {Cluster, useDeleteClusterMutation} from "../../../generated/graphql";
import InlineError from "../../alert/InlineError";

const useStyles = makeStyles()(() => ({
	button: {
		fontFamily: "Manrope",
		fontWeight: 600,
		fontSize: 13,
		textTransform: "none",
		minHeight: 28,
		height: 28
	}
}));

interface Props {
	cluster: Cluster | null;
}

const ClusterSettingsView: React.FC<Props> = ({cluster}): JSX.Element => {
	// hooks
	const {classes} = useStyles();
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
		variant="outlined"
		sx={{p: 2}}>
		<List>
			<ListItem>
				<ListItemText
					primary="Delete cluster"
					secondary="Permanently remove this cluster and all its resources."
				/>
				<ListItemSecondaryAction>
					<Button
						className={classes.button}
						color="error"
						variant="outlined"
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
					className={classes.button}
					variant="outlined"
					onClick={() => setShowDelete(() => false)}
					disabled={loading}>
					Cancel
				</Button>
				<Button
					className={classes.button}
					variant="outlined"
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
