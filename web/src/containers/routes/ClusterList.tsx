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
import {Link, useNavigate, useParams} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";
import InlineNotFound from "../alert/InlineNotFound";
import {Cluster, TenantPhase, useApproveTenancyMutation, useClusterListQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import Loadable from "../data/Loadable";
import BackButton from "../layout/BackButton";
import ClusterVersionIndicator from "./cluster/ClusterVersionIndicator";

const ClusterList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();
	const navigate = useNavigate();
	const tenantName = params["tenant"] || "";

	const clusterList = useClusterListQuery({
		variables: {tenant: tenantName},
		skip: !tenantName
	});
	const [approveTenant] = useApproveTenancyMutation();

	const tenantApproved = clusterList.data?.tenant.status.phase === TenantPhase.Ready;

	const onApproveTenant = (): void => {
		approveTenant({
			variables: {tenant: tenantName}
		}).then(r => {
			if (r.data) {
				navigate("/tenants");
			}
		});
	}

	const clusterData = useMemo(() => {
		if (clusterList.loading || clusterList.error || !clusterList.data)
			return [];
		return (clusterList.data.clustersInTenant as Cluster[]).map(c => (<TableRow
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
	}, [clusterList]);

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
		<BackButton
			title="Back to tenants"
			to="/tenants"
		/>
		<Card>
			<CardHeader
				title="Clusters"
				titleTypographyProps={{fontFamily: "Figtree"}}
				action={<Button
					sx={{fontFamily: "Manrope", fontSize: 13, fontWeight: 600, height: 24, minHeight: 24, textTransform: "none"}}
					variant="outlined"
					component={Link}
					to={`/new/cluster/${tenantName}`}
					disabled={clusterList.loading || clusterList.error != null || !tenantApproved || !clusterList.data?.hasTenantAccess}>
					Create
				</Button>}
			/>
			{!clusterList.loading && clusterList.error && <InlineError
				message="Unable to load tenant"
				error={clusterList.error}
			/>}
			<Collapse
				in={!clusterList.loading && !tenantApproved && !clusterList.error}>
				<Alert
					severity="warning">
					Kubernetes clusters cannot be provisioned as this tenancy is awaiting approval from an administrator or policy agent.
				</Alert>
				{clusterList.data?.hasRole === true && <Alert
					action={<Button
						onClick={() => onApproveTenant()}>
						Approve
					</Button>}
					severity="info">
					You are an administrator and can approve this tenancy.
				</Alert>}
			</Collapse>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell>Name</TableCell>
						<TableCell align="right">Version</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					<Loadable data={clusterData} skeleton={loadingData()} loading={clusterList.loading}/>
				</TableBody>
			</Table>
			{!clusterList.loading && clusterList.error && <InlineError
				message="Unable to load clusters"
				error={clusterList.error}
			/>}
			{clusterList.data?.clustersInTenant.length === 0 && <InlineNotFound
				title="No clusters"
				subtitle="This tenancy is empty. Create a cluster to get started"
			/>}
		</Card>
	</StandardLayout>
}
export default ClusterList;
