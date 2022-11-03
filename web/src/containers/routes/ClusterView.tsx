import React from "react";
import StandardLayout from "../layout/StandardLayout";
import {useParams} from "react-router-dom";
import {useClusterQuery} from "../../generated/graphql";
import {Card, CardHeader, Typography} from "@mui/material";

const ClusterView: React.FC = (): JSX.Element => {
	const params = useParams();
	const clusterName = params["name"];
	const tenantName = params["tenant"];

	const {data, loading, error} = useClusterQuery({
		variables: {tenant: tenantName || "", cluster: clusterName || ""},
		skip: !clusterName || !tenantName
	});

	return <StandardLayout>
		<Typography>
			Cluster overview
		</Typography>
		<Card>
			<CardHeader>
				{data?.cluster.name}
			</CardHeader>
		</Card>
	</StandardLayout>
}
export default ClusterView;
