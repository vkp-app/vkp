import React from "react";
import {Card, CardHeader, IconButton, ListSubheader} from "@mui/material";
import {Link, useParams} from "react-router-dom";
import {ArrowLeft} from "tabler-icons-react";
import StandardLayout from "../layout/StandardLayout";
import {AccessRef, useCanEditClusterQuery, useClusterQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import AccessorList from "./access/AccessorList";

const ClusterAccessorList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();

	const clusterName = params["name"] || "";
	const tenantName = params["tenant"] || "";

	const cluster = useClusterQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !tenantName
	});
	const canEditCluster = useCanEditClusterQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !tenantName || !clusterName
	});

	return <StandardLayout>
		<ListSubheader
			sx={{display: "flex", alignItem: "center"}}>
			<IconButton
				size="small"
				component={Link}
				to={`/clusters/${tenantName}/cluster/${clusterName}`}>
				<ArrowLeft
					size={18}
				/>
			</IconButton>
			Back to cluster
		</ListSubheader>
		<Card
			variant="outlined"
			sx={{p: 2}}>
			<CardHeader
				title="Cluster permissions"
				titleTypographyProps={{fontFamily: "Figtree", fontWeight: 400, mb: 1}}
				subheader="Control who can access this cluster and what they can do."
				subheaderTypographyProps={{fontSize: 14}}
			/>
			{!cluster.loading && cluster.error && <InlineError
				message="Unable to load permissions"
				error={cluster.error}
			/>}
		</Card>
		<AccessorList
			accessors={(cluster.data?.cluster.accessors || []) as AccessRef[]}
			loading={cluster.loading}
			error={cluster.error}
			readOnly={!canEditCluster.data?.hasClusterAccess ?? true}
			onUpdate={() => {}}
		/>
	</StandardLayout>
}
export default ClusterAccessorList;
