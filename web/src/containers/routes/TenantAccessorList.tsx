import React, {useMemo} from "react";
import {Card, CardHeader, IconButton, ListSubheader} from "@mui/material";
import {Link, useParams} from "react-router-dom";
import {ArrowLeft} from "tabler-icons-react";
import StandardLayout from "../layout/StandardLayout";
import {AccessRef, useSetTenantAccessorsMutation, useTenantAccessQuery} from "../../generated/graphql";
import InlineError from "../alert/InlineError";
import AccessorList from "./access/AccessorList";

const TenantAccessorList: React.FC = (): JSX.Element => {
	// hooks
	const params = useParams();

	const tenantName = params["tenant"] || "";

	const tenant = useTenantAccessQuery({
		variables: {tenant: tenantName},
		skip: !tenantName
	});

	const [setAccessors, setAccessorsResult] = useSetTenantAccessorsMutation();

	const handleSetAccessors = (r: AccessRef[]): void => {
		setAccessors({
			variables: {tenant: tenantName, accessors: r}
		}).then(r => {
			if (r.data) {
				void tenant.refetch();
			}
		})
	}

	const accessors = useMemo(() => {
		if (tenant.data == null)
			return [];
		return (tenant.data?.tenant.accessors || []) as AccessRef[];
	}, [tenant]);

	return <StandardLayout>
		<ListSubheader
			sx={{display: "flex", alignItem: "center"}}>
			<IconButton
				size="small"
				component={Link}
				to={`/clusters/${tenantName}`}>
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
				title="Tenant permissions"
				titleTypographyProps={{fontFamily: "Figtree", fontWeight: 400, mb: 1}}
				subheader="Control who can access this tenant and what they can do."
				subheaderTypographyProps={{fontSize: 14}}
			/>
			{!tenant.loading && tenant.error && <InlineError
				message="Unable to load permissions"
				error={tenant.error}
			/>}
		</Card>
		<AccessorList
			accessors={accessors}
			loading={setAccessorsResult.loading || tenant.loading}
			error={setAccessorsResult.error || tenant.error}
			readOnly={!tenant.data?.hasTenantAccess ?? true}
			onUpdate={handleSetAccessors}
		/>
	</StandardLayout>
}
export default TenantAccessorList;
