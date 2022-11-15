import React, {useMemo} from "react";
import {
	Button,
	Card,
	CardHeader,
	Link as MuiLink,
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
		<Card
			variant="outlined">
			<CardHeader
				title="Tenants"
				action={<Button
					component={Link}
					to="/new/tenant">
					Create
				</Button>}
			/>
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
			{!loading && error && <InlineError
				message="Unable to load tenants"
				error={error}
			/>}
			{data?.tenants.length === 0 && <InlineNotFound
				title="No tenants"
				subtitle="Create a tenant to get started"
			/>}
		</Card>
	</StandardLayout>
}
export default TenantList;
