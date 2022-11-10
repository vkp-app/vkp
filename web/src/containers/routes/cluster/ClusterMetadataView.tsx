import React from "react";
import {Alert, Card, List, ListItem, ListItemText, Skeleton} from "@mui/material";
import {Cluster} from "../../../generated/graphql";

interface Props {
	cluster: Cluster | null;
	loading: boolean;
}

const ClusterMetadataView: React.FC<Props> = ({cluster, loading}): JSX.Element => {
	return <Card
		variant="outlined"
		sx={{p: 2}}>
		<Alert
			severity="warning">
			This information is not normally needed but may be useful for debugging or support-related purposes.
		</Alert>
		<List>
			<ListItem>
				<ListItemText
					primary="Release track"
					secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : cluster?.track || "Unknown"}
				/>
			</ListItem>
			<ListItem>
				<ListItemText
					primary="Addons"
					secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : JSON.stringify(cluster?.addons || [])}
				/>
			</ListItem>
			<ListItem>
				<ListItemText
					primary="Management API"
					secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : "VCluster"}
				/>
			</ListItem>
			<ListItem>
				<ListItemText
					primary="Management namespace"
					secondary={loading ? <Skeleton variant="text" height={20} width="40%"/> : cluster?.tenant || "Unknown"}
				/>
			</ListItem>
		</List>
	</Card>
}
export default ClusterMetadataView;
