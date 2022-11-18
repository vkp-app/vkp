import React, {useEffect, useMemo, useState} from "react";
import {Box, Button, Card, Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import {ApolloError} from "@apollo/client";
import {AccessRef} from "../../../generated/graphql";
import InlineNotFound from "../../alert/InlineNotFound";
import InlineError from "../../alert/InlineError";
import AccessorItem from "./AccessorItem";
import AccessorDialog from "./AccessorDialog";

interface Props {
	accessors: AccessRef[];
	loading: boolean;
	error?: ApolloError;
	readOnly: boolean;
	onUpdate: (a: AccessRef[]) => void;
}

const AccessorList: React.FC<Props> = ({accessors, loading, error, readOnly, onUpdate}): JSX.Element => {
	// local state
	const [showAdd, setShowAdd] = useState<boolean>(false);
	const [newData, setNewData] = useState<AccessRef[]>(accessors);
	const [dirty, setDirty] = useState<boolean>(false);

	const handleReset = (): void => {
		setNewData(() => accessors);
		setDirty(() => false);
	}

	const handleUpdate = (): void => {
		onUpdate(newData);
	}

	useEffect(() => {
		handleReset();
	}, [accessors]);

	const accessData = useMemo(() => {
		if (loading || error)
			return [];
		const data = (newData || accessors);
		return data.map((c, idx) => <AccessorItem
			key={idx}
			item={c}
			readOnly={readOnly || data.length === 1}
			onDelete={() => {
				setNewData((d) => d.filter(i => i !== c));
				setDirty(() => true);
			}}
		/>);
	}, [loading, error, accessors, newData]);

	const loadingData = (): JSX.Element[] => {
		const items = [];
		for (let i = 0; i < 6; i++) {
			items.push(<AccessorItem
				key={i}
				item={null}
				readOnly
				onDelete={() => {}}
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
			{!loading && error && <InlineError
				error={error}
			/>}
			{accessors.length === 0 && <InlineNotFound
				title="No permissions"
			/>}
		</Card>
		<Box
			sx={{display: "flex", p: 1, pt: 2}}>
			<Button
				disabled={readOnly}
				onClick={() => setShowAdd(() => true)}>
				Add
			</Button>
			<Box sx={{flexGrow: 1}}/>
			<Button
				disabled={!dirty}
				onClick={handleReset}>
				Reset
			</Button>
			<Button
				sx={{ml: 1}}
				disabled={!dirty}
				onClick={handleUpdate}>
				Update
			</Button>
		</Box>
		<AccessorDialog
			open={showAdd}
			setOpen={setShowAdd}
			onAdd={v => {
				setNewData(d => [...d, v]);
				setDirty(() => true);
			}}
		/>
	</React.Fragment>
}
export default AccessorList;
