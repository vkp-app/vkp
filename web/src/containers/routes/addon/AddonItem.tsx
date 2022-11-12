import React from "react";
import {
	Avatar,
	Box,
	Button,
	CardActions,
	CardContent,
	CardHeader,
	CircularProgress,
	Skeleton,
	Typography
} from "@mui/material";
import Grid2 from "@mui/material/Unstable_Grid2";
import {makeStyles} from "tss-react/mui";
import {AddonSource, ClusterAddon} from "../../../generated/graphql";
import AddonSourceChip from "./AddonSourceChip";
import AddonChip from "./AddonChip";

const useStyles = makeStyles()(() => ({
	button: {
		fontFamily: "Manrope",
		fontWeight: 600,
		fontSize: 13,
		textTransform: "none",
		minHeight: 24,
		height: 24
	}
}));

interface Props {
	item: ClusterAddon | null;
	installed?: boolean;
	loading?: boolean;
	onInstall: () => void;
	onUninstall: () => void;
}

const AddonItem: React.FC<Props> = ({item, installed, loading, onInstall, onUninstall}): JSX.Element => {
	// hooks
	const {classes} = useStyles();

	const handleClick = (): void => {
		if (installed) {
			onUninstall();
			return;
		}
		onInstall();
	}

	return <Grid2
		xs={6}>
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
					{installed && <AddonChip/>}
				</span> : <Skeleton
					variant="text"
					width="60%"
				/>}
				titleTypographyProps={{fontFamily: "Manrope", fontWeight: "bold"}}
				subheaderTypographyProps={{fontSize: 14}}
				avatar={item != null ? <Avatar
					src={item.logo}
					alt={`${item.displayName} logo`}
					variant={item.logo !== "" ? "square" : "circular"}
				/> : <Skeleton
					variant="circular"
					width={48}
					height={48}
				/>}
				action={item != null ? <Button
					className={classes.button}
					variant="outlined"
					disabled={loading}
					startIcon={loading ? <CircularProgress size={14}/> : undefined}
					onClick={handleClick}>
					{installed ? "Uninstall" : "Install"}
				</Button> : undefined}
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
			{(item == null || item?.sourceURL !== "") && <CardActions>
				{item != null ? <Button
					className={classes.button}
					variant="text"
					href={item?.sourceURL || ""}
					target="_blank">
					View source
				</Button> : <Skeleton
					width="20%"
				/>}
			</CardActions>}
		</Box>
	</Grid2>
}
export default AddonItem;
