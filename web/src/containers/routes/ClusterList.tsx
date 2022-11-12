import React, {useMemo} from "react";
import {
	Alert,
	Button,
	Card,
	CardHeader,
	Collapse,
	Link as MuiLink,
	Skeleton,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow
} from "@mui/material";
import {Link, useParams} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";
import InlineNotFound from "../alert/InlineNotFound";
import {Cluster, TenantPhase, useClustersQuery, useTenantQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import ClusterVersionIndicator from "./cluster/ClusterVersionIndicator";

const ClusterList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();
	const tenantName = params["tenant"] || "";

	const clusters = useClustersQuery({
		variables: {tenant: tenantName},
		skip: !tenantName
	});

	const tenant = useTenantQuery({
		variables: {tenant: tenantName},
		skip: !tenantName
	});

	const tenantApproved = tenant.data?.tenant.status.phase === TenantPhase.Ready;

	const clusterData = useMemo(() => {
		if (clusters.loading || clusters.error || !clusters.data)
			return [];
		return (clusters.data.clustersInTenant as Cluster[]).map(c => (<TableRow
			key={c.name}>
			<TableCell
				component="th"
				scope="row">
				<MuiLink
					component={Link}
					to={`/clusters/${c.tenant}/cluster/${c.name}`}>
					{c.name}
				</MuiLink>
			</TableCell>
			<TableCell
				sx={{display: "flex", alignItems: "center"}}
				align="right">
				<ClusterVersionIndicator version={c.status.kubeVersion}/>
			</TableCell>
		</TableRow>))
	}, [clusters]);

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 5; i++) {
			items.push(<TableRow
				key={i}>
				<TableCell
					component="th"
					scope="row">
					<Skeleton
						variant="text"
						width="60%"
					/>
				</TableCell>
				<TableCell
					align="right">
					<Skeleton
						sx={{float: "right"}}
						variant="text"
						width="10%"
					/>
				</TableCell>
			</TableRow>);
		}
		return items;
	}


	return <StandardLayout>
		<Card
			variant="outlined">
			<CardHeader
				title="Clusters"
				action={<Button
					sx={{fontFamily: "Manrope", fontSize: 13, fontWeight: 600, height: 24, minHeight: 24, textTransform: "none"}}
					variant="outlined"
					component={Link}
					to={`/new/cluster/${tenantName}`}
					disabled={tenant.loading || tenant.error != null || !tenantApproved}>
					Create
				</Button>}
			/>
			{!tenant.loading && tenant.error && <InlineError
				message="Unable to load tenant"
				error={tenant.error}
			/>}
			<Collapse
				in={!tenant.loading && !tenantApproved && !tenant.error}>
				<Alert
					severity="warning">
					Kubernetes clusters cannot be provisioned as this tenancy is awaiting approval from an administrator or policy agent.
				</Alert>
			</Collapse>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell>Name</TableCell>
						<TableCell align="right">Version</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{clusters.loading ? loadingData() : clusterData}
				</TableBody>
			</Table>
			{!clusters.loading && clusters.error && <InlineError
				message="Unable to load clusters"
				error={clusters.error}
			/>}
			{clusters.data?.clustersInTenant.length === 0 && <InlineNotFound
				title="No clusters"
				subtitle="This tenancy is empty. Create a cluster to get started"
			/>}
		</Card>
	</StandardLayout>
}
export default ClusterList;
