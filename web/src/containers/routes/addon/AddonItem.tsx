import React from "react";
import {Avatar, Box, CardContent, CardHeader, Skeleton, Typography} from "@mui/material";
import Grid2 from "@mui/material/Unstable_Grid2";
import {AddonSource, ClusterAddon} from "../../../generated/graphql";
import AddonSourceChip from "./AddonSourceChip";

interface Props {
	item: ClusterAddon | null;
}

const AddonItem: React.FC<Props> = ({item}): JSX.Element => {
	return <Grid2
		xs={4}>
		<Box>
			<CardHeader
				disableTypography={item == null}
				title={item?.displayName || <Skeleton
					variant="text"
					width="40%"
				/>}
				subheader={item?.maintainer ? <span>
					{item.maintainer}
					<AddonSourceChip source={item?.source || AddonSource.Unknown}/>
				</span> : <Skeleton
					variant="text"
					width="60%"
				/>}
				subheaderTypographyProps={{fontSize: 14}}
				avatar={item != null ? <Avatar
					src={item.logo}
					alt={`${item.displayName} logo`}
					variant="square"
				/> : <Skeleton
					variant="circular"
					width={48}
					height={48}
				/>}
			/>
			<CardContent>
				{item != null ? <Typography
					variant="body2"
					color="text.secondary">
					{item.description || <i>No description provided.</i>}
				</Typography> : <Skeleton
					variant="text"
					width="40%"
				/>}
			</CardContent>
		</Box>
	</Grid2>
}
export default AddonItem;
