import React, {ReactNode} from "react";
import {Card, CircularProgress, Grid, ListItem, ListItemText} from "@mui/material";
import SparkLine from "../../data/SparkLine";
import {Cluster, MetricValue, useMetricsClusterQuery} from "../../../generated/graphql";
import {formatBytes} from "../../../utils/fmt";

interface Props {
	cluster: Cluster | null;
	loading: boolean;
}

const ClusterMetricsView: React.FC<Props> = ({cluster, loading}): JSX.Element => {
	const {data} = useMetricsClusterQuery({
		variables: {tenant: cluster?.tenant || "", cluster: cluster?.name || ""},
		skip: !cluster
	});

	/**
	 * formats a number as bytes
	 */
	const fmtBytes = (n: number): string => {
		return formatBytes(n, 0);
	}

	/**
	 * formats a number in CPU millicores.
	 * 1 or more cores switches to decimal values
	 * @param n
	 */
	const fmtCPU = (n: number): string => {
		if (!+n || !Number.isFinite(n)) {
			return "0m";
		}
		const cores = Math.floor(n);
		const ms = ((n % 1) * 1000).toFixed(1);

		if (cores === 0) {
			return `${ms}m`;
		}
		return `${n.toFixed(2)} CPU`;
	}

	/**
	 * formats a number as-is
	 */
	const fmtPlain = (n: number): string => {
		if (!+n || !Number.isFinite(n)) {
			return "0";
		}
		return `${n}`;
	}

	const metric = (data: MetricValue[], name: string, bz: boolean, fmt: (n: number) => string): ReactNode => {
		const numData = data.map(i => Number(i.value));
		const last = data.length === 0 ? 0 : Number(data[data.length - 1].value);
		const max = Math.max(...numData);
		return <ListItem>
			<ListItemText
				secondary={`${name} (${fmt(last)}/${fmt(max)})`}
				primary={loading ? <CircularProgress/> : <SparkLine
					width={300}
					data={numData}
					baseZero={bz}
				/>}
			/>
		</ListItem>;
	}

	return <Card
		variant="outlined"
		sx={{p: 2}}>
		<Grid container>
			<Grid item xs={6}>
				{metric(data?.clusterMetricMemory || [], "Memory usage", false, fmtBytes)}
			</Grid>
			<Grid item xs={6}>
				{metric(data?.clusterMetricCPU || [], "CPU usage", false, fmtCPU)}
			</Grid>
			<Grid item xs={6}>
				{metric(data?.clusterMetricPods || [], "Pod count", false, fmtPlain)}
			</Grid>
			<Grid item xs={6}>
				{metric(data?.clusterMetricNetReceive || [], "Network received", true, fmtBytes)}
			</Grid>
		</Grid>
	</Card>
}
export default ClusterMetricsView;
