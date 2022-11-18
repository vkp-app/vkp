import React from "react";
import {Checkbox, IconButton, Skeleton, TableCell, TableRow} from "@mui/material";
import {Trash} from "tabler-icons-react";
import {AccessRef} from "../../../generated/graphql";

interface Props {
	item: AccessRef | null;
	readOnly: boolean;
	onDelete: () => void;
}

const AccessorItem: React.FC<Props> = ({item, readOnly, onDelete}): JSX.Element => {
	const refName = item?.user || item?.group;
	const refType = item?.user ? "User" : "Group";

	return <TableRow>
		<TableCell
			component="th"
			scope="row">
			{item != null ? refName : <Skeleton width="70%"/>}
		</TableCell>
		<TableCell
			align="right">
			{item != null ? refType : <Skeleton sx={{float: "right"}} width="20%"/>}
		</TableCell>
		<TableCell
			sx={{display: "flex", alignItems: "center"}}
			align="right">
			{item ? <Checkbox
				checked={item?.readOnly ?? true}
				disabled
			/> : <Skeleton
				width={32}
				height={32}
				variant="circular"
			/>}
		</TableCell>
		<TableCell
			align="right">
			<IconButton
				disabled={readOnly}
				onClick={onDelete}
				color="error">
				<Trash/>
			</IconButton>
		</TableCell>
	</TableRow>
}
export default AccessorItem;
