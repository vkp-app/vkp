import React from "react";
import {Card, Grid, ListItem, ListItemText} from "@mui/material";
import SparkLine from "../../data/SparkLine";
import {Cluster, useMetricsClusterQuery} from "../../../generated/graphql";

interface Props {
	cluster: Cluster | null;
	loading: boolean;
}

const ClusterMetricsView: React.FC<Props> = ({cluster, loading}): JSX.Element => {
	const {data} = useMetricsClusterQuery({
		variables: {tenant: cluster?.tenant || "", cluster: cluster?.name || ""},
		skip: !cluster
	});

	return <Card
		variant="outlined"
		sx={{p: 2}}>
		<Grid container>
			<Grid item xs={6}>
				<ListItem>
					<ListItemText
						secondary="Cluster memory usage"
						primary={<SparkLine
							width={300}
							data={data?.clusterMetricMemory.map(i => Number(i.value)) || []}
						/>}
					/>
				</ListItem>
			</Grid>
			<Grid item sx={6}>
				<ListItem>
					<ListItemText
						secondary="Cluster CPU usage"
						primary={<SparkLine
							width={300}
							data={data?.clusterMetricCPU.map(i => Number(i.value)) || []}
						/>}
					/>
				</ListItem>
			</Grid>
			<Grid item sx={6}>
				<ListItem>
					<ListItemText
						secondary="Pod count"
						primary={<SparkLine
							width={300}
							data={data?.clusterMetricPods.map(i => Number(i.value)) || []}
						/>}
					/>
				</ListItem>
			</Grid>
		</Grid>
	</Card>
}
export default ClusterMetricsView;
