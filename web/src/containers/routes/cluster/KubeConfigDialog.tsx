import React, {useState} from "react";
import {
	Alert,
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
import Icon from "@mdi/react";
import {mdiContentCopy} from "@mdi/js";

interface Props {
	open: boolean;
	config: string;
	onClose: () => void;
}

const KubeConfigDialog: React.FC<Props> = ({open, config, onClose}): JSX.Element => {
	// local state
	const [data, setData] = useState<string>("");

	return <Dialog
		open={open}
		TransitionProps={{
			onEnter: () => setData(() => config),
			onExited: () => setData(() => "")
		}}
		onClose={onClose}
		scroll="paper">
		<DialogTitle>
			Kubeconfig
		</DialogTitle>
		<DialogContent
			dividers>
			<DialogContentText>
				<Alert
					severity="info">
					This kubeconfig is intended for use by a human operator.
					For automated interaction, create a Service Account inside the cluster.
				</Alert>
				<Card
					variant="outlined"
					sx={{
						fontSize: 12,
						p: 1,
						pl: 2,
						color: "text.primary",
						overflowX: "auto",
						borderRadius: 1
					}}
					component="pre">
					<code>
						{data}
					</code>
				</Card>
				<Box>
					<Box sx={{display: "flex"}}>
						<Typography
							color="textPrimary">
							Setup guide
						</Typography>
						<Box sx={{flexGrow: 1}}/>
						<Button
							onClick={() => {
								void navigator.clipboard.writeText(data);
							}}
							startIcon={<Icon
								path={mdiContentCopy}
								size={0.7}
							/>}>
							Copy
						</Button>
					</Box>
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
