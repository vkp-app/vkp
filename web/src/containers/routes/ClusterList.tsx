import React, {useMemo} from "react";
import StandardLayout from "../layout/StandardLayout";
import {
	Box,
	Button,
	Card,
	Link as MuiLink,
	ListSubheader,
	Skeleton,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow
} from "@mui/material";
import InlineNotFound from "../alert/InlineNotFound";
import {Cluster, useClustersQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import {Link, useParams} from "react-router-dom";
import ClusterVersionIndicator from "./cluster/ClusterVersionIndicator";

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
				<ClusterVersionIndicator version={c.status.kubeVersion}/>
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
		<ListSubheader
			sx={{display: "flex", alignItems: "center"}}>
			Clusters
			<Box sx={{flexGrow: 1}}/>
			<Button
				sx={{fontFamily: "Manrope", fontSize: 13, fontWeight: 600, height: 24, minHeight: 24, textTransform: "none"}}
				variant={"outlined"}
				component={Link}
				to={"/new/cluster"}>
				Create
			</Button>
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
