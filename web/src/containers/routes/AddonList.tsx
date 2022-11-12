import React, {useMemo} from "react";
import {Avatar, Box, Card, CardContent, CardHeader, IconButton, ListSubheader, Typography} from "@mui/material";
import {Link, useParams} from "react-router-dom";
import {ArrowLeft} from "tabler-icons-react";
import Grid2 from "@mui/material/Unstable_Grid2";
import StandardLayout from "../layout/StandardLayout";
import {ClusterAddon, useAllAddonsQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";

const AddonList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();

	const clusterName = params["name"] || "";
	const tenantName = params["tenant"] || "";

	const addons = useAllAddonsQuery({
		variables: {tenant: tenantName},
		skip: !tenantName
	});

	const addonData = useMemo(() => {
		if (addons.loading || addons.error || !addons.data)
			return [];
		return (addons.data.clusterAddons as ClusterAddon[]).map(c => (<Grid2
			xs={4}
			key={c.name}>
			<Box>
				<CardHeader
					title={c.displayName}
					subheader={c.maintainer}
					subheaderTypographyProps={{fontSize: 14}}
					avatar={<Avatar
						src={c.logo}
						alt={`${c.displayName} logo`}
						variant="square"
					/>}
				/>
				<CardContent>
					<Typography
						variant="body2"
						color="text.secondary">
						{c.description || <i>No description provided.</i>}
					</Typography>
				</CardContent>
			</Box>
		</Grid2>))
	}, [addons]);

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
				{addonData}
			</Grid2>
		</Card>
	</StandardLayout>
}
export default AddonList;
