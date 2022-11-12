import React from "react";
import {
	Box,
	Button,
	Card,
	Dialog,
	DialogActions,
	DialogContent,
	DialogContentText,
	DialogTitle,
	Typography
} from "@mui/material";

interface Props {
	open: boolean;
	config: string;
	onClose: () => void;
}

const KubeConfigDialog: React.FC<Props> = ({open, config, onClose}): JSX.Element => {
	return <Dialog
		open={open}
		onClose={onClose}
		scroll="paper">
		<DialogTitle>
			Kubeconfig
		</DialogTitle>
		<DialogContent
			dividers>
			<DialogContentText>
				<Card
					variant="outlined"
					sx={{
						fontSize: 12,
						p: 1,
						pl: 2,
						color: "text.primary"
					}}
					component="pre">
					<code>
						{config}
					</code>
				</Card>
				<Box>
					<Typography
						color="textPrimary">
						Setup guide
					</Typography>
					<ol>
						<li>Create the <code>$HOME/.kube/config</code> file containing the above data.</li>
						<li>Run a command (e.g. <code>kubectl get pods</code>) to trigger the login.</li>
						<li>Follow the instructions in the browser page that opens.</li>
					</ol>
					For more information about client configuration, read the <a href="https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/" target="_blank" rel="noreferrer">Kubernetes documentation</a>.
				</Box>
			</DialogContentText>
		</DialogContent>
		<DialogActions>
			<Button
				onClick={onClose}>
				Close
			</Button>
		</DialogActions>
	</Dialog>
};
export default KubeConfigDialog;
