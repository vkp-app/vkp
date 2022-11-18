import React, {useMemo} from "react";
import {Box, Button, Card, Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import {ApolloError} from "@apollo/client";
import {AccessRef} from "../../../generated/graphql";
import InlineNotFound from "../../alert/InlineNotFound";
import AccessorItem from "./AccessorItem";

interface Props {
	accessors: AccessRef[];
	loading: boolean;
	error?: ApolloError;
	readOnly: boolean;
}

const AccessorList: React.FC<Props> = ({accessors, loading, error, readOnly}): JSX.Element => {

	const accessData = useMemo(() => {
		if (loading || error)
			return [];
		return accessors.map((c, idx) => <AccessorItem
			key={idx}
			item={c}
			readOnly={readOnly}
		/>);
	}, [loading, error, accessors]);

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 6; i++) {
			items.push(<AccessorItem
				key={i}
				item={null}
				readOnly
			/>);
		}
		return items;
	}

	return <React.Fragment>
		<Card
			variant="outlined"
			sx={{p: 2, mt: 2}}>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell>Name</TableCell>
						<TableCell align="right">Type</TableCell>
						<TableCell align="right">Read-only</TableCell>
						<TableCell align="right">Actions</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{loading ? loadingData() : accessData}
				</TableBody>
			</Table>
			{accessors.length === 0 && <InlineNotFound
				title="No permissions"
			/>}
		</Card>
		<Box
			sx={{display: "flex", p: 1, pt: 2}}>
			<Button
				disabled>
				Add
			</Button>
			<Box sx={{flexGrow: 1}}/>
			<Button
				disabled>
				Update
			</Button>
		</Box>
	</React.Fragment>
}
export default AccessorList;
