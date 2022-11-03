import React from "react";
import {Cluster} from "../../../generated/graphql";
import {Card, List, ListItem, ListItemText, Skeleton} from "@mui/material";

interface Props {
	cluster: Cluster | null;
	loading: boolean;
}

const ClusterMetadataView: React.FC<Props> = ({cluster, loading}): JSX.Element => {
	return <Card
		variant={"outlined"}
		sx={{p: 2}}>
		<List>
			<ListItem>
				<ListItemText
					primary={"Management API"}
					secondary={loading ? <Skeleton variant={"text"} height={20} width={"40%"}/> : "VCluster"}
				/>
			</ListItem>
		</List>
	</Card>
}
export default ClusterMetadataView;
