import React, {ReactNode} from "react";
import {Card, CardHeader, Grid, Skeleton, Typography} from "@mui/material";
import SparkLine from "../../data/SparkLine";
import {Cluster, MetricValue, useMetricsClusterQuery} from "../../../generated/graphql";
import {formatBytes} from "../../../utils/fmt";

interface Props {
	cluster: Cluster | null;
	loading: boolean;
	refresh?: boolean;
}

const ClusterMetricsView: React.FC<Props> = ({cluster, loading, refresh}): JSX.Element => {
	const {data} = useMetricsClusterQuery({
		variables: {tenant: cluster?.tenant || "", cluster: cluster?.name || ""},
		skip: !cluster,
		pollInterval: refresh ? 10_000 : 0
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

	const getFormatter = (metric: string): (n: number) => string => {
		if (metric.includes("bytes")) {
			return fmtBytes;
		}
		if (metric.includes("seconds") && metric.includes("cpu")) {
			return fmtCPU;
		}

		return fmtPlain;
	}

	const metric = (data: MetricValue[], name: string, bz: boolean, fmt: (n: number) => string): ReactNode => {
		const numData = data.map(i => Number(i.value));
		const last = data.length === 0 ? 0 : Number(data[data.length - 1].value);
		const max = Math.max(...numData);
		return <Grid
			key={name}
			item
			xs={6}>
			<CardHeader
				title={<SparkLine
					width={1000}
					data={numData}
					baseZero={bz}
				/>}
				subheader={<Typography
					sx={{fontSize: 14}}
					color="text.secondary">
					{name}&nbsp;({fmt(last)}/{fmt(max)})
				</Typography>}
				disableTypography
			/>
		</Grid>
	}

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 4; i++) {
			items.push(<Grid
				key={i}
				item
				xs={6}>
				<CardHeader
					title={<Skeleton height={35}/>}
					subheader={<Skeleton width="50%"/>}
					disableTypography
				/>
			</Grid>);
		}
		return items;
	}

	return <Card
		variant="outlined"
		sx={{p: 2}}>
		<Grid container>
			{loading && loadingData()}
			{!loading && data?.clusterMetrics.map(m => metric(m.values as MetricValue[], m.name, false, getFormatter(m.metric)))}
		</Grid>
	</Card>
}
export default ClusterMetricsView;
