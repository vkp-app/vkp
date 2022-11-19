import React from "react";
import {Card, CardHeader} from "@mui/material";
import {useParams} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";
import {
	AccessRef,
	useCanEditClusterQuery,
	useClusterQuery,
	useSetClusterAccessorsMutation
} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import BackButton from "../layout/BackButton";
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

	const [setAccessors, setAccessorsResult] = useSetClusterAccessorsMutation();

	const handleSetAccessors = (r: AccessRef[]): void => {
		setAccessors({
			variables: {tenant: tenantName, cluster: clusterName, accessors: r}
		}).then(r => {
			if (r.data) {
				void cluster.refetch();
			}
		});
	}

	return <StandardLayout>
		<BackButton
			title="Back to cluster"
			to={`/clusters/${tenantName}/cluster/${clusterName}`}
		/>
		<Card
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
			loading={setAccessorsResult.loading || cluster.loading}
			error={setAccessorsResult.error || cluster.error}
			readOnly={!canEditCluster.data?.hasClusterAccess ?? true}
			onUpdate={handleSetAccessors}
		/>
	</StandardLayout>
}
export default ClusterAccessorList;
