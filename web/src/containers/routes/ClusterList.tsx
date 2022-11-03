import React, {useMemo} from "react";
import StandardLayout from "../layout/StandardLayout";
import {
	Card,
	IconButton,
	ListSubheader,
	Skeleton,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow,
	Tooltip,
	Link as MuiLink
} from "@mui/material";
import InlineNotFound from "../alert/InlineNotFound";
import {Cluster, useClustersQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import {Help} from "tabler-icons-react";
import {Link, useParams} from "react-router-dom";

const ClusterList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();
	const tenantName = params["tenant"];

	const {data, loading, error} = useClustersQuery({
		variables: {tenant: tenantName || ""},
		skip: !tenantName
	});

	const clusterData = useMemo(() => {
		if (loading || error || !data)
			return [];
		return (data.clustersInTenant as Cluster[]).map(c => (<TableRow
			key={c.name}>
			<TableCell
				component={"th"}
				scope={"row"}>
				<MuiLink
					component={Link}
					to={`/clusters/${c.tenant}/cluster/${c.name}`}>
					{c.name}
				</MuiLink>
			</TableCell>
			<TableCell
				sx={{display: "flex", alignItems: "center"}}
				align={"right"}>
				<Tooltip title={"View Kubernetes version information (external)"}>
					<IconButton
						sx={{ml: 1}}
						size={"small"}
						href={`https://kubernetes.io/releases/#release-v${c.status.kubeVersion.replace(".", "-")}`}
						target={"_blank"}>
						<Help
							size={18}
						/>
					</IconButton>
				</Tooltip>
				{c.status.kubeVersion}
			</TableCell>
		</TableRow>))
	}, [data, loading, error]);

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 5; i++) {
			items.push(<TableRow
				key={i}>
				<TableCell
					component={"th"}
					scope={"row"}>
					<Skeleton
						variant={"text"}
						width={"60%"}
					/>
				</TableCell>
				<TableCell
					align={"right"}>
					<Skeleton
						sx={{float: "right"}}
						variant={"text"}
						width={"10%"}
					/>
				</TableCell>
			</TableRow>);
		}
		return items;
	}


	return <StandardLayout>
		<ListSubheader>
			Clusters
		</ListSubheader>
		<Card
			variant={"outlined"}>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell>Name</TableCell>
						<TableCell align="right">Version</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{loading ? loadingData() : clusterData}
				</TableBody>
			</Table>
			{!loading && error && <InlineError error={error}/>}
			{data?.clustersInTenant.length === 0 && <InlineNotFound
				title={"No clusters"}
				subtitle={"This tenancy is empty. Create a cluster to get started"}
			/>}
		</Card>
	</StandardLayout>
}
export default ClusterList;
