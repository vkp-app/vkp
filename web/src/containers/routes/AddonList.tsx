import React, {useMemo} from "react";
import {Card, IconButton, ListSubheader} from "@mui/material";
import {Link, useParams} from "react-router-dom";
import {ArrowLeft} from "tabler-icons-react";
import Grid2 from "@mui/material/Unstable_Grid2";
import StandardLayout from "../layout/StandardLayout";
import {ClusterAddon, useAllAddonsQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import AddonItem from "./addon/AddonItem";

const AddonList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();

	const clusterName = params["name"] || "";
	const tenantName = params["tenant"] || "";

	const addons = useAllAddonsQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !tenantName
	});

	const addonData = useMemo(() => {
		if (addons.loading || addons.error || !addons.data)
			return [];
		return (addons.data.clusterAddons as ClusterAddon[]).map(c => <AddonItem
			key={c.name}
			item={c}
			installed={addons.data?.clusterInstalledAddons.find(i => i === c.name) != null}
		/>);
	}, [addons]);

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 9; i++) {
			items.push(<AddonItem key={i} item={null}/>);
		}
		return items;
	}

	return <StandardLayout>
		<ListSubheader
			sx={{display: "flex", alignItem: "center"}}>
			<IconButton
				size="small"
				centerRipple={false}
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
			{!addons.loading && addons.error && <InlineError
				message="Unable to load addons"
				error={addons.error}
			/>}
			<Grid2
				container
				spacing={2}>
				{addons.loading ? loadingData() : addonData}
			</Grid2>
		</Card>
	</StandardLayout>
}
export default AddonList;
