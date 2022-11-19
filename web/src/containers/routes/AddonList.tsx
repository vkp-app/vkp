import React, {useMemo} from "react";
import {Card, CardHeader, IconButton, ListSubheader} from "@mui/material";
import {Link, useParams} from "react-router-dom";
import {ArrowLeft} from "tabler-icons-react";
import Grid2 from "@mui/material/Unstable_Grid2";
import StandardLayout from "../layout/StandardLayout";
import {
	ClusterAddon,
	useAllAddonsQuery,
	useCanEditClusterQuery,
	useInstallAddonMutation,
	useUninstallAddonMutation
} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import Loadable from "../data/Loadable";
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
	const canEditCluster = useCanEditClusterQuery({
		variables: {tenant: tenantName, cluster: clusterName},
		skip: !tenantName || !clusterName
	});

	const [installAddon, installData] = useInstallAddonMutation();
	const [uninstallAddon, uninstallData] = useUninstallAddonMutation();

	const onInstallAddon = (name: string): void => {
		installAddon({
			variables: {tenant: tenantName, cluster: clusterName, addon: name}
		}).then(r => {
			if (r.data != null) {
				void addons.refetch();
			}
		});
	}

	const onUninstallAddon = (name: string): void => {
		void uninstallAddon({
			variables: {tenant: tenantName, cluster: clusterName, addon: name}
		}).then(r => {
			if (r.data != null) {
				void addons.refetch();
			}
		});
	}

	const addonData = useMemo(() => {
		if (addons.loading || addons.error || !addons.data)
			return [];
		return (addons.data.clusterAddons as ClusterAddon[]).map(c => <AddonItem
			key={c.name}
			item={c}
			phase={addons.data?.clusterInstalledAddons.find(i => i.name === `${clusterName}-${c.name}`)?.phase}
			loading={installData.loading || uninstallData.loading}
			readOnly={!canEditCluster.data?.hasClusterAccess ?? true}
			onInstall={() => onInstallAddon(c.name)}
			onUninstall={() => onUninstallAddon(c.name)}
		/>);
	}, [addons, installData, uninstallData]);

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 6; i++) {
			items.push(<AddonItem
				key={i}
				item={null}
				onInstall={() => {}}
				onUninstall={() => {}}
			/>);
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
			sx={{p: 2}}>
			<CardHeader
				title="Addon marketplace"
				titleTypographyProps={{fontFamily: "Figtree", fontWeight: 400, mb: 1}}
				subheader="Browse, install and modify pre-packaged applications, configuration and capability provided by us, your administrators and the community."
				subheaderTypographyProps={{fontSize: 14}}
			/>
			{!addons.loading && addons.error && <InlineError
				message="Unable to load addons"
				error={addons.error}
			/>}
		</Card>
		<Grid2
			sx={{mt: 2}}
			container
			spacing={2}>
			<Loadable data={addonData} skeleton={loadingData()} loading={addons.loading}/>
		</Grid2>
	</StandardLayout>
}
export default AddonList;
