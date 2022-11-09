import React, {useMemo} from "react";
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
import {Link} from "react-router-dom";
import StandardLayout from "../layout/StandardLayout";
import InlineNotFound from "../alert/InlineNotFound";
import {Tenant, useTenantsQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";

const TenantList: React.FC = (): JSX.Element => {
	// hooks
	const {data, loading, error} = useTenantsQuery();

	const tenants = useMemo(() => {
		if (loading || error || !data)
			return [];
		return (data.tenants as Tenant[]).map(c => (<TableRow
			key={c.name}>
			<TableCell
				component="th"
				scope="row">
				<MuiLink
					component={Link}
					to={`/clusters/${c.name}`}>
					{c.name}
				</MuiLink>
			</TableCell>
			<TableCell
				align="right">
				{c.owner}
			</TableCell>
			<TableCell
				sx={{display: "flex", alignItems: "center"}}
				align="right">
				{c.status.phase}
			</TableCell>
		</TableRow>))
	}, [data, loading, error]);

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
		<ListSubheader
			sx={{display: "flex", alignItems: "center"}}>
			Tenants
			<Box sx={{flexGrow: 1}}/>
			<Button
				sx={{fontFamily: "Manrope", fontSize: 13, fontWeight: 600, height: 24, minHeight: 24, textTransform: "none"}}
				variant="outlined"
				component={Link}
				to="/new/tenant">
				Create
			</Button>
		</ListSubheader>
		<Card
			variant="outlined">
			<Table>
				<TableHead>
					<TableRow>
						<TableCell>Name</TableCell>
						<TableCell align="right">Owner</TableCell>
						<TableCell align="right">Status</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{loading ? loadingData() : tenants}
				</TableBody>
			</Table>
			{!loading && error && <InlineError error={error}/>}
			{data?.tenants.length === 0 && <InlineNotFound
				title="No tenants"
				subtitle="Create a tenant to get started"
			/>}
		</Card>
	</StandardLayout>
}
export default TenantList;
