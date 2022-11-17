import React, {ReactNode} from "react";
import {Card, CardHeader, Grid, Skeleton, Typography} from "@mui/material";
import {ApolloError} from "@apollo/client";
import SparkLine, {SparklineColours} from "../../data/SparkLine";
import {Cluster, MetricFormat, MetricValue, useMetricsClusterQuery} from "../../../generated/graphql";
import {formatBytes} from "../../../utils/fmt";
import InlineError from "../../alert/InlineError";

interface Props {
	cluster: Cluster | null;
	clusterError?: ApolloError;
	refresh?: boolean;
}

const colours: SparklineColours[] = ["primary", "success", "warning", "error"];

const ClusterMetricsView: React.FC<Props> = ({cluster, refresh, clusterError}): JSX.Element => {
	const {data, loading, error} = useMetricsClusterQuery({
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

	const getFormatter = (f: MetricFormat): (n: number) => string => {
		switch (f) {
			case MetricFormat.Bytes:
				return fmtBytes;
			case MetricFormat.Cpu:
				return fmtCPU;
			case MetricFormat.Time:
			case MetricFormat.Plain:
			default:
				return fmtPlain;
		}
	}

	const metric = (colour: SparklineColours, data: MetricValue[], name: string, bz: boolean, fmt: (n: number) => string): ReactNode => {
		const numData = data.map(i => Number(i.value));
		const last = data.length === 0 ? 0 : Number(data[data.length - 1].value);
		const max = Math.max(...numData);
		return <Grid
			key={name}
			item
			xs={6}>
			<CardHeader
				title={<SparkLine
					color={colour}
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
		{!loading && (error || clusterError) && <InlineError
			message="Unable to load cluster metrics"
			error={error || clusterError}
		/>}
		<Grid container>
			{loading && loadingData()}
			{!loading && data?.clusterMetrics.map((m, i) => metric(colours[i % colours.length], m.values as MetricValue[], m.name, false, getFormatter(m.format)))}
		</Grid>
	</Card>
}
export default ClusterMetricsView;
